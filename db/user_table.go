package db

import (
	"nero/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func CreateUserTable() {
	tableName := "Users"
	attributeDefinitions := []*dynamodb.AttributeDefinition{
		{
			AttributeName: aws.String("phone"),
			AttributeType: aws.String("N"),
		},
	}

	keySchema := []*dynamodb.KeySchemaElement{
		{
			AttributeName: aws.String("phone"),
			KeyType:       aws.String("HASH"),
		},
	}

	provisionedThroughput := &dynamodb.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(10),
		WriteCapacityUnits: aws.Int64(10),
	}

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions:  attributeDefinitions,
		KeySchema:             keySchema,
		ProvisionedThroughput: provisionedThroughput,
		TableName:             aws.String(tableName),
	}
	_, err := db.CreateTable(input)
	if err != nil {
		utils.LogError("db/db.go", err, "line-57, func-Init")
		return
	}
	utils.ColoredPrintln("Table "+tableName+" created successfully", utils.CYellow)
}

func DeleteUserTable() {
	del := &dynamodb.DeleteTableInput{
		TableName: aws.String("Users"),
	}
	db.DeleteTable(del)
	utils.ColoredPrintln("User table deleted", utils.CRed)
}
