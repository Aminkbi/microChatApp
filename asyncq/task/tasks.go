package task

import (
	"github.com/hibiken/asynq"
)

// A list of task types.
const (
	TypeArchiveMessages = "message:archive"
	TypeReportMessages  = "message:report"
)

func CreateArchiveMessagesTask() *asynq.Task {
	// Add any payload data if needed
	return asynq.NewTask(TypeArchiveMessages, []byte("archive"))
}

func CreateReportMessagesTask() *asynq.Task {

	return asynq.NewTask(TypeReportMessages, []byte("report"))
}
