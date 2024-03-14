package models

import (
	"errors"
	"nero/db"
	"nero/utils"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         string `form:"user_id,omitempty"`
	Name       string `form:"name" validate:"min=3"`
	Email      string `form:"email" validate:"required,email"`
	ProflieUrl string `form:"profile" binding:"required"`
	Password   string `form:"password" validate:"required,min=8"`
	Phone      int    `form:"phone" validate:"required,numeric,min=10"`
	Age        int    `form:"age" validate:"required,gt=12"`
	Gender     string `form:"gender" validate:"required,oneof=male female"`
	Active     bool   `form:"active"`
	CreatedAt  int64
	UpdatedAt  int64
	DeletedAt  int64
}

func (u *User) Create() (*User, error) {
	// initialize database
	db := db.GetDB()
	pass, err := u.hashPassword()
	if err != nil {
		return nil, errors.New("Could not hash password, something went wrong")
	}
	id := uuid.New()
	u.ID = id.String()
	u.Password = pass
	u.Active = true
	u.CreatedAt = time.Now().Unix()
	u.UpdatedAt = time.Now().Unix()
	u.DeletedAt = time.Now().Unix()
	item, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		utils.LogError("models/user.go", err, "line-52, converting user data")
		return nil, errors.New("error when try to convert user data to dynamodbattribute")
	}
	params := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String("Users"),
	}
	if _, err := db.PutItem(params); err != nil {
		utils.LogError("models/user.go", err, "line-60, inserting user data")
		return nil, errors.New("error when inserting user to data to database")
	}
	return u, nil

}

func GetByPhone(phone int) (*User, error) {
	db := db.GetDB()

	params := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"phone": {
				N: aws.String(strconv.Itoa(phone)),
			},
		},
		TableName:      aws.String("Users"),
		ConsistentRead: aws.Bool(true),
	}

	result, err := db.GetItem(params)
	if err != nil {
		return nil, err
	}

	var user *User
	if err := dynamodbattribute.UnmarshalMap(result.Item, &user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) Validate() error {
	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			msg := err.Field() + err.Tag()
			utils.ColoredPrintln(msg, utils.CRed)
		}
	}
	return err
}

func (u *User) hashPassword() (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	return string(bytes), err
}

func ComparePassword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateUserObj(c *gin.Context) *User {
	phone, _ := strconv.Atoi(c.Request.FormValue("phone"))
	age, _ := strconv.Atoi(c.Request.FormValue("age"))
	active, _ := strconv.ParseBool(c.Request.FormValue("name"))

	user := &User{
		Name:     c.Request.FormValue("name"),
		Email:    c.Request.FormValue("email"),
		Password: c.Request.FormValue("password"),
		Phone:    phone,
		Age:      age,
		Gender:   c.Request.FormValue("gender"),
		Active:   active,
	}
	return user
}
