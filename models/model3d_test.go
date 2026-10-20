package models

import (
	"encoding/binary"
	"errors"
	"strings"
	"testing"

	"ships3d/i18n"
)

// glb wraps a glTF JSON document in a binary glTF file.
func glb(document string) []byte {
	for len(document)%4 != 0 {
		document += " "
	}
	data := make([]byte, 20, 20+len(document))
	binary.LittleEndian.PutUint32(data[0:], glbMagic)
	binary.LittleEndian.PutUint32(data[4:], 2)
	binary.LittleEndian.PutUint32(data[8:], uint32(20+len(document)))
	binary.LittleEndian.PutUint32(data[12:], uint32(len(document)))
	binary.LittleEndian.PutUint32(data[16:], glbChunkJSON)
	return append(data, document...)
}

const triangleDocument = `{"asset":{"version":"2.0"},"buffers":[{"byteLength":36}],
"accessors":[{"count":3}],"meshes":[{"primitives":[{"attributes":{"POSITION":0}}]}]}`

func messageKey(t *testing.T, err error) string {
	t.Helper()
	var message i18n.Message
	if !errors.As(err, &message) {
		t.Fatalf("error %v is not a message", err)
	}
	return message.Key
}

func TestValidateGLBAcceptsATriangle(t *testing.T) {
	if err := ValidateGLB(glb(triangleDocument)); err != nil {
		t.Fatalf("valid model rejected: %v", err)
	}
}

func TestValidateGLBRejects(t *testing.T) {
	broken := glb(triangleDocument)
	binary.LittleEndian.PutUint32(broken[8:], uint32(len(broken)+4))
	cases := map[string]struct {
		data []byte
		key  string
	}{
		"not a glb":          {[]byte("{\"asset\":{\"version\":\"2.0\"}}"), "invalidModel"},
		"wrong length":       {broken, "invalidModel"},
		"glTF 1":             {glb(strings.Replace(triangleDocument, `"2.0"`, `"1.0"`, 1)), "invalidModel"},
		"bad JSON":           {glb(`{"asset":`), "invalidModel"},
		"no meshes":          {glb(`{"asset":{"version":"2.0"}}`), "modelEmpty"},
		"bad accessor":       {glb(strings.Replace(triangleDocument, `"POSITION":0`, `"POSITION":3`, 1)), "invalidModel"},
		"external file":      {glb(strings.Replace(triangleDocument, `{"byteLength":36}`, `{"uri":"https://example.com/a.bin","byteLength":36}`, 1)), "modelExternalFile"},
		"external image":     {glb(strings.Replace(triangleDocument, `"buffers"`, `"images":[{"uri":"tex.png"}],"buffers"`, 1)), "modelExternalFile"},
		"draco":              {glb(strings.Replace(triangleDocument, `"buffers"`, `"extensionsRequired":["KHR_draco_mesh_compression"],"buffers"`, 1)), "modelUnsupported"},
		"too many triangles": {glb(strings.Replace(triangleDocument, `"count":3`, `"count":300003`, 1)), "modelTooManyTriangles"},
		"too large":          {make([]byte, MaxModelBytes+1), "modelTooLarge"},
	}
	for name, test := range cases {
		err := ValidateGLB(test.data)
		if err == nil {
			t.Errorf("%s: accepted", name)
			continue
		}
		if key := messageKey(t, err); key != test.key {
			t.Errorf("%s: got %s, want %s", name, key, test.key)
		}
	}
}

func TestValidateGLBAcceptsEmbeddedData(t *testing.T) {
	document := strings.Replace(triangleDocument, `{"byteLength":36}`, `{"uri":"data:application/octet-stream;base64,AAAA","byteLength":36}`, 1)
	if err := ValidateGLB(glb(document)); err != nil {
		t.Fatalf("embedded buffer rejected: %v", err)
	}
}

func TestPrimitiveTriangles(t *testing.T) {
	strip, fan, lines := 5, 6, 1
	if got := primitiveTriangles(nil, 9); got != 3 {
		t.Errorf("triangles = %d, want 3", got)
	}
	if got := primitiveTriangles(&strip, 6); got != 4 {
		t.Errorf("strip = %d, want 4", got)
	}
	if got := primitiveTriangles(&fan, 1); got != 0 {
		t.Errorf("short fan = %d, want 0", got)
	}
	if got := primitiveTriangles(&lines, 10); got != 0 {
		t.Errorf("lines = %d, want 0", got)
	}
}
