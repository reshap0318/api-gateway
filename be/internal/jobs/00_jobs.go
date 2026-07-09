package jobs

import (
	"github.com/hibiken/asynq"

	"github.com/reshap0318/api-gateway/internal/services"
)

// Task type constants — used by both scheduler and server mux.
const (
	TypeClearTmp = "cleartmp"
)

// Jobs holds all background job handlers.
// Follows the same single-struct pattern as Handlers.
type Jobs struct {
	svcs *services.Services
}

// NewJobs creates a new Jobs instance.
func NewJobs(svcs *services.Services) *Jobs {
	return &Jobs{svcs: svcs}
}

// Register wires all job handlers to the asynq ServeMux.
func (j *Jobs) Register(mux *asynq.ServeMux) {
	mux.HandleFunc(TypeClearTmp, j.HandleClearTmp)
}
