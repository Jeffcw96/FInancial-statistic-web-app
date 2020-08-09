package statistic

import (
	"encoding/json"
	"fmt"
	"goAgain/cms"
	"goAgain/db"
	"net/http"
	"strconv"

	"github.com/mitchellh/mapstructure"
)

func GetFinancialStatistic(w http.ResponseWriter, r *http.Request) {
	var a interface{}
	getExpensesMonths, _ := db.Client.LRange("expenses:month", 0, -1).Result()
	for ind, months := range getExpensesMonths {
		fmt.Println("ind", ind)
		fmt.Println("current month", months)
		statisticData := cms.StatisticData{}
		getExpensesDate, _ := db.Client.LRange("expenses:"+months+":all", 0, -1).Result()
		for _, date := range getExpensesDate {
			var totalDailyCost float64
			expensesValueArr := []cms.ExpensesString{}
			expensesValue := cms.ExpensesString{}
			getExpensesValue, _ := db.Client.HGetAll("expenses:" + months + ":" + date).Result()
			fmt.Println("getExpensesValue", getExpensesValue)

			for _, value := range getExpensesValue {
				expensesInFloat, _ := strconv.ParseFloat(value, 64)
				totalDailyCost += expensesInFloat
			}

			_ = mapstructure.Decode(getExpensesValue, &expensesValue)
			fmt.Println("expensesValue", expensesValue)
			expensesValueArr = append(expensesValueArr, expensesValue)
			statisticData.ExpensesInfo = expensesValueArr
			statisticData.TotalPrice = totalDailyCost
		}
		a = statisticData
	}

	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(a)
}
