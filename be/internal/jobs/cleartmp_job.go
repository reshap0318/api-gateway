package jobs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hibiken/asynq"

	"github.com/reshap0318/api-gateway/internal/helpers"
)

// HandleClearTmp processes the cleartmp job.
// Deletes files in storage/tmp older than TMP_AGE_HOURS (default 24h).
func (j *Jobs) HandleClearTmp(ctx context.Context, t *asynq.Task) error {
	dir := helpers.GetEnv("TMP_DIR", "storage/tmp")
	ageHours := helpers.GetEnvInt("TMP_AGE_HOURS", 24)
	threshold := time.Now().Add(-time.Duration(ageHours) * time.Hour)

	j.svcs.Logger.LogStart("HandleClearTmp", "Scanning %s (files older than %dh)", dir, ageHours)

	deleted, errCount := 0, 0

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			errCount++
			return nil
		}
		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			errCount++
			return nil
		}

		if info.ModTime().After(threshold) {
			return nil
		}

		if err := helpers.DeleteFile(path); err != nil {
			j.svcs.Logger.LogWarn("HandleClearTmp", "Failed to delete %s: %v", path, err)
			errCount++
			return nil
		}

		deleted++
		j.svcs.Logger.LogStep("HandleClearTmp", "Deleted: %s", path)
		return nil
	})

	if err != nil {
		j.svcs.Logger.LogEndWithError("HandleClearTmp", "Walk failed: %v", err)
		return fmt.Errorf("%s: walk %s: %w", TypeClearTmp, dir, err)
	}

	j.svcs.Logger.LogEnd("HandleClearTmp", "Done — deleted: %d, errors: %d", deleted, errCount)
	return nil
}
