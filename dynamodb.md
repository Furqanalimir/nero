# [Step-1]
    download dynamodb docker image to run locally
    $ docker pull amazon/dynamodb-local

# [Step-2]
    create docker compose file (docker-compode.yaml)
        version: '3.7'
        services:
          dynamodb:
            image: amazon/dynamodb-local
            container_name: my-dynamodb
            hostname: dynamodb
            restart: always
            ports:
              - "8000:8000"
            environment:
              - SERVICE_NAME=my-dynamodb
              - SERVICE_TAGS=dev
              - SERVICE_PORT_NUMBER=8000
              - AWS_REGION=us-west-2
              - AWS_ACCESS_KEY_ID=dummy
              - AWS_SECRET_ACCESS_KEY=dummy

#   Start above image
        $ docker compose up -d

#   [Basics]
#   Start with ERD (Enterprise Relation Diagram)    
#   Define your access patterns
#   Design your primary and secondary indexes
#   https://www.youtube.com/watch?v=DIQVJqiSUkE (Data modeling with Amazon DynamoDB)

# [Step-3]
#   Setup db connection in golang

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
        input := &dynamodb.CreateTableInput{
            AttributeDefinitions: []*dynamodb.AttributeDefinition{
                {
                    AttributeName: aws.String("id"),
                    AttributeType: aws.String("S"),
                },
                {
                    AttributeName: aws.String("name"),
                    AttributeType: aws.String("S"),
                },
            },
            KeySchema: []*dynamodb.KeySchemaElement{
                {
                    AttributeName: aws.String("id"),
                    KeyType:       aws.String("HASH"),
                },
                {
                    AttributeName: aws.String("name"),
                    KeyType:       aws.String("RANGE"),
                },
            },
            ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
                ReadCapacityUnits:  aws.Int64(10),
                WriteCapacityUnits: aws.Int64(10),
            },
            TableName: aws.String("Users"),
        }

        result, err := db.CreateTable(input)
        if err != nil {
            utils.LogError("db/db.go", err, "line-57, func-Init")
            return
        }
        // call clist all tables function to list tables created
        listAllTables()
        utils.ColoredPrintln(result.GoString(), utils.CGreen)
        // res, err := db.CreateTable("Users")
        utils.ColoredPrintln("Database connected!", utils.CGreen)
    }

    func GetDB() *dynamodb.DynamoDB {
        return db
    }

    func listAllTables() {
        // create the input configuration instance
        input := &dynamodb.ListTablesInput{}

        utils.ColoredPrintln("Tables: ", utils.CYellow)

        for {
            // Get the list of tables
            result, err := db.ListTables(input)
            if err != nil {
                if aerr, ok := err.(awserr.Error); ok {
                    switch aerr.Code() {
                    case dynamodb.ErrCodeInternalServerError:
                        // fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
                        utils.LogError("db/db.go, line-85", err, aerr.Error())
                    default:
                        utils.LogError("db/db.go", err, "line-86, func-listAllTables")
                    }
                } else {
                    // Print the error, cast err to awserr.Error to get the Code and
                    // Message from an error.
                    // fmt.Println(err.Error())
                    utils.LogError("db/db.go", err, "line-92, func-listAllTables")
                }
                return
            }

            for _, n := range result.TableNames {
                fmt.Println(*n)
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

