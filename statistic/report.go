package statistic

import (
	"encoding/json"
	"fmt"
	"goAgain/db"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type Month struct {
	Month string `json:"month"`
	Value int64  `json:"value"`
}

type MonthlyExpenses struct {
	Month         string  `json:"month"`
	TotalExpenses float64 `json:"totalExpenses"`
	Saving        float64 `json:"saving"`
}

type MonthlyReport struct {
	Report []MonthlyExpenses `json:"report"`
}

type Summary struct {
	Expenses string  `json:"expenses"`
	Value    float64 `json:"value"`
}

func GetFinancialStatistic(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("UserId")
	Vars := mux.Vars(r)
	year := strings.ToLower(Vars["year"])
	getAllMonth := db.Client.HGetAll("expenses:" + userId + ":months").Val()
	fmt.Println(getAllMonth)
	sortedMonth := []Month{}

	for month, value := range getAllMonth {
		monthStruct := Month{}
		monthStruct.Month = month

		intMonth, _ := strconv.ParseInt(value, 10, 64)
		monthStruct.Value = intMonth

		sortedMonth = append(sortedMonth, monthStruct)
	}
	sort.Slice(sortedMonth, func(i, j int) bool { return sortedMonth[i].Value < sortedMonth[j].Value })
	fmt.Println("sortedMonth", sortedMonth)

	monthlyReport := MonthlyReport{}
	expensesInfoArr := []MonthlyExpenses{}
	for _, monthData := range sortedMonth {
		month := strconv.FormatInt(monthData.Value, 10)
		if len(month) == 1 {
			month = "0" + month
		}
		getAllSaving := db.Client.HGetAll("saving:" + userId + ":" + year + "-" + month).Val()
		fmt.Println("getAllSaving", getAllSaving)
		fmt.Println("monthData.Month", monthData.Month)

		getMonthlyExpenses := db.Client.HGetAll("expenses:" + userId + ":" + year + "-" + month).Val()
		//fmt.Println("getMonthlyExpenses", getMonthlyExpenses)
		var totalMonthExpenses float64
		expensesInfo := MonthlyExpenses{}

		for _, expenses := range getMonthlyExpenses {
			expensesMap := make(map[string]interface{})
			json.Unmarshal([]byte(expenses), &expensesMap)

			for object, value := range expensesMap {
				if !strings.Contains(object, "R-") {
					stringExp := fmt.Sprintf("%v", value)
					floatExpense, _ := strconv.ParseFloat(stringExp, 64)
					totalMonthExpenses += floatExpense
				}
			}
		}
		//fmt.Println(monthData.Month+" expenses >>", totalMonthExpenses)
		expensesInfo.Month = monthData.Month
		expensesInfo.TotalExpenses = totalMonthExpenses
		formatedSaving, _ := strconv.ParseFloat(getAllSaving["saving"], 64)
		expensesInfo.Saving = formatedSaving
		expensesInfoArr = append(expensesInfoArr, expensesInfo)
	}
	monthlyReport.Report = expensesInfoArr

	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(monthlyReport)
}

func GenerateExpensesSummary(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("UserId")
	Vars := mux.Vars(r)
	month := strings.ToLower(Vars["month"])
	year := strings.ToLower(Vars["year"])
	summary := Summary{}
	summaryArr := []Summary{}
	var totalValue float64

	monthLabel := strings.Title(month)
	month = db.Client.HGet("expenses:"+userId+":months", monthLabel).Val()

	if len(month) == 1 {
		month = "0" + month
	}

	expensesOption := db.Client.HGetAll("expenses:" + userId + ":options").Val()
	for _, option := range expensesOption {
		summary.Expenses = strings.ToLower(option)
		summaryArr = append(summaryArr, summary)
	}

	getAllExpenses := db.Client.HGetAll("expenses:" + userId + ":" + year + "-" + month).Val()
	fmt.Println("getAllExpenses", getAllExpenses)

	for _, expenses := range getAllExpenses {
		jsonExpenses := make(map[string]float64)
		json.Unmarshal([]byte(expenses), &jsonExpenses)

		for object, value := range jsonExpenses {
			for i := 0; i < len(summaryArr); i++ {
				if summaryArr[i].Expenses == object {
					fmt.Println("value", value)
					// fltValue, _ := strconv.ParseFloat(value, 64)
					// fmt.Println("fltValue", fltValue)
					summaryArr[i].Value += value
					totalValue += value
				}
			}
		}
	}
	fmt.Println("total Value", totalValue)
	for i := 0; i < len(summaryArr); i++ {
		expensesPercentage := (summaryArr[i].Value / totalValue) * 100
		expensesPercentage = math.Floor(expensesPercentage*100) / 100

		// the format string change the float to string in 2 decimal places, 2f means 2 decimal
		summaryArr[i].Expenses = summaryArr[i].Expenses + " ~" + fmt.Sprintf("%.2f", expensesPercentage) + "%"
		fmt.Println("expensesPercentage", expensesPercentage)
	}

	fmt.Println("summaryArr", summaryArr)

	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(summaryArr)

}
