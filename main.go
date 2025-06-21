package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sts"
)

type Response struct {
	UnattachedVolumes []string `json:"unattached_volumes"`
	DeletedSanpshots  []string `json:"deleted_sanpshots"`
}

func handler(ctx context.Context) (Response, error) {
	sessi := session.Must(session.NewSession())
	ec2Client := ec2.New(sessi)
	snsClient := sns.New(sessi)

	snsTopicArn := os.Getenv("SNS_TOPIC_ARN")

	var unattached_volumes []string
	var deleted_sanpshots []string

	// describe all volumes
	volumeOutput, err := ec2Client.DescribeVolumes(nil)
	if err != nil {
		log.Fatalf("Unable to list volumes : %v", err)
	}

	for _, vol := range volumeOutput.Volumes {
		if len(vol.Attachments) == 0 {
			volumeID := aws.StringValue(vol.VolumeId)
			unattached_volumes = append(unattached_volumes, volumeID)

			// Sending Notification
			message := fmt.Sprintf("Unattached EBS Volumes detected : %s in %s", volumeID, aws.StringValue(vol.AvailabilityZone))
			_, err := snsClient.Publish(&sns.PublishInput{
				Message:  aws.String(message),
				TopicArn: aws.String(snsTopicArn),
				Subject:  aws.String("Unattached EBS Volume Found"),
			})
			if err != nil {
				log.Printf("Failed to send SNS notification for volume %s : %v", volumeID, err)
			}

			// Finding and Deleting snapshots for unattached volume
			snapOutput, err := ec2Client.DescribeSnapshots(&ec2.DescribeSnapshotsInput{
				Filters: []*ec2.Filter{
					{
						Name:   aws.String("volume-id"),
						Values: []*string{aws.String(volumeID)},
					},
				},
				OwnerIds: []*string{aws.String(getAccountID(sessi))},
			})
			if err != nil {
				log.Fatalf("Failed to describe sanpshot for volume %s : %v", volumeID, err)
				continue
			}

			for _, snap := range snapOutput.Snapshots {
				snapID := aws.StringValue(snap.SnapshotId)
				_, err := ec2Client.DeleteSnapshot(&ec2.DeleteSnapshotInput{
					SnapshotId: aws.String(snapID),
				})
				if err != nil {
					log.Printf("Failed to delete snapshot %s : %v", snapID, err)
				} else {
					deleted_sanpshots = append(deleted_sanpshots, snapID)
				}
			}
		}
	}

	return Response{
		UnattachedVolumes: unattached_volumes,
		DeletedSanpshots:  deleted_sanpshots,
	}, nil
}

func getAccountID(sess *session.Session) string {
	stsClient := sts.New(sess)
	resp, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatalf("Unable to get AWS account ID: %v", err)
	}
	return *resp.Account
}

func main() {
	lambda.Start(handler)
}
