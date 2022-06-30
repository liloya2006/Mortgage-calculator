package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
)

type Data struct {
	InitialLoan  float64
	DownPayment  float64
	DesiredTerm  float64
	InterestRate float64
}

//
//func (b *BanksPage) Loan(initialLoan float64) error {
//	if initialLoan >= 1000000 {
//		return errors.New("The bank cannot provide such an amount")
//	}
//	b.InitialLoan = initialLoan
//	return nil
//}
//
//func (b *BanksPage) Payment(downPayment float64) error {
//	if downPayment < (b.InitialLoan * 0.2) {
//		return errors.New("Very small contribution")
//	}
//	b.DownPayment = downPayment
//	return nil
//}
//
//func (b *BanksPage) Term(desiredTerm float64) error {
//	if desiredTerm >= 120 {
//		return errors.New("Too long")
//	}
//	b.DesiredTerm = desiredTerm
//	return nil
//}

//type Calendar struct {
//	DesiredTerm    float64
//	InitialLoan    float64
//	MonthlyPayment float64
//}

type ConditionsBank struct {
	InterestRate       float64
	MaximumLoan        float64
	MinimumDownPayment float64
	LoanTerm           float64
}

func (c ConditionsBank) Condition() string {
	return fmt.Sprintf("Interest Rate: %f, Maximum Loan: %f, Minimum Down Payment: %f, Loan Term: %f", c.InterestRate, c.MaximumLoan, c.MinimumDownPayment, c.LoanTerm)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	Bank1 := ConditionsBank{0.21, 1000000.0, 0.2, 120.0}
	tmpl, err := template.ParseFiles("HomePage.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, Bank1)
}

func BankPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("BankPage.html")

	tmpl.Execute(w, nil)
}

func Save(w http.ResponseWriter, r *http.Request) {
	initialLoan := r.FormValue("initialLoan")
	downPayment := r.FormValue("downPayment")
	desiredTerm := r.FormValue("desiredTerm")
	interestRate := r.FormValue("interestRate")

	InitialLoanF, _ := strconv.ParseFloat(initialLoan, 64)
	DownPaymentF, _ := strconv.ParseFloat(downPayment, 64)
	DesiredTermF, _ := strconv.ParseFloat(desiredTerm, 64)
	InterestRateF, _ := strconv.ParseFloat(interestRate, 64)

	if initialLoan == "" || downPayment == "" || desiredTerm == "" || interestRate == "" {
		fmt.Fprintf(w, "There are free fields")
	} else {

		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889/golang")
		if err != nil {
			panic(err)
		}

		defer db.Close()

		insert, err := db.Query(fmt.Sprintf("INSERT INTO `Data` (`interestRate`, `initialLoan`, `downPayment`, `loanTerm`) VALUES('%d', '%d', '%d', '%d')", interestRate, initialLoan, downPayment, loanTerm))
		if err != nil {
			panic(err)
		}
		defer insert.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}

func Formula() {

	db, err := sql.Open("sqlite3", "calendar.db")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	rows, err := db.Query("select * from Data")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	Datas := []Data{}

	for rows.Next() {
		d := Data{}
		err := rows.Scan(&d.DesiredTerm, &d.InitialLoan, &d.DownPayment, &d.InterestRate)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	for _, d := range Datas {
		fmt.Println(d.DownPayment, d.DesiredTerm, d.InitialLoan, d.InterestRate)
		AmountBorrowed := d.InitialLoan - d.DownPayment
		x := math.Pow((1 + d.InterestRate), d.DesiredTerm)
		MounthlyPayment := (AmountBorrowed * d.InterestRate * x) / (x - 1)
		fmt.Println(MounthlyPayment)
	}
}

func HandleRequest() {
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/bank/", BankPage)
	http.HandleFunc("/save/", Save)
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}

func main() {
	Formula()
	HandleRequest()
}
