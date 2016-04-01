package main

type Message interface {
	GetBody() *string
	GetReceiptHandle() *string
	GetQueueUrl() *string
}

type SqsMessage struct {
	Body          *string
	ReceiptHandle *string
	QueueURL      *string
}

func (sm SqsMessage) GetBody() *string {
	return sm.Body
}

func (sm SqsMessage) GetReceiptHandle() *string {
	return sm.ReceiptHandle
}

func (sm SqsMessage) GetQueueUrl() *string {
	return sm.QueueURL
}
