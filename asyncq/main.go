package main

import (
	"github.com/aminkbi/microChatApp/asyncq/handler"
	"github.com/aminkbi/microChatApp/asyncq/task"
	"github.com/aminkbi/microChatApp/internal/util"
	"github.com/hibiken/asynq"
	"time"
)

func main() {
	util.InitLogger()

	err := util.ConnectMongoDB()
	if err != nil {
		util.Logger.Fatal("Can't connect to mongo: ", err)
	}
	// Redis connection
	r := asynq.RedisClientOpt{Addr: "localhost:6379"}
	client := asynq.NewClient(r)
	defer client.Close()

	// Scheduler
	scheduler := asynq.NewScheduler(r, &asynq.SchedulerOpts{
		Location: time.UTC,
	})

	// Register periodic task
	if _, err = scheduler.Register("@every 24h", task.CreateArchiveMessagesTask()); err != nil {
		util.Logger.Fatalf("could not register archive messages task: %v", err)
	}
	util.Logger.Println("Registered archive messages task to run every 10 seconds")

	if _, err = scheduler.Register("@every 24h", task.CreateReportMessagesTask()); err != nil {
		util.Logger.Fatalf("could not register report messages task: %v", err)
	}
	util.Logger.Println("Registered report messages task to run every 10 seconds")

	err = scheduler.Start()
	defer scheduler.Shutdown()
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
