package handler

import (
	"context"
	"github.com/aminkbi/microChatApp/asyncq/util"
	"github.com/hibiken/asynq"
)

func HandleArchiveMessagesTask(ctx context.Context, t *asynq.Task) error {
	util.Logger.Printf("Handling report messages task %v", string(t.Payload()))

	return nil
}

func HandleReportMessagesTask(ctx context.Context, t *asynq.Task) error {
	util.Logger.Printf("Handling report messages task %v", string(t.Payload()))
	return nil
}
