package main

import (
	"github.com/aminkbi/microChatApp/asyncq/handler"
	"github.com/aminkbi/microChatApp/asyncq/task"
	"github.com/aminkbi/microChatApp/asyncq/util"
	"time"

	"github.com/hibiken/asynq"
)

func main() {
	util.InitLogger()
	// Redis connection
	r := asynq.RedisClientOpt{Addr: "localhost:6379"}
	client := asynq.NewClient(r)
	defer client.Close()

	// Scheduler
	scheduler := asynq.NewScheduler(r, &asynq.SchedulerOpts{
		Location: time.UTC,
	})

	// Register periodic task
	if _, err := scheduler.Register("@every 10s", task.CreateArchiveMessagesTask()); err != nil {
		util.Logger.Fatalf("could not register archive messages task: %v", err)
	}
	util.Logger.Println("Registered archive messages task to run every 10 seconds")

	if _, err := scheduler.Register("@every 10s", task.CreateReportMessagesTask()); err != nil {
		util.Logger.Fatalf("could not register report messages task: %v", err)
	}
	util.Logger.Println("Registered report messages task to run every 10 seconds")

	err := scheduler.Start()
	if err != nil {
		util.Logger.Fatal(err)
	}

	// Start a worker to process task
	mux := asynq.NewServeMux()
	mux.HandleFunc(task.TypeArchiveMessages, handler.HandleArchiveMessagesTask)
	mux.HandleFunc(task.TypeReportMessages, handler.HandleReportMessagesTask)

	worker := asynq.NewServer(r, asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"default": 1,
		},
	})

	util.Logger.Println("Starting worker to process task")
	if err = worker.Run(mux); err != nil {
		util.Logger.Fatalf("could not start worker: %v", err)
	}

}
