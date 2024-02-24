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
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name" validate:"min=3"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=8"`
	Phone      int    `json:"phone" validate:"required,numeric,min=10"`
	Age        int    `json:"age" validate:"required,gt=12"`
	DOB        string `json:"dob"`
	Gender     string `json:"gender" validate:"required,oneof=male female"`
	ProfileUrl string `json:"profile_url"`
	Active     bool   `json:"active"`
	CreatedAt  int64  `json:"createdAt"`
	UpdatedAt  int64  `json:"-"`
	DeletedAt  int64  `json:"-"`
}

func (u *User) Create() (*User, error) {
	// initialize database
	db := db.GetDB()
	// initialize uuid to get id
	user, _ := GetByPhone(u.Phone)
	if len(user.ID) > 0 || user.Phone > 0 {
		return nil, errors.New("User alrady exists.")
	}
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
		utils.LogError("models/user.go", err, "line-49, converting user data")
		return nil, errors.New("error when try to convert user data to dynamodbattribute")
	}
	params := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String("Users"),
	}
	if _, err := db.PutItem(params); err != nil {
		utils.LogError("models/user.go", err, "line-64, inserting user data")
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
	resp, err := db.GetItem(params)
	if err != nil {
		return nil, err
	}
	var user *User
	if err := dynamodbattribute.UnmarshalMap(resp.Item, &user); err != nil {
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
