package main

//import "fmt"

var qz = make(map[string]string, 99999)
var hz = make(map[string]string, 99999)

var lastPosts = make(map[string]int, 99999)

/*
type have struct {
//	user string
	seen map[string]int
}

*/

var wl = make(map[string]map[string]int, 99999)

var Database = "test.db"

var Invite = "."
//var Invite = "/"

var pageTrigger = 20
//this defines number of posts allowed per page

//var forumIp = "127.0.0.1"
var ipPort = "1.2.3.4:80"//replace with your actual IP on VPS
var basePath = "/0/uf4"

var privTitles, privTids [999]string

var topicTitles = make(map[string]string, 9999)
var pageTitles = make(map[string]string, 9999)
var linkedFrom = make(map[string]string, 9999)
var lastPages = make(map[string]int, 9999)

func init() {
//	wl["00"] = make(map[) make(map[string]int, 99999) }

	wl["Test"] = map[string]int{ "9" : 1 }
	wl["Mon"] = map[string]int{ "9" : 1 }
	wl["Anon"] = map[string]int{ "9" : 1 }
	wl["Y"] = map[string]int{ "9" : 1 }

	topicTitles["0"] = "Test Topic"
	topicTitles["1"] = "UF discussion & feature requests"
	
	pageTitles["test"] = "Test Page"
	pageTitles["0"] = "Zero"


//	linkedFrom["1"] = "en/Sandbox"
	linkedFrom["0"] = "thread/0/z"
	linkedFrom["1"] = "thread/1/z"

}


