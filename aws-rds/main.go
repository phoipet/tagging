// Created by Phoipet <ifrusman@gmail.com>

package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

func main() {
	rdsAddTag()
}

func rdsAddTag() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := rds.New(sess)
	projectName := ""

	result, _ := svc.DescribeDBInstances(nil)
	for x := 0; x < len(result.DBInstances); x++ {
		name := *result.DBInstances[x].DBInstanceIdentifier
		arnName := *result.DBInstances[x].DBInstanceArn
		fmt.Println(name, arnName)
		matched, err := regexp.MatchString(os.Args[1], name)
		if matched {
			projectDummy, _ := regexp.MatchString("dummy", name)

			// condition for project name
			if projectDummy {
				projectName = "Dummy"
			} else {
				projectName = "=Unknown"
			}

			input := &rds.AddTagsToResourceInput{

				ResourceName: aws.String(arnName),
				Tags: []*rds.Tag{
					{
						Key:   aws.String("AWS"),
						Value: aws.String("rds"),
					},
					{
						Key:   aws.String("Stage"),
						Value: aws.String(os.Args[2]),
					},
					{
						Key:   aws.String("Product"),
						Value: aws.String(projectName),
					},
				},
			}

			result, err := svc.AddTagsToResource(input)
			fmt.Println(name, projectName, arnName)
			fmt.Println(result)
			if err != nil {
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case rds.ErrCodeDBInstanceNotFoundFault:
						fmt.Println(rds.ErrCodeDBInstanceNotFoundFault, aerr.Error())
					default:
						fmt.Println(aerr.Error())
					}
				} else {
					fmt.Println(err.Error())
				}
				return
			}

		}
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case rds.ErrCodeDBInstanceNotFoundFault:
					fmt.Println(rds.ErrCodeDBInstanceNotFoundFault, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				fmt.Println(err.Error())
			}
			return
		}
	}
}