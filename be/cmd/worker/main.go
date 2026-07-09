package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"

	"github.com/reshap0318/api-gateway/internal/di"
	"github.com/reshap0318/api-gateway/internal/jobs"
	pkgasynq "github.com/reshap0318/api-gateway/internal/pkg/asynq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	container, err := di.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	defer container.Close()

	redisOpt := pkgasynq.RedisOpt()

	// Server — processes jobs dequeued from Redis
	srv := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 5,
		Queues: map[string]int{
			"default": 1,
		},
	})

	mux := asynq.NewServeMux()
	j := jobs.NewJobs(container.Services)
	j.Register(mux)

	if err := srv.Start(mux); err != nil {
		log.Fatalf("Failed to start worker server: %v", err)
	}
	defer srv.Shutdown()

	// Scheduler — enqueues jobs based on cron schedule
	scheduler := asynq.NewScheduler(redisOpt, nil)

	if _, err := scheduler.Register("30 2 * * *", asynq.NewTask(jobs.TypeClearTmp, nil)); err != nil {
		log.Fatalf("Failed to register %s schedule: %v", jobs.TypeClearTmp, err)
	}

	if err := scheduler.Start(); err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}
	defer scheduler.Shutdown()

	log.Println("Worker started. Registered jobs: cleartmp (02:30)")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker...")
}
