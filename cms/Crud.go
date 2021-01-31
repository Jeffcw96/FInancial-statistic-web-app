package cms

import (
	"encoding/json"
	"fmt"
	"goAgain/db"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
)

type ReportStatistic struct {
	Jul []StatisticData `json:"jul"`
	Aug []StatisticData `json:"aug"`
}

type StatisticData struct {
	ExpensesInfo []ExpensesString `json:"expensesInfo"`
	TotalPrice   float64          `json:"totalPrice"`
	Date         string           `json:"date"`
}

type AllExpensesOption struct {
	Expenses []ExpensesOption `json:"expenses"`
}

type ExpensesOption struct {
	Name          string `json:"name"`
	Id            int32  `json:"id"`
	CurrentValue  string `json:"currentValue"`
	CurrentRemark string `json:"currentRemark"`
}

//Salary struct
type Salary struct {
	Saving float64 `json:"saving"`
	Month  string  `json:"month"`
}

//Expenses struct
type ExpensesInfo struct {
	ExpensesOption string  `json:"expensesOption,omitempty"`
	ExpensesValue  float64 `json:"expensesValue,omitempty"`
	ExpensesRemark string  `json:"expensesRemark,omitempty"`
}

//Expenses Arr struct
type ExpensesArr struct {
	AllExpenses []ExpensesInfo `json:"allExpenses"`
}

//This need to be changed
type ExpensesString struct {
	Food          string `json:"food"`
	Transport     string `json:"transport"`
	Entertainment string `json:"entertainment"`
	Loan          string `json:"loan"`
	Family        string `json:"family"`
}

type ResponseStatus struct {
	Status string `json:"status"`
}

func CreateNewExpensesObject(w http.ResponseWriter, r *http.Request) {
	fmt.Println("add expenses option")
	userId := r.Header.Get("UserId")

	jsonFeed, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println("err in jsonFeed")
	}

	optionName := ExpensesOption{}
	json.Unmarshal([]byte(jsonFeed), &optionName)

	optionId := db.Client.Incr("expenses:ids").Val()

	db.Client.HSet("expenses:"+userId+":options", strconv.FormatInt(optionId, 10), optionName.Name)

	response := ResponseStatus{}
	response.Status = "00"

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func ReadExpensesObject(w http.ResponseWriter, r *http.Request) {
	year, month, date := GenerateMonthAndDate()
	userId := r.Header.Get("UserId")
	getAllExpensesOption, _ := db.Client.HGetAll("expenses:" + userId + ":options").Result()
	getTodayExpensesValue := db.Client.HGet("expenses:"+userId+":"+year+"-"+month, date).Val()

	expensesOptionArr := []ExpensesOption{}
	for k, v := range getAllExpensesOption {
		//fmt.Println("expenses option info", v)

		hm := make(map[string]interface{})

		hm["id"], _ = strconv.ParseInt(k, 10, 32)
		hm["name"] = v

		expenses := make(map[string]interface{})
		json.Unmarshal([]byte(getTodayExpensesValue), &expenses)

		for object, val := range expenses {
			//fmt.Println("expensesVal", val)
			//fmt.Println("object", object)
			if strings.ToLower(object) == strings.ToLower(v) {
				//fmt.Println("FOund !!!")
				hm["currentValue"] = val
				hm["currentRemark"] = expenses["R-"+strings.ToLower(object)]
			}
		}
		expensesOption := ExpensesOption{}
		_ = mapstructure.Decode(hm, &expensesOption)
		fmt.Println("&expensesOption", &expensesOption)
		expensesOptionArr = append(expensesOptionArr, expensesOption)
	}

	sort.SliceStable(expensesOptionArr, func(i, j int) bool {
		return expensesOptionArr[i].Id < expensesOptionArr[j].Id
	})

	fmt.Println("expenses Option Arr", expensesOptionArr)
	allExpenses := AllExpensesOption{}
	allExpenses.Expenses = expensesOptionArr
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(expensesOptionArr)
}

func DeleteExpensesOption(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("UserId")
	v := mux.Vars(r)
	optionId := v["id"]
	db.Client.HDel("expenses:"+userId+":options", optionId)

	fmt.Println("optionId", optionId)
	response := ResponseStatus{}
	response.Status = "00"

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func GenerateMonthAndDate() (string, string, string) {
	monthAndDate := time.Now().Format("2006-01-02")
	getMonthDateSplit := strings.Split(monthAndDate, "-")
	year := getMonthDateSplit[0]
	month := getMonthDateSplit[1]
	date := getMonthDateSplit[2]

	if len(month) == 1 {
		month = "0" + month
	}

	if len(date) == 1 {
		date = "0" + date
	}

	return year, month, date
}
