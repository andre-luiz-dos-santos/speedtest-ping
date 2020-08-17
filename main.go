package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var (
	webBindAddr    string
	allowFrom      IPNetList
	telegramChatID string
	telegramToken  string
	ixcDSN         string

	ixcDB *sql.DB
)

func main() {
	var err error

	flag.StringVar(&webBindAddr, "web-bind", ":8080", "web bind address")
	flag.StringVar(&ixcDSN, "ixc-dsn", "", "IXC database")
	flag.StringVar(&telegramToken, "telegram-token", "", "Telegram token")
	flag.StringVar(&telegramChatID, "telegram-chat-id", "", "Telegram chat ID")
	allowFromStr := flag.String("allow-from", "", "limit permitted IPs")
	flag.Parse()

	err = allowFrom.ParseArg(*allowFromStr)
	if err != nil {
		log.Fatal(err)
	}

	ixcDB, err = sql.Open("mysql", ixcDSN)
	if err != nil {
		log.Fatal(err)
	}
	ixcDB.SetMaxOpenConns(5)

	err = ixcDB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/ping", pingHandler)

	log.Fatal(http.ListenAndServe(webBindAddr, allowNetworkWrapper(http.DefaultServeMux)))
}
