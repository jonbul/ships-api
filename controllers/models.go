package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"net/http"
	"path"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"

	dataaccess "ships3d/dataAccess"
	"ships3d/i18n"
	"ships3d/models"
)

// maxModelBodyBytes caps an upload: the largest model allowed, plus one
// byte so a larger one is seen as too large rather than cut short.
const maxModelBodyBytes = models.MaxModelBytes + 1

const glbContentType = "model/gltf-binary"

func (s *Server) registerModelRoutes(router *gin.Engine) {
	router.POST("/models", s.requireUser, s.uploadModel)
	// Public and unauthenticated: every player loads the models of the
	// ships they see.
	router.GET("/models/:id", s.getModel)
}

// uploadModel stores a .glb file sent as the request body. ?name= is the
// file's name, kept to tell models apart.
func (s *Server) uploadModel(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			errorResponse(c, http.StatusRequestEntityTooLarge, i18n.New("modelTooLarge", models.MaxModelMegabytes))
			return
		}
		errorResponse(c, http.StatusBadRequest, i18n.New("invalidBody"))
		return
	}
	if err := models.ValidateGLB(data); err != nil {
		var message i18n.Message
		if !errors.As(err, &message) {
			message = i18n.New("invalidModel")
		}
		errorResponse(c, http.StatusBadRequest, message)
		return
	}
	hash := sha256.Sum256(data)
	model, err := s.store.InsertModel(&models.Model3d{
		UserId:      contextUser(c).IdAsString(),
		Name:        modelName(c.Query("name")),
		Hash:        hex.EncodeToString(hash[:]),
		Data:        data,
		DateCreated: time.Now().UnixMilli(),
	})
	if errors.Is(err, dataaccess.ErrTooManyModels) {
		errorResponse(c, http.StatusBadRequest, i18n.New("tooManyModels", models.MaxModelsPerUser))
		return
	}
	if err != nil {
		log.Println("uploadModel:", err)
		errorResponse(c, http.StatusInternalServerError, i18n.New("saveModelFailed"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "model": model})
}

func (s *Server) getModel(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusNotFound, i18n.New("modelNotFound"))
		return
	}
	data, err := s.store.GetModelData(id)
	if errors.Is(err, dataaccess.ErrNotFound) {
		errorResponse(c, http.StatusNotFound, i18n.New("modelNotFound"))
		return
	}
	if err != nil {
		log.Println("getModel:", err)
		errorResponse(c, http.StatusInternalServerError, i18n.New("modelNotFound"))
		return
	}
	// A model never changes once uploaded (a new file gets a new id).
	c.Header("Cache-Control", "public, max-age=31536000, immutable")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Data(http.StatusOK, glbContentType, data)
}

// modelName turns an uploaded file's name into a model's: no directories or
// .glb extension, trimmed, at most MaxModelNameLength characters.
func modelName(name string) string {
	name = path.Base(strings.ReplaceAll(name, "\\", "/"))
	name = strings.TrimSpace(strings.TrimSuffix(name, path.Ext(name)))
	if name == "." || name == "/" || !utf8.ValidString(name) {
		return ""
	}
	if runes := []rune(name); len(runes) > models.MaxModelNameLength {
		name = strings.TrimSpace(string(runes[:models.MaxModelNameLength]))
	}
	return name
}
