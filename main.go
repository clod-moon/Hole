package main

import (
	"Hole/crawler"
	"github.com/clod-moon/goconf"
	"strconv"
	"log"
)


func init() {
	conf := iniconf.InitConfig("./conf.ini")

	crawler.Username = conf.GetValue("mysql", "username")

	crawler.Password = conf.GetValue("mysql", "password")
	crawler.DbName = conf.GetValue("mysql", "dbName")
	crawler.Host = conf.GetValue("mysql", "host")
	mySqlPort, err := strconv.Atoi(conf.GetValue("mysql", "port"))
	if err != nil {
		log.Fatal("mysql port conf err!")
	}
	crawler.Port = mySqlPort

	crawler.Cookie = conf.GetValue("session", "cookie")
	crawler.Accept = conf.GetValue("session", "Accept")
	crawler.Coding = conf.GetValue("session", "Accept-Encoding")
	crawler.Language = conf.GetValue("session", "Accept-Language")
	crawler.Control = conf.GetValue("session", "Cache-Control")
	crawler.Connection = conf.GetValue("session", "Connection")
	crawler.CNDVHost = conf.GetValue("session", "Host")
	crawler.Agent = conf.GetValue("session", "User-Agent")
}

func main() {

	crawler.StartCrawler()

}
