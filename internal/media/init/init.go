package main

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/eniac/Beldi/internal/media/core"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/lithammer/shortuuid"
	"io/ioutil"
	"os"
)

var services = []string{"CastInfo", "ComposeReview", "Frontend", "MovieId", "MovieInfo", "MovieReview", "Page",
	"Plot", "Rating", "ReviewStorage", "Text", "UniqueId", "User", "UserReview"}

func tables(baseline bool) {
	if baseline {
		for _, service := range services {
			tablename := fmt.Sprintf("b%s", service)
			for ; ; {
				beldilib.CreateBaselineTable(tablename)
				if beldilib.WaitUntilActive(tablename) {
					break
				}
			}
		}
	} else {
		for _, service := range services {
			for ; ; {
				beldilib.CreateLambdaTables(service)
				if beldilib.WaitUntilAllActive([]string{service, fmt.Sprintf("%s-collector", service), fmt.Sprintf("%s-log", service)}) {
					break
				}
			}
		}
	}
}

func delete_tables(baseline bool) {
	if baseline {
		for _, service := range services {
			beldilib.DeleteTable(fmt.Sprintf("b%s", service))
			beldilib.WaitUntilDeleted(fmt.Sprintf("b%s", service))
		}
	} else {
		for _, service := range services {
			beldilib.DeleteLambdaTables(service)
			beldilib.WaitUntilAllDeleted([]string{service, fmt.Sprintf("%s-collector", service), fmt.Sprintf("%s-log", service)})
		}
	}
}

func user(baseline bool) {
	for i := 0; i < 1000; i++ {
		userId := fmt.Sprintf("user%d", i)
		username := fmt.Sprintf("username_%d", i)
		password := fmt.Sprintf("password_%d", i)
		hasher := sha512.New()
		salt := shortuuid.New()
		hasher.Write([]byte(password + salt))
		passwordHash := hex.EncodeToString(hasher.Sum(nil))
		user := core.User{
			UserId:    userId,
			FirstName: "firstname",
			LastName:  "lastname",
			Username:  username,
			Password:  passwordHash,
			Salt:      salt,
		}
		beldilib.Populate("User", username, user, baseline)
	}
}

func movie(baseline bool, file string) {
	data, err := ioutil.ReadFile(file)
	beldilib.CHECK(err)
	var movies []core.MovieInfo
	err = json.Unmarshal(data, &movies)
	beldilib.CHECK(err)
	for _, movie := range movies {
		beldilib.Populate("MovieInfo", movie.MovieId, movie, baseline)
		beldilib.Populate("Plot", movie.MovieId, aws.JSONValue{"plotId": movie.MovieId, "plot": "plot"}, baseline)
		beldilib.Populate("MovieId", movie.Title, aws.JSONValue{"movieId": movie.MovieId, "title": movie.Title}, baseline)
	}
}

func populate(baseline bool, file string) {
	user(baseline)
	movie(baseline, file)
}

func main() {
	option := os.Args[1]
	baseline := os.Args[2] == "baseline"
	if option == "create" {
		tables(baseline)
	} else if option == "populate" {
		populate(baseline, os.Args[3])
	} else if option == "clean" {
		delete_tables(baseline)
	}
}
