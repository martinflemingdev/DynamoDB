package main

import (
    "context"
    "fmt"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials/stscreds"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
    "github.com/aws/aws-sdk-go-v2/service/sts"
)

// NewDynamoClient creates a new instance of the DynamoDB client
func NewDynamoClient(ctx context.Context, roleARN string) (*dynamodb.Client, error) {
    // Load the Shared AWS Configuration (~/.aws/config)
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        return nil, fmt.Errorf("unable to load SDK config, %w", err)
    }

    // Assume an IAM role using the STS service
    stsClient := sts.NewFromConfig(cfg)
    creds := stscreds.NewAssumeRoleProvider(stsClient, roleARN)

    // Create a DynamoDB client with the assumed role
    dynamoClient := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
        o.Credentials = aws.NewCredentialsCache(creds)
    })

    return dynamoClient, nil
}

// QueryDynamo queries a DynamoDB table using a partition key
func QueryDynamo(ctx context.Context, client *dynamodb.Client, tableName, searchString string) (*dynamodb.QueryOutput, error) {
    input := &dynamodb.QueryInput{
        TableName: aws.String(tableName),
        KeyConditionExpression: aws.String("searchString = :value"),
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":value": &types.AttributeValueMemberS{Value: searchString},
        },
    }

    result, err := client.Query(ctx, input)
    if err != nil {
        return nil, fmt.Errorf("failed to query items, %w", err)
    }

    return result, nil
}


func main() {
    // Example usage
    ctx := context.TODO()
    roleARN := "arn:aws:iam::123456789012:role/example-role"

    client, err := NewDynamoClient(ctx, roleARN)
    if err != nil {
        fmt.Printf("Failed to create DynamoDB client: %s\n", err)
        return
    }

    tableName := "YourDynamoDBTable"
    partitionKey := "YourPartitionKey"
    partitionValue := "YourPartitionValue"

    result, err := QueryDynamo(ctx, client, tableName, partitionKey, partitionValue)
    if err != nil {
        fmt.Printf("Failed to query DynamoDB: %s\n", err)
        return
    }

    fmt.Printf("Query result: %v\n", result)
	
}
