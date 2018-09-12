// Created by Phoipet <ifrusman@gmail.com> 

package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
)

func main() {
	ebAddTag()
}

func ebAddTag() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := elasticbeanstalk.New(sess)
	projectName := ""

	result, _ := svc.DescribeEnvironments(nil)
	for x := 0; x < len(result.Environments); x++ {
		appName := *result.Environments[x].ApplicationName
		envName := *result.Environments[x].EnvironmentName
		arnName := *result.Environments[x].EnvironmentArn

		matched, _ := regexp.MatchString(os.Args[1], envName)
		if matched {

			projectTest, _ := regexp.MatchString("test", appName)

			// condition for project name
			if projectTest {
				projectName = "Test Project"
			} else {
				projectName = "unknown"
			}

			input := &elasticbeanstalk.UpdateTagsForResourceInput{
				ResourceArn: aws.String(arnName),
				TagsToAdd: []*elasticbeanstalk.Tag{
					{
						Key:   aws.String("AWS"),
						Value: aws.String("ec2"),
					},
					{
						Key:   aws.String("Stage"),
						Value: aws.String(os.Args[2]),
					},
					{
						Key:   aws.String("Project"),
						Value: aws.String(projectName),
					},
				},
			}

			result, _ := svc.UpdateTagsForResource(input)
			fmt.Println(projectName, appName, envName)
			fmt.Println(result)
		}
	}
}