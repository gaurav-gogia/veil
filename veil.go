package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {
	f, err := os.OpenFile("./log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	log.SetOutput(f)

	fmt.Printf("Starting server at: http://localhost:%s\n", PORT)
	registerHandlers()
	err = http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func registerHandlers() {
	http.HandleFunc("/", index)
	http.HandleFunc("/unsafe", unsafe)
	http.HandleFunc("/waf", waf)
	http.HandleFunc("/decepticon", decepticon)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./index.html")
}

func unsafe(w http.ResponseWriter, r *http.Request) {
	atk := "ello!"

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}
		atk = r.Form.Get("atk")
	}

	_, err := fmt.Fprintf(w, atk)
	if err != nil {
		log.Fatalln(err)
	}
}

func waf(w http.ResponseWriter, r *http.Request) {
	atk := "ello!"

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}
		atk = r.Form.Get("atk")
		log.Println(r.RemoteAddr, r.RequestURI, r.UserAgent(), "|------> ", atk)

		if isSQLInjection(atk) || isXSSAttack(atk) {
			atk = ""
		}
	}

	_, err := fmt.Fprintf(w, atk)
	if err != nil {
		log.Fatalln(err)
	}
}

func decepticon(w http.ResponseWriter, r *http.Request) {
	atk := "ello!"

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}
		atk = r.Form.Get("atk")
		log.Println(r.RemoteAddr, r.RequestURI, r.UserAgent(), "|------> ", atk)

		if isSQLInjection(atk) {
			atk = "<script>alert(155157 rows affected)</script>"
		}

		if isXSSAttack(atk) {
			inside := strings.Split(atk, "alert(")[1]
			inside = strings.Split(inside, ")")[0]
			if len(inside) < 10 {
				atk = "<script>alert(" + inside + ")</script>"
			} else {
				atk = "<script>alert(1)</script>"
			}
		}
	}

	_, err := fmt.Fprintf(w, atk)
	if err != nil {
		log.Fatalln(err)
	}
}

func isSQLInjection(input string) bool {
	sqlInjectionPattern := `(\b(?:drop|delete|truncate|insert|update|union|alter|select|create|exec)\b|\b(?:;|--))`
	match, _ := regexp.MatchString(sqlInjectionPattern, input)
	return match
}

func isXSSAttack(input string) bool {
	xssPattern := `(<\s*script\s*>)|(<\s*\/\s*script\s*>)`
	match, _ := regexp.MatchString(xssPattern, strings.ToLower(input))
	return match
}
