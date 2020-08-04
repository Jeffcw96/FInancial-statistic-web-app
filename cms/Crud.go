package cms

import (
	"encoding/json"
	"fmt"
	"goAgain/db"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"

	"github.com/mitchellh/mapstructure"
)

type AllExpensesOption struct {
	Expenses []ExpensesOption `json:"expenses"`
}

type ExpensesOption struct {
	Name string `json:"name"`
	Id   int32  `json:"id"`
}

//Salary struct
type Salary struct {
	Saving float64 `json:"saving"`
	Month  string  `json:"month"`
}

//Expenses struct
type Expenses struct {
	Food          float64 `json:"food,omitempty"`
	Transport     float64 `json:"transport,omitempty"`
	Entertainment float64 `json:"entertainment,omitempty"`
	Loan          float64 `json:"loan,omitempty"`
	Family        float64 `json:"family,omitempty"`
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

	getAllExpensesOption, _ := db.Client.HGetAll("expenses:option").Result()

	expensesOptionArr := []ExpensesOption{}
	for k, v := range getAllExpensesOption {
		fmt.Println("expenses option info", v)

		hm := make(map[string]interface{})

		hm["id"], _ = strconv.ParseInt(k, 10, 32)
		hm["name"] = v
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
