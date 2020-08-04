package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	"goAgain/cms"
	"goAgain/db"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// create a global variable to access html file in router
var templates *template.Template

func main() {

	templates = template.Must(template.ParseGlob("template/*.html"))
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/addSaving", AddMonthlySaving).Methods("POST")
	r.HandleFunc("/addExpenses", AddDailyExpenses).Methods("POST")
	r.HandleFunc("/addExpensesOption", cms.CreateNewExpensesObject).Methods("POST")
	r.HandleFunc("/readExpensesObject", cms.ReadExpensesObject).Methods("GET")

	corsOpts := cors.New(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"GET", "POST", "HEAD", "OPTIONS", "PUT"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"Content-Type", "Access-Control-Allow-Origin"},
		OptionsPassthrough: true,
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

	monthAndDate := time.Now().Format("2006-Jan-02")
	getSalaryMonth := strings.Split(monthAndDate, "-")
	fmt.Println(getSalaryMonth[1])
	m := make(map[string]interface{})
	m["saving"] = salary.Saving
	m["month"] = getSalaryMonth[1]
	fmt.Println("m", m)

	db.Client.HMSet("saving:"+getSalaryMonth[1], m)

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

	dailyExpenses := cms.Expenses{}
	json.Unmarshal([]byte(jsonFeed), &dailyExpenses)
	fmt.Println("daily Expenses 666", dailyExpenses)
	totalDailyCost := dailyExpenses.Food + dailyExpenses.Entertainment + dailyExpenses.Transport

	fmt.Println("total dailly cost >>", totalDailyCost)
	monthAndDate := time.Now().Format("2006-Jan-02")
	getMonthDateSplit := strings.Split(monthAndDate, "-")
	month := getMonthDateSplit[1]
	date := getMonthDateSplit[2]

	hm := make(map[string]interface{})

	if dailyExpenses.Food != 0.0 {
		hm["food"] = dailyExpenses.Food
	}

	if dailyExpenses.Entertainment != 0.0 {
		hm["entertainment"] = dailyExpenses.Entertainment
	}

	if dailyExpenses.Transport != 0.0 {
		hm["transport"] = dailyExpenses.Transport
	}

	if dailyExpenses.Loan != 0.0 {
		hm["loan"] = dailyExpenses.Loan
	}

	if dailyExpenses.Family != 0.0 {
		hm["family"] = dailyExpenses.Family
	}

	db.Client.HMSet("expenses:"+month+":"+date, hm)

	getStatus := cms.ResponseStatus{}
	getStatus.Status = "00"
	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	w.WriteHeader((http.StatusOK))
	json.NewEncoder(w).Encode(getStatus)
}
