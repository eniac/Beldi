package core

import (
	"crypto/sha512"
	"encoding/hex"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/dgrijalva/jwt-go"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/lithammer/shortuuid"
	"github.com/mitchellh/mapstructure"
	"time"
)

func RegisterUserWithUserId(env *beldilib.Env, firstName string, lastName string, username string, password string,
	userId string) {
	hasher := sha512.New()
	salt := shortuuid.New()
	hasher.Write([]byte(password + salt))
	passwordHash := hex.EncodeToString(hasher.Sum(nil))
	user := User{
		UserId:    userId,
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
		Password:  passwordHash,
		Salt:      salt,
	}
	beldilib.Write(env, TUser(), username, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V"): expression.Value(user),
	})
}

func RegisterUser(env *beldilib.Env, firstName string, lastName string, username string, password string) {
	RegisterUserWithUserId(env, firstName, lastName, username, password, shortuuid.New())
}

func Login(env *beldilib.Env, username string, password string) string {
	item := beldilib.Read(env, TUser(), username)
	var user User
	beldilib.CHECK(mapstructure.Decode(item, &user))
	hasher := sha512.New()
	hasher.Write([]byte(password + user.Salt))
	passwordHash := hex.EncodeToString(hasher.Sum(nil))
	if passwordHash == user.Password {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id":   user.UserId,
			"timestamp": time.Now().Format("20060102150405"),
			"TTL":       "60000",
		})
		tokenString, err := token.SignedString("secret")
		beldilib.CHECK(err)
		return tokenString
	} else {
		panic("Password not correct")
	}
}

func UploadUser(env *beldilib.Env, reqId string, username string) {
	item := beldilib.Read(env, TUser(), username)
	var user User
	beldilib.CHECK(mapstructure.Decode(item, &user))
	beldilib.AsyncInvoke(env, TComposeReview(), RPCInput{
		Function: "UploadUserId",
		Input: aws.JSONValue{
			"reqId":  reqId,
			"userId": user.UserId,
		},
	})
}
