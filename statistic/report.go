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
	getAllMonth := db.Client.HGetAll("expenses:month").Val()
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
		getAllSaving := db.Client.HGetAll("saving:" + monthData.Month).Val()
		fmt.Println("getAllSaving", getAllSaving)
		fmt.Println("monthData.Month", monthData.Month)
		getAllDate := db.Client.HGetAll("expenses:" + monthData.Month + ":all").Val()
		var totalMonthExpenses float64
		expensesInfo := MonthlyExpenses{}

		for date, _ := range getAllDate {
			getDailyExpense := db.Client.HGetAll("expenses:" + monthData.Month + ":" + date).Val()
			//fmt.Println("expenses:"+monthData.Month+":"+date, getDailyExpense)
			for object, value := range getDailyExpense {
				if !strings.Contains(object, "R-") {
					floatExpense, _ := strconv.ParseFloat(value, 64)
					totalMonthExpenses += floatExpense
				}
			}
		}
		fmt.Println(monthData.Month+" expenses >>", totalMonthExpenses)
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
	Vars := mux.Vars(r)
	month := strings.ToLower(Vars["month"])
	summary := Summary{}
	summaryArr := []Summary{}
	var totalValue float64

	expensesOption := db.Client.HGetAll("expenses:option").Val()
	for _, option := range expensesOption {
		summary.Expenses = strings.ToLower(option)
		summaryArr = append(summaryArr, summary)
	}

	getAllExpenses := db.Client.HGetAll("expenses:" + strings.Title(month) + ":all").Val()
	fmt.Println("getAllExpenses", getAllExpenses)

	for date, _ := range getAllExpenses {
		dailyExpenses := db.Client.HGetAll("expenses:" + strings.Title(month) + ":" + date).Val()
		for object, value := range dailyExpenses {
			for i := 0; i < len(summaryArr); i++ {
				if summaryArr[i].Expenses == object {
					fltValue, _ := strconv.ParseFloat(value, 64)
					fmt.Println("fltValue", fltValue)
					summaryArr[i].Value += fltValue
					totalValue += fltValue
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
