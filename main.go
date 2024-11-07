package main

import (
	"log"
	"os/exec"

	"github.com/dhruv-assessment/app-tier/service"
)

func main() {
	for {
		noOfMsgReqQueue, err := service.GetNoOfMessagesInRequestQueue()
		if err != nil {
			log.Printf("Unable to get no. of messages in request queue: %v", err)
			break
		}
		if noOfMsgReqQueue == 0 {
			log.Println("No. of messages are zero in req so terminating this machine")
			break
		} else {
			log.Println("Reading from Request queue")
			messageReq, messageID, messageReceiptID, err := service.ReadMessagesInRequestQueue()
			if err != nil {
				log.Printf("Unable to read message from request queue: %v", err)
				break
			}
			log.Printf("Message->%v --- Message ID->%v\n", messageReq, messageID)
			log.Println("Deleting the message from request queue")
			err = service.DeleteMessageFromSQS(messageReceiptID)
			if err != nil {
				log.Printf("Unable to delete the message from request queue")
				break
			}
			log.Println("Successfully deleted the message from request queue")
			log.Println("Downloading the image from S3")
			if err := service.DownloadFromS3(messageReq); err != nil {
				log.Printf("Download from S3 failed: %v\n", err)
				break
			}
			log.Printf("Running the model on %v\n", messageReq)
			prediction, err := exec.Command("python3", "./model/face_recognition.py", messageReq).Output()
			if err != nil {
				log.Printf("Not able to run the model: %v\n", err)
				break
			}
			log.Printf("Successfully determined the %v as %v\n", messageReq, string(prediction))
			log.Println("Sending message to response queue")
			_, err = service.SendMessageToSQS(string(prediction), messageID)
			if err != nil {
				log.Printf("Unable to send message to response queue: %v", err)
				break
			}
			log.Println("Successfully sent the prediction to response queue")
		}
	}
	log.Println("Getting this instance's ID")
	instanceID, err := service.GetInstanceID()
	if err != nil {
		log.Printf("Unable to get instance-id: %v\n", err)
		return
	}
	log.Printf("ID of this instance is %v\n", instanceID)
	log.Println("Terminating this instance")
	err = service.TerminateEC2(instanceID)
	if err != nil {
		log.Printf("Unable to terminate ec2: %v", err)
		return
	}
	log.Println("Successfully terminated this instance")
}

// Things to do in AMU: install go
// app_tier code should exist
// #!/bin/bash
// cd /root
// ./app_tier
