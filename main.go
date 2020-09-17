package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
	"text/template"

	"goAgain/cms"
	"goAgain/db"
	"goAgain/statistic"

	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/cors"
)

// create a global variable to access html file in router
var templates *template.Template

func main() {

	templates = template.Must(template.ParseGlob("template/*.html"))
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/second", SecondFunction).Methods("GET")
	r.HandleFunc("/addSaving", AddMonthlySaving).Methods("POST")
	r.HandleFunc("/addExpenses", AddDailyExpenses).Methods("POST")
	r.HandleFunc("/addExpensesOption", cms.CreateNewExpensesObject).Methods("POST")
	r.HandleFunc("/deleteExpensesOption/{id}", cms.DeleteExpensesOption).Methods("POST")
	r.HandleFunc("/readExpensesObject", cms.ReadExpensesObject).Methods("GET")
	r.HandleFunc("/getFinancialStatistic", statistic.GetFinancialStatistic).Methods("GET")
	r.HandleFunc("/generateExpensesSummary/{month}", statistic.GenerateExpensesSummary).Methods("GET")

	corsOpts := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},                                     //Because now we are using the same port and domain, so it does not really matter
		AllowedMethods: []string{"GET", "POST", "HEAD", "OPTIONS", "PUT"}, // specify what method are allows
		AllowedHeaders: []string{"Content-Type", "Access-Control-Allow-Origin", "X-Requested-With", "Authorization"},
	})
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	http.Handle("/", r)
	handler := corsOpts.Handler(r)
	http.ListenAndServe(":8000", handler)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	//Call the index html template
	templates.ExecuteTemplate(w, "index.html", nil)
}

//AddMonthlySaving function
func AddMonthlySaving(w http.ResponseWriter, r *http.Request) {
	jsonFeed, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("read people err: ", err)
	}
	salary := cms.Salary{}
	json.Unmarshal([]byte(jsonFeed), &salary)
	fmt.Println("salary", salary)

	month, _ := cms.GenerateMonthAndDate()

	m := make(map[string]interface{})
	m["saving"] = salary.Saving
	m["month"] = month
	fmt.Println("m", m)

	db.Client.HMSet("saving:"+month, m)

	getStatus := cms.ResponseStatus{}
	getStatus.Status = "00"
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(getStatus)
}

//AddDailyExpenses function
func AddDailyExpenses(w http.ResponseWriter, r *http.Request) {

	jsonFeed, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println("nil", nil)
	}

	dailyExpenses := cms.ExpensesArr{}
	json.Unmarshal([]byte(jsonFeed), &dailyExpenses)
	fmt.Println("daily Expenses 666", dailyExpenses)
	month, date := cms.GenerateMonthAndDate()
	for _, expensesData := range dailyExpenses.AllExpenses {
		expensesInfo := cms.ExpensesInfo{}
		_ = mapstructure.Decode(expensesData, &expensesInfo)

		fmt.Println("expenses Info", expensesInfo)
		fmt.Println("type of expenses value", reflect.TypeOf(expensesData.ExpensesValue))
		db.Client.HSet("expenses:"+month+":"+date, expensesData.ExpensesOption, expensesData.ExpensesValue)
		db.Client.HSet("expenses:"+month+":"+date, "R-"+strings.ToLower(expensesData.ExpensesOption), expensesData.ExpensesRemark)
	}

	hm := make(map[string]interface{})
	hm[date] = 1
	db.Client.HMSet("expenses:"+month+":all", hm)

	getStatus := cms.ResponseStatus{}
	getStatus.Status = "00"
	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	w.WriteHeader((http.StatusOK))
	json.NewEncoder(w).Encode(getStatus)
}

func SecondFunction(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "second.html", nil)
}
