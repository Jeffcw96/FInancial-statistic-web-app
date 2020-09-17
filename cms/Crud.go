package cms

import (
	"encoding/json"
	"fmt"
	"goAgain/db"
	"io/ioutil"
	"net/http"
	"reflect"
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

	jsonFeed, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println("err in jsonFeed")
	}

	optionName := ExpensesOption{}
	json.Unmarshal([]byte(jsonFeed), &optionName)

	optionId := db.Client.Incr("expenses:ids").Val()

	hm := make(map[string]interface{})
	hm[strconv.FormatInt(optionId, 10)] = optionName.Name

	m := make(map[string]interface{})
	m["id"] = optionId
	m["name"] = optionName.Name

	db.Client.HMSet("expenses:option:"+strconv.FormatInt(optionId, 10), m)
	db.Client.HMSet("expenses:option", hm)

	response := ResponseStatus{}
	response.Status = "00"

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func ReadExpensesObject(w http.ResponseWriter, r *http.Request) {
	month, date := GenerateMonthAndDate()

	getAllExpensesOption, _ := db.Client.HGetAll("expenses:option").Result()
	getTodayExpensesValue := db.Client.HGetAll("expenses:" + month + ":" + date).Val()

	expensesOptionArr := []ExpensesOption{}
	for k, v := range getAllExpensesOption {
		fmt.Println("expenses option info", v)

		hm := make(map[string]interface{})

		hm["id"], _ = strconv.ParseInt(k, 10, 32)
		hm["name"] = v

		for expensesObject, expensesVal := range getTodayExpensesValue {
			fmt.Println("expensesVal", expensesVal)
			fmt.Println("object", expensesObject)
			if strings.ToLower(expensesObject) == strings.ToLower(v) {
				fmt.Println("FOund !!!")
				hm["currentValue"] = expensesVal
				hm["currentRemark"] = getTodayExpensesValue["R-"+strings.ToLower(expensesObject)]
			}
			fmt.Println("type of v", reflect.TypeOf(v))
			fmt.Println("string object", reflect.TypeOf(strings.ToLower(expensesObject)))
			fmt.Println("string v", reflect.TypeOf(strings.ToLower(v)))
			fmt.Println(reflect.TypeOf(expensesObject))
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

	v := mux.Vars(r)
	optionId := v["id"]

	getAllOptions, _ := db.Client.HGetAll("expenses:option").Result()

	for id, _ := range getAllOptions {
		if id == optionId {
			db.Client.HDel("expenses:option", id)
			db.Client.Del("expenses:option:" + id)
		}
	}

	fmt.Println("optionId", optionId)
	response := ResponseStatus{}
	response.Status = "00"

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func GenerateMonthAndDate() (string, string) {
	monthAndDate := time.Now().Format("2006-Jan-02")
	getMonthDateSplit := strings.Split(monthAndDate, "-")
	month := getMonthDateSplit[1]
	date := getMonthDateSplit[2]
	return month, date
}
