package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// JWKSGetKeys returns the JWKS public keys
func (h *Handlers) JWKSGetKeys(c *gin.Context) {
	jwks := h.svcs.JWKSManager.GetJWKS()

	c.JSON(http.StatusOK, jwks)
}
