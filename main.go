package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
)

type BanksPage struct {
	InitialLoan float64
	DownPayment float64
	DesiredTerm float64
}

func (b *BanksPage) Loan(initialLoan float64) error {
	if initialLoan > 1000000 {
		return errors.New("The bank cannot provide such an amount")
	}
	b.InitialLoan = initialLoan
	return nil
}

func (b *BanksPage) Payment(downPayment float64) error {
	if downPayment < (b.InitialLoan * 0.2) {
		return errors.New("Very small contribution")
	}
	b.DownPayment = downPayment
	return nil
}

func (b *BanksPage) Term(desiredTerm float64) error {
	if desiredTerm > 120 {
		return errors.New("Too long")
	}
	b.DesiredTerm = desiredTerm
	return nil
}

type ConditionsBank struct {
	InterestRate       float64
	MaximumLoan        float64
	MinimumDownPayment float64
	LoanTerm           float64
}

func (c ConditionsBank) Condition() string {
	return fmt.Sprintf("Interest Rate: %d, Maximum Loan: %d, Minimum Down Payment: %d, Loan Term: %d", c.InterestRate, c.MaximumLoan, c.MinimumDownPayment, c.LoanTerm)
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
	//tmpl, _ := template.ParseFiles("BankPage.html")
}

func Save(w http.ResponseWriter, r *http.Request) {
	interestRate := r.FormValue("interestRate")
	initialLoan := r.FormValue("initialLoan")
	downPayment := r.FormValue("downPayment")
	loanTerm := r.FormValue("loanTerm")

	if interestRate == "" || initialLoan == "" || downPayment == "" || loanTerm == "" {
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

func HandleRequest() {
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/bank/", BankPage)
	http.HandleFunc("/save/", Save)
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}

func main() {
	HandleRequest()
}
