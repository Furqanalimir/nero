package db

import (
	"nero/config"
	"nero/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	dynamodbsession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var db *dynamodb.DynamoDB

func Init() {
	sess, err := dynamodbsession.NewSession(&aws.Config{
		Region:      aws.String(config.StrictEnvVars("DB_REGION")),
		Credentials: credentials.NewEnvCredentials(),
		Endpoint:    aws.String(config.StrictEnvVars("DB_ENDPOINT")),
		DisableSSL:  aws.Bool(true),
	})
	if err != nil {
		panic(err)
	}
	db = dynamodb.New(sess)
	if config.EnvVars("CREATE_TABLES") == "true" {
		// Delete tables
		DeleteUserTable()
		DeleteOrderTable()

		// Create tables
		CreateUserTable()
		CreateOrderTable()
	}
	listAllTables()
	utils.ColoredPrintln("Database connected!", utils.CBlue)
}

func GetDB() *dynamodb.DynamoDB {
	return db
}

func listAllTables() {
	// create the input configuration instance
	input := &dynamodb.ListTablesInput{}
	for {
		// Get the list of tables
		result, err := db.ListTables(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeInternalServerError:
					// fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
					utils.LogError("db/db.go "+dynamodb.ErrCodeInternalServerError, err, "line-54")
				default:
					utils.LogError("db/db.go", err, "line-55, func-listAllTables")
				}
			} else {
				utils.LogError("db/db.go", err, "line-59, func-listAllTables")
			}
			return
		}

		utils.ColoredPrintln("Tables Name: ", utils.CBlue)
		for _, n := range result.TableNames {
			utils.ColoredPrintln(""+*n, utils.CBlue)
		}

		// assign the last read tablename as the start for our next call to the ListTables function
		// the maximum number of table names returned in a call is 100 (default), which requires us to make
		// multiple calls to the ListTables function to retrieve all table names
		input.ExclusiveStartTableName = result.LastEvaluatedTableName

		if result.LastEvaluatedTableName == nil {
			break
		}
	}

}
