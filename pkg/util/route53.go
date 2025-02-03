package util

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	log "github.com/sirupsen/logrus"
)

func InvokeRecordActionsFromVIPS(ctx context.Context, action types.ChangeAction, vips []string, domainName string) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("failed to load configuration, %v", err)
	}

	// Create a Route 53 client.
	svc := route53.NewFromConfig(cfg)

	// Specify the hosted zone ID and the domain name.
	hostedZoneID := os.Getenv("HOSTED_ZONE_ID")

	recordType := types.RRTypeA // The record type you want to check (e.g., A, CNAME, TXT).

	// List the records and check if the specific record exists.
	checkRecordName := aws.String(fmt.Sprintf("api.%s", domainName))
	lInput := &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(hostedZoneID),
		StartRecordName: checkRecordName,
		StartRecordType: recordType,
		MaxItems:        aws.Int32(1), // Limit the search to the specific record.
	}

	lResult, err := svc.ListResourceRecordSets(context.TODO(), lInput)
	if err != nil {
		return fmt.Errorf("failed to list resource record sets, %v", err)
	}

	// Check if the record exists.
	if len(lResult.ResourceRecordSets) > 0 &&
		*lResult.ResourceRecordSets[0].Name == *checkRecordName &&
		lResult.ResourceRecordSets[0].Type == recordType {

		if action == types.ChangeActionUpsert {
			log.Printf("Record already exists: %+v\n", lResult.ResourceRecordSets[0])
			return nil
		}

	} else {
		if action == types.ChangeActionDelete {
			log.Println("Record does not exist.")
		}
	}

	// Create the A record.
	input := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(hostedZoneID),
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: action,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: aws.String(fmt.Sprintf("api.%s", domainName)),
						Type: types.RRTypeA,
						TTL:  aws.Int64(300),
						ResourceRecords: []types.ResourceRecord{
							{
								Value: aws.String(vips[0]),
							},
						},
					},
				},
				{
					Action: action,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: aws.String(fmt.Sprintf("api-int.%s", domainName)),
						Type: types.RRTypeA,
						TTL:  aws.Int64(300),
						ResourceRecords: []types.ResourceRecord{
							{
								Value: aws.String(vips[0]),
							},
						},
					},
				},
				{
					Action: action,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: aws.String(fmt.Sprintf("*.apps.%s", domainName)),
						Type: types.RRTypeA,
						TTL:  aws.Int64(300),
						ResourceRecords: []types.ResourceRecord{
							{
								Value: aws.String(vips[1]),
							},
						},
					},
				},
			},
			Comment: aws.String("Creating A record for " + domainName),
		},
	}
	log.Printf("%v route53 records for %s", action, domainName)
	// Make the API call to create the record.
	result, err := svc.ChangeResourceRecordSets(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to %v A record, %v", action, err)
	}

	changeID := result.ChangeInfo.Id
	log.Infof("Waiting for DNS change to complete...")

	for {
		statusResult, err := svc.GetChange(context.TODO(), &route53.GetChangeInput{
			Id: changeID,
		})
		if err != nil {
			log.Fatalf("failed to get change status, %v", err)
		}

		log.Debugf("Change status: %s\n", statusResult.ChangeInfo.Status)

		if statusResult.ChangeInfo.Status == types.ChangeStatusInsync {
			log.Infof("DNS change is complete.")
			break
		}

		// Wait for a while before checking again.
		time.Sleep(10 * time.Second)
	}

	log.Debugf("ChangeInfo: %+v\n", result.ChangeInfo)
	return nil
}
