package alexa

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"time"
	"os"
	"sync"
	"fmt"
	"encoding/json"
	"github.com/markmester/rokuTVControl/rokuAPI"
)

//PollQueue: Module for polling AWS SQS for alexa events and calling requested function e.g. powering on device pr
//launching an application
func PollQueue(wg *sync.WaitGroup, queue_name string, timeout int64) {
	defer wg.Done()

	fmt.Println(">>> Starting SQS polling...")

	// Initialize a session that the SDK will use to load configuration,
	// credentials, and region from the shared config file. (~/.aws/config).
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create a SQS service client.
	svc := sqs.New(sess)

	// Need to convert the queue name into a URL. Make the GetQueueUrl
	// API call to retrieve the URL. This is needed for receiving messages
	// from the queue.
	resultURL, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queue_name),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == sqs.ErrCodeQueueDoesNotExist {
			exitErrorf("Unable to find queue %q.", queue_name)
		}
		exitErrorf("Unable to queue %q, %v.", queue_name, err)
	}

	// Receive messages from the SQS queue with long polling enabled.
	for {
		time.Sleep(500 * time.Millisecond)
		result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl: resultURL.QueueUrl,
			AttributeNames: aws.StringSlice([]string{
				"SentTimestamp",
			}),
			MaxNumberOfMessages: aws.Int64(1),
			WaitTimeSeconds: aws.Int64(timeout),
			VisibilityTimeout:   aws.Int64(0),
		})

		if err != nil {
			exitErrorf("Unable to receive message from queue %q, %v.", queue_name, err)
		}

		fmt.Printf("Received %d messages.\n", len(result.Messages))

		type Message struct {
			Data map[string]string
			Command string
		}
		var m Message

		if len(result.Messages) > 0 {
			for i := 0; i < len(result.Messages); i++ {
				msg := result.Messages[i]

				// decode body
				body := *msg.Body
				println(body)
				json.Unmarshal([]byte(body), &m)

				// make commanded request
				command := m.Command
				if command == "launch_app" {
					app := m.Data["app"]
					println(">>> launching ", app )
					rokuAPI.LaunchAppEndpoint(app)
				} else if command == "power" {
					println(">>> powering roku device")
					rokuAPI.PowerEndpoint()
				} else {
					println(">>> unrecognized command ")
				}

				// delete message
				_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:	resultURL.QueueUrl,
					ReceiptHandle: msg.ReceiptHandle,
				})

				if err != nil {
					fmt.Println("Delete Error", err)
					return
				}
			}
		}
	}
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}