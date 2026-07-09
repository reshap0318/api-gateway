package handlers

import (
	"errors"
	"strings"

	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/reshap0318/api-gateway/internal/proxy"
	"github.com/reshap0318/api-gateway/internal/services"
)

// Handlers holds all HTTP handlers.
type Handlers struct {
	svcs     *services.Services
	Validate *validator.Validate
	trans    ut.Translator

	// RouteManager is wired post-construction from di.Container (needed by the manual
	// cache refresh/status endpoints — see gateway_cache_handler.go).
	RouteManager *proxy.RouteManager
}

// NewHandlers creates a new Handlers instance.
func NewHandlers(svcs *services.Services, validate *validator.Validate, trans ut.Translator) *Handlers {
	return &Handlers{
		svcs:     svcs,
		Validate: validate,
		trans:    trans,
	}
}

// getErrorsMap parses validation errors into a formatted map using translator.
func (h *Handlers) getErrorsMap(err error) map[string][]string {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		errorsMap := make(map[string][]string)
		for _, e := range validationErrs {
			fieldKey := e.Field()
			humanName := toHumanReadable(fieldKey)

			msg := e.Translate(h.trans)
			msg = strings.Replace(msg, e.Field(), humanName, 1)

			errorsMap[fieldKey] = append(errorsMap[fieldKey], msg)
		}
		return errorsMap
	}
	return nil
}

// toHumanReadable converts snake_case to human-readable format.
func toHumanReadable(s string) string {
	words := strings.Split(s, "_")
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}
