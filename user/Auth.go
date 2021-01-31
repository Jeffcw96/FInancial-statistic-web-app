package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"goAgain/db"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/practice/cms"
	"golang.org/x/crypto/bcrypt"
)

type JwtWrapper struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

// JwtClaim adds email as a claim to the token
type JwtClaim struct {
	Email string
	Id    string
	jwt.StandardClaims
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func DoRegiser(w http.ResponseWriter, r *http.Request) {

	jsonFeed, _ := ioutil.ReadAll(r.Body)
	user := User{}

	json.Unmarshal([]byte(jsonFeed), &user)
	hm := make(map[string]interface{})
	bytePassword := []byte(user.Password)

	id := db.Client.Incr("user:ids").Val()

	hm["id"] = id
	hm["email"] = user.Email
	hm["password"] = hashAndSalt(bytePassword)

	fmt.Println("hm", hm)

	db.Client.HSet("user:all", user.Email, id)
	db.Client.HMSet("user:"+strconv.FormatInt(id, 10)+":data", hm)

	monthMap := make(map[string]interface{})
	monthMap["Jan"] = "01"
	monthMap["Feb"] = "02"
	monthMap["Mar"] = "03"
	monthMap["Apr"] = "04"
	monthMap["May"] = "05"
	monthMap["Jun"] = "06"
	monthMap["Jul"] = "07"
	monthMap["Aug"] = "08"
	monthMap["Sep"] = "09"
	monthMap["Oct"] = "10"
	monthMap["Nov"] = "11"
	monthMap["Dec"] = "12"

	defaultOption := make(map[string]interface{})
	defaultOption["1"] = "Food"
	defaultOption["2"] = "Transport"
	defaultOption["3"] = "Entertainment"
	defaultOption["4"] = "Family"
	defaultOption["5"] = "Loan"
	defaultOption["6"] = "Other"

	db.Client.IncrBy("expenses:"+strconv.FormatInt(id, 10)+":ids", 7)

	db.Client.HMSet("expenses:"+strconv.FormatInt(id, 10)+":months", monthMap)
	db.Client.HMSet("expenses:"+strconv.FormatInt(id, 10)+":options", defaultOption)

	response := cms.ResponseStatus{}
	response.Status = "00"

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func DoLogin(w http.ResponseWriter, r *http.Request) {

	jsonFeed, _ := ioutil.ReadAll(r.Body)
	user := User{}

	json.Unmarshal([]byte(jsonFeed), &user)

	existEmail := db.Client.HExists("user:all", user.Email).Val()

	if !existEmail {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("User not found")
		return
	}

	userId := db.Client.HGet("user:all", user.Email).Val()
	userData := db.Client.HGetAll("user:" + userId + ":data").Val()
	password := userData["password"]

	isMatch := comparePasswords(password, []byte(user.Password))
	fmt.Println("isMatch", isMatch)

	if !isMatch {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Invalid Credential")
		return
	}

	jwtWrapper := JwtWrapper{
		SecretKey:       "jeffdevslife",
		Issuer:          "AuthService",
		ExpirationHours: 24,
	}

	token, err := jwtWrapper.GenerateToken(userData["email"], userId)

	if err != nil {
		fmt.Println("error in generating token", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)
}

func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}
func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func (j *JwtWrapper) GenerateToken(email, id string) (signedToken string, err error) {
	claims := &JwtClaim{
		Email: email,
		Id:    id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(j.ExpirationHours)).Unix(),
			Issuer:    j.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return
	}

	return
}

//ValidateToken validates the jwt token
func (j *JwtWrapper) ValidateToken(signedToken string) (claims *JwtClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*JwtClaim)
	if !ok {
		err = errors.New("Couldn't parse claims")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("JWT is expired")
		return
	}

	return

}
