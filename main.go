package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"github.com/practice/db"
	"github.com/practice/statistic"

	"github.com/practice/cms"

	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"github.com/practice/user"
	"github.com/rs/cors"
)

// create a global variable to access html file in router
var templates *template.Template

func main() {

	templates = template.Must(template.ParseGlob("template/*.html"))
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", loginPage).Methods("GET")
	r.HandleFunc("/home", indexHandler).Methods("GET")
	r.HandleFunc("/second", SecondFunction).Methods("GET")
	r.HandleFunc("/auth/register", user.DoRegiser).Methods("POST")
	r.HandleFunc("/auth/login", user.DoLogin).Methods("POST")
	r.HandleFunc("/auth/forgotPassword", user.ForgotPassword).Methods("POST")
	r.HandleFunc("/api/addSaving", AddMonthlySaving).Methods("POST")
	r.HandleFunc("/api/addExpenses", AddDailyExpenses).Methods("POST")
	r.HandleFunc("/api/addExpensesOption", cms.CreateNewExpensesObject).Methods("POST")
	r.HandleFunc("/api/deleteExpensesOption/{id}", cms.DeleteExpensesOption).Methods("POST")
	r.HandleFunc("/api/readExpensesObject", cms.ReadExpensesObject).Methods("GET")
	r.HandleFunc("/api/getFinancialStatistic/{year}", statistic.GetFinancialStatistic).Methods("GET")
	r.HandleFunc("/api/generateExpensesSummary/{year}/{month}", statistic.GenerateExpensesSummary).Methods("GET")

	corsOpts := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},                                     //Because now we are using the same port and domain, so it does not really matter
		AllowedMethods: []string{"GET", "POST", "HEAD", "OPTIONS", "PUT"}, // specify what method are allows
		AllowedHeaders: []string{"Content-Type", "Access-Control-Allow-Origin", "X-Requested-With", "Authorization"},
	})
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	http.Handle("/", r)
	handler := corsOpts.Handler(r)

	port := os.Getenv("PORT")

	http.ListenAndServe(port, Middleware(handler))
	http.ListenAndServe(":8000", Middleware(handler))
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestUrl := r.URL.Path //return "/createUser"
		method := r.Method       //return "POST"

		//In case if only allow POST and GET request
		if method != "POST" && method != "GET" {
			//Response bad request if doesn't pass the condition
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
			return
		}

		if strings.Contains(requestUrl, "api") {
			fmt.Println("is API")
			token := r.Header.Get("Authorization")

			if token == "" {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Bad Request"))
				return
			}

			jwtWrapper := user.JwtWrapper{
				SecretKey: "jeffdevslife",
				Issuer:    "AuthService",
			}

			jwtClaims := &user.JwtClaim{}
			claims, err := jwtWrapper.ValidateToken(token)

			if claims == jwtClaims {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Bad Request"))
				return
			}

			fmt.Println("clam token err", err)

			fmt.Println("id, email", claims.Id, claims.Email)

			r.Header.Set("UserId", claims.Id)
			r.Header.Set("UserEmail", claims.Email)

		}

		next.ServeHTTP(w, r)

		fmt.Println("after api call")
	})
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	//Call the index html template
	templates.ExecuteTemplate(w, "login.html", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	//Call the index html template
	templates.ExecuteTemplate(w, "index.html", nil)
}

//AddMonthlySaving function
func AddMonthlySaving(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("UserId")
	jsonFeed, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("read people err: ", err)
	}
	salary := cms.Salary{}
	json.Unmarshal([]byte(jsonFeed), &salary)
	fmt.Println("salary", salary)

	year, month, _ := cms.GenerateMonthAndDate()
	getAllMonths := db.Client.HGetAll("expenses:" + userId + ":months").Val()
	monthLabel := ""

	for labelMonth, monthNum := range getAllMonths {
		intMonthNum, _ := strconv.ParseInt(monthNum, 10, 64)
		intMonth, _ := strconv.ParseInt(month, 10, 64)
		diff := intMonthNum - intMonth

		if diff == 0 {
			monthLabel = labelMonth
		}
	}

	m := make(map[string]interface{})
	m["saving"] = salary.Saving
	m["month"] = monthLabel
	fmt.Println("m", m)

	db.Client.HMSet("saving:"+userId+":"+year+"-"+month, m)

	getStatus := cms.ResponseStatus{}
	getStatus.Status = "00"
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(getStatus)
}

//AddDailyExpenses function
func AddDailyExpenses(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("UserId")
	jsonFeed, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println("nil", nil)
	}

	dailyExpenses := cms.ExpensesArr{}
	json.Unmarshal([]byte(jsonFeed), &dailyExpenses)
	fmt.Println("daily Expenses 666", dailyExpenses)
	year, month, date := cms.GenerateMonthAndDate()

	expensesHash := make(map[string]interface{})
	for _, expensesData := range dailyExpenses.AllExpenses {
		expensesInfo := cms.ExpensesInfo{}
		_ = mapstructure.Decode(expensesData, &expensesInfo)

		fmt.Println("expenses Info", expensesInfo)
		fmt.Println("type of expenses value", reflect.TypeOf(expensesData.ExpensesValue))
		expensesHash[expensesData.ExpensesOption] = expensesData.ExpensesValue
		expensesHash["R-"+strings.ToLower(expensesData.ExpensesOption)] = expensesData.ExpensesRemark

	}

	jsonHash, _ := json.Marshal(expensesHash)

	db.Client.HSet("expenses:"+userId+":"+year+"-"+month, date, jsonHash)

	getStatus := cms.ResponseStatus{}
	getStatus.Status = "00"
	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	w.WriteHeader((http.StatusOK))
	json.NewEncoder(w).Encode(getStatus)
}

func SecondFunction(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "second.html", nil)
}
