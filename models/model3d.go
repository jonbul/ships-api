package models

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"

	"ships3d/i18n"
)

// Model3d is a document of the `models3d` collection: a 3D model a user
// imported as a .glb file, shown by ship parts of type PartModel. Every
// player in a game downloads the models of the ships they see, so a model is
// public once uploaded (its id is all it takes to fetch it) and never
// changes, which lets browsers cache it for good.
type Model3d struct {
	Id     bson.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserId string        `json:"userId" bson:"userId"`
	Name   string        `json:"name" bson:"name"`
	// Hash is the SHA-256 of Data, so importing the same file twice reuses
	// the first upload.
	Hash        string `json:"-" bson:"hash"`
	Data        []byte `json:"-" bson:"data"`
	DateCreated int64  `json:"dateCreated" bson:"dateCreated"`
}

// Limits on imported models. A model is sent to every player who sees the
// ship and drawn by phones, so it has to stay small. Mirrored in
// ships-vue-3d's shipModel.ts.
const (
	MaxModelMegabytes  = 4
	MaxModelBytes      = MaxModelMegabytes << 20
	MaxModelTriangles  = 100_000
	MaxModelsPerUser   = 100
	MaxModelNameLength = 60
)

// Extensions three.js' GLTFLoader can only read with decoders the client
// doesn't load. A model that requires one would never show.
var unsupportedExtensions = map[string]bool{
	"KHR_draco_mesh_compression": true,
	"EXT_meshopt_compression":    true,
	"KHR_meshopt_compression":    true,
	"KHR_texture_basisu":         true,
}

const (
	glbMagic          = 0x46546C67 // "glTF"
	glbChunkJSON      = 0x4E4F534A // "JSON"
	glbHeader         = 12
	glbChunk          = 8
	glTFTriangles     = 4
	glTFTriangleStrip = 5
	glTFTriangleFan   = 6
)

type gltfDocument struct {
	Asset struct {
		Version string `json:"version"`
	} `json:"asset"`
	Buffers []struct {
		URI string `json:"uri"`
	} `json:"buffers"`
	Images []struct {
		URI string `json:"uri"`
	} `json:"images"`
	Accessors []struct {
		Count int64 `json:"count"`
	} `json:"accessors"`
	Meshes []struct {
		Primitives []struct {
			Attributes map[string]int `json:"attributes"`
			Indices    *int           `json:"indices"`
			Mode       *int           `json:"mode"`
		} `json:"primitives"`
	} `json:"meshes"`
	ExtensionsRequired []string `json:"extensionsRequired"`
}

// ValidateGLB checks that data is a binary glTF 2.0 file every player's
// browser can draw: self-contained (no links to other files, which would
// make every player's browser fetch them from wherever they point), within
// the size and triangle limits, and without compression the client can't
// decode. It reads the file's JSON chunk only; the client parses the rest.
func ValidateGLB(data []byte) error {
	if len(data) > MaxModelBytes {
		return i18n.New("modelTooLarge", MaxModelMegabytes)
	}
	if len(data) < glbHeader+glbChunk ||
		binary.LittleEndian.Uint32(data[0:]) != glbMagic ||
		binary.LittleEndian.Uint32(data[4:]) != 2 ||
		int(binary.LittleEndian.Uint32(data[8:])) != len(data) {
		return i18n.New("invalidModel")
	}
	chunkLength := int(binary.LittleEndian.Uint32(data[glbHeader:]))
	if binary.LittleEndian.Uint32(data[glbHeader+4:]) != glbChunkJSON ||
		chunkLength > len(data)-glbHeader-glbChunk {
		return i18n.New("invalidModel")
	}
	chunk := data[glbHeader+glbChunk : glbHeader+glbChunk+chunkLength]
	// The spec pads the chunk with spaces; some exporters use zeros.
	chunk = bytes.TrimRight(chunk, " \x00")
	var doc gltfDocument
	if err := json.Unmarshal(chunk, &doc); err != nil || !strings.HasPrefix(doc.Asset.Version, "2.") {
		return i18n.New("invalidModel")
	}
	for _, extension := range doc.ExtensionsRequired {
		if unsupportedExtensions[extension] {
			return i18n.New("modelUnsupported", extension)
		}
	}
	for _, buffer := range doc.Buffers {
		if buffer.URI != "" && !strings.HasPrefix(buffer.URI, "data:") {
			return i18n.New("modelExternalFile")
		}
	}
	for _, image := range doc.Images {
		if image.URI != "" && !strings.HasPrefix(image.URI, "data:") {
			return i18n.New("modelExternalFile")
		}
	}
	triangles := int64(0)
	for _, mesh := range doc.Meshes {
		for _, primitive := range mesh.Primitives {
			accessor, ok := primitive.Attributes["POSITION"]
			if primitive.Indices != nil {
				accessor, ok = *primitive.Indices, true
			}
			if !ok || accessor < 0 || accessor >= len(doc.Accessors) {
				return i18n.New("invalidModel")
			}
			triangles += primitiveTriangles(primitive.Mode, doc.Accessors[accessor].Count)
		}
	}
	if triangles == 0 {
		return i18n.New("modelEmpty")
	}
	if triangles > MaxModelTriangles {
		return i18n.New("modelTooManyTriangles", MaxModelTriangles)
	}
	return nil
}

// primitiveTriangles counts the triangles drawn from count vertices. Points
// and lines (modes 0 to 3) count as none.
func primitiveTriangles(mode *int, count int64) int64 {
	drawMode := glTFTriangles
	if mode != nil {
		drawMode = *mode
	}
	switch drawMode {
	case glTFTriangles:
		return count / 3
	case glTFTriangleStrip, glTFTriangleFan:
		return max(count-2, 0)
	}
	return 0
}
