package main

import (
	"code.google.com/p/gorilla/mux"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"log"
	"strconv"
	"strings"
	"time"
	"flag"
	//	"code.google.com/p/gorilla/sessions"
	"code.google.com/p/gorilla/securecookie"
	"github.com/boltdb/bolt"
)

//var reip = "(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$"

var Mastercode = "1e70ed81c45a494f835f45a21ce9ccf0cbdab50b24ff7713beadadbbc39fececaa95e2d1c133b766eee0da93f6a97e727f8a562a15dcfd594ce5eff4a84c1be6"

var rescript = "SRC|src|<script|</script>|<SCRIPT|</SCRIPT>|eval[(]"

//<script[^>]*>[^<]*</script>"

//<script (.|\n)*>(.|\n)*?</script>"
//s/<script[^>]*>.*?<\/script>"//igs;"
// "/<script\ .*?<\/.*?script>/i"

var resc = regexp.MustCompile(rescript)

//var rei = regexp.MustCompile(reip)


var sw = make(map[string]string, 4)

var buck = "tracks"

//var store = sessions.NewCookieStore([]byte("hhjjfuur784844ffjhhdhdgllgd++test__secrettt"))

//var secret = []byte("topppp-secccret98484848mkljfsjiuwu432kljlkusadasj;d")

var hashKey = []byte("C0%%*H$duulwaejpo2p38Y*<8g12kj3j;22432,42;jd;sa;d;al;k'k**hJ8U^t")
var blockKey = []byte("#-l76TGdlk;la;li3B*8VQ43")

//var maxAge = 0
var secure = securecookie.New(hashKey, blockKey)
var tids = [...]string{ "test", "1", "0"  }

var posts = ""
var upid = "XXX"
var cnt = 20
var last, old string

var localRun = "0"

var db *bolt.DB

func init() {
// Initialize db.

	db, err := bolt.Open(Database, 0644, &bolt.Options{Timeout: 1 * time.Second})
	db.Update(func(tx *bolt.Tx) error {
		_, err1 := tx.CreateBucketIfNotExists([]byte("tracks"))
		_, err1 = tx.CreateBucketIfNotExists([]byte("flats"))
		_, err1 = tx.CreateBucketIfNotExists([]byte("watch"))
		_, err1 = tx.CreateBucketIfNotExists([]byte("users"))
		if err1 != nil {
			return fmt.Errorf("creating bucket <?>: %s", err1)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	db.Close()

	flag.Parse()
	if flag.Arg(0) == "0" { localRun = "1" }
}


func Log(v ...interface{}) {
	fmt.Println(v)
}

func main() {

	mainLoop()
}

func handlerRoot(w http.ResponseWriter, r *http.Request) {

	titles := ""
	for key, val := range topicTitles {
		if val == "" {
			continue
		}
		if strings.Contains(key, "7") || strings.Contains(key, "9") {
			continue
		}
		titles += "<div class=\"r\" id=\"i\"><div id=\"t\"><a href=\"/t/" +
			key + "/1\"><strong>" + val + "</strong></a></div></div><br>"
	}
	fmt.Fprintf(w, "<html><head>"+
		"<link rel=\"stylesheet\" href=\"/s/1.css\" />"+
		"</head><body><p class=\"b\"><strong>List Of Titled Topics</strong></p><br>"+titles+
		"</form></html>")
	return
}

func handlerFront(w http.ResponseWriter, r *http.Request) {

	titles := ""
	for key, val := range pageTitles {
		if val == "" {
			continue
		}
		if strings.Contains(key, "7") || strings.Contains(key, "9") {
			continue
		}
		titles += "<div class=\"r\" id=\"i\"><div id=\"t\"><a href=\"/page/" +
			key + "\"><strong>" + val + "</strong></a></div></div><br>"
	}
	fmt.Fprintf(w, "<html><head>"+
		"<link rel=\"stylesheet\" href=\"/s/1.css\" />"+
		"</head><body><p class=\"b\"><strong>List Of All Pages</strong></p><br>"+titles+
		"</form></html>")
	return
}

func handlerWatchlist(w http.ResponseWriter, r *http.Request) {
	user := "Y"
	cookie, err := r.Cookie("data")
	if err == nil {
		value := make(map[string]string)
		if err = secure.Decode("data", cookie.Value, &value); err == nil {
			user = value["user"]
			fmt.Println("user from Cookie in Watch() == ", user)
		}
	}
	if user == "Y" {
		http.Redirect(w, r, "/+", http.StatusFound)
		return
	}
	titles := ""
	for _, tid := range tids {
		if wl[user][tid] >= lastPosts[tid] { continue }
		left := strconv.Itoa(wl[user][tid] + 1)
		delta := strconv.Itoa(lastPosts[tid] - wl[user][tid])
		right := strconv.Itoa(lastPosts[tid])
//		fmt.Println("User: "+user+" had not seen in "+tid+" post # "+left)
		titles += "<div class=\"r\" id=\"i\"><div id=\"t\">"+
			"[ Topic : /"+tid+" : <em>"+delta+
			"</em> new post(s) ] <a href=\"/t/"+tid+
			"/"+left+"#"+left+"\">1st Unread Post</a>  ||  <a href=\"/uf/thread/"+
			tid + "/"+right+"#"+right+"\"> Last Post</a></div></div><br>"
	}
//	}
	fmt.Fprintf(w, "<html><head>"+
		"<link rel=\"stylesheet\" href=\"/s/1.css\" />"+
		"</head><body><p class=\"b\">List Of Unread Posts by: <em>"+user+"</em></p><br>"+titles+
		"<br><br><p>&nbsp;<a href=\"/+\">Log me in/out</a></p></body></html>")
	return
}

func handlerBrokenTrack(w http.ResponseWriter, r *http.Request) {

	thread := mux.Vars(r)["thr"]
	http.Redirect(w, r, "/page/"+thread, http.StatusFound)
	return
}

func handlerMegaBrokenTrack(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "/page/0", http.StatusFound)
	return
}

func handlerTrackShow(w http.ResponseWriter, r *http.Request) {

	by := make([]byte, 22)
	rand.Read(by)
	challenge := hex.EncodeToString(by)
//	hsh := r.FormValue("hsh")
//	chal := []byte(r.FormValue("ch"))
	//	ref := r.Referer()
	user := ""
	tid := mux.Vars(r)["thr"]
//	tidstrip := thread
//	post := mux.Vars(r)["post"]
//	postint, _ := strconv.Atoi(post)
//	pageint := postint/20 + 1
//	fmt.Prairentln("page= ", page)
//	if pageint > lastPages[thread] {
//	}
//	fmt.Println("post= ", post)
//	fmt.Println("thread= ", tidstrip)
//	page := strconv.Itoa(pageint)
//	tid := thread
	fmt.Println("tid= ", tid)
//	from := linkedFrom[page]


//	topic := getStaticPrivTopic(page)
//	if topic == "e" || len(topic) == 0 {
//		topic = getPrivTopic(tid)
//	}


	topic := getTopic(tid)


//	fmt.Println(topic)
/*	title := topicTitles[tidstrip]
	if title == "" {
		title = "Untitled Topic"
	}
	ttl := title + " # " + page
	title += " [ page № " + page + " ]" */
	cookie, err := r.Cookie("data")
	if err == nil {
		value := make(map[string]string)
		if err = secure.Decode("data", cookie.Value, &value); err == nil {
			user = value["user"]
			fmt.Println("user from Cookie == ", user)
		}
	}
//	fmt.Println("user == ", user)
	status := " ~ "+user+" ~"
	if user == "Anon" || user == "" {
		status = "<em>.......</em>"
		user = "Anon"
	}
	if strings.Contains(tid, "$") && user == "Anon" {
		http.Redirect(w, r, "/t/0", http.StatusFound)
	}
/*
	if tidstrip == "B" {
		uflag := 0
		if user == "Moon" || user == "Test" {
			uflag = 1
		}
		if uflag == 0 {
			http.Redirect(w, r, "/t/0/z", http.StatusFound)
			return
		}
	}
	if tidstrip == "X1" {
		uflag := 0
		if user == "Moon" || user == "Test" {
			uflag = 1
		}
		if uflag == 0 {
			http.Redirect(w, r, "/t/0", http.StatusFound)
			return
		}
	}
	fmt.Println(user+" is watching post #"+strconv.Itoa(wl[user][tidstrip]))
	fmt.Println("pageint=", pageint)
*/
//	replybox := ""
	replybox := "<textarea name=\"txt\" id=\"rb\" cols=\"80\" rows=\"8\"></textarea>" +
	"<br><br><input type=\"submit\" value=\"Submit\"><br>"
	if user == "Anon" || user == "" { replybox = "" }
	fmt.Fprint(w, "<html><head>"+//<title>"+ttl+"</title>"+
		"<link rel=\"stylesheet\" href=\"/s/1.css\" />"+
		"<script src=\"/s/Whirlpool.min.js\"></script>\n"+
		"<script type=\"text/javascript\">\n"+
		"<!--\n"+
		"function pack(form){\n"+
		"if (form.sig.value != '')\n"+
		"{   var aa = Wh(Wh(form.sig.value).toLowerCase()).toLowerCase();\n"+
		" form.sig.value = '';\n"+
		"form.hsh.value = Wh(Wh(Wh(form.ch.value).toLowerCase() + aa).toLowerCase()).toLowerCase();\n"+
		"aa = ''; }\n"+
		"}\n"+
		"//-->\n"+
		"</script>\n"+
		"</head><body>"+//<div id=\"x\" class=\"b\" style=\"border: 2px solid #5D478B;\">"+title+
//		"</div>"+
//		<br><div class=\"b\"><a href=\""+prev+"\">&laquo;&laquo;</a>&nbsp;&nbsp;&nbsp;<a href=\"/"+from+
//		"\">Exit</a>&nbsp;&nbsp;&nbsp;<a href=\"/en/ForumTopics\">Topics</a>&nbsp;&nbsp;&nbsp;<a href=\""+next+"\">&raquo;&raquo;</a></div>"+
		"<br>"+topic+
//		"<div id=\"l\" class=\"b\"><a href=\""+prev+"\">&laquo;&laquo;</a>"+
//		"&nbsp;[ page № "+page+" ]&nbsp;<a href=\""+next+"\">&raquo;&raquo;</a></div>"+
		//	"<p class=\"h\"><em><a>"+title+"</a></em></p>"+topic+
		"<br><form action=\"/submit\" method=\"POST\" onsubmit=\"return pack(this);\">"+
		replybox+
		"<input type=\"hidden\" name=\"ch\" value=\""+challenge+"\">"+
		"<input type=\"hidden\" name=\"hsh\" value=\"\">"+
		"<input type=\"hidden\" name=\"user\" value=\""+user+"\">"+
		"<input type=\"hidden\" name=\"tid\" value=\""+tid+"\"></form>"+
		"<div class=\"b\"><a href=\"/t/0\">Home</a>&nbsp;&nbsp;[ "+status+
		"&nbsp;&nbsp;] <a href=\"/list\">List</a></div><p>&nbsp;<a href=\"/door\">Log me in/out</a></p></body></html>")

	return
}


func handlerTrackSubmit(w http.ResponseWriter, r *http.Request) {
	txt := r.FormValue("txt")
	tid := r.FormValue("tid")
//	hsh := r.FormValue("hsh")
	mlen, _ := strconv.Atoi(r.FormValue("mlen"))
	editnum := r.FormValue("num")
//	numint, _ :=strconv.Atoi(num)
	owner := r.FormValue("own")
//	chal := []byte(r.FormValue("ch"))
	
	
	if len(txt) < 1 {
		http.Redirect(w, r, "/t/"+tid, http.StatusFound)
		return
	}
	
	if txt[0] == '$' {
		bytetxt := []byte(txt)
		bytetxt = bytetxt[1:]
		txt = "<em>©</em><br><br>" + string(bytetxt)
	}
	user := ""
	if cookie, err := r.Cookie("data"); err == nil {
		value := make(map[string]string)
		if err = secure.Decode("data", cookie.Value, &value); err == nil {
			user = value["user"]
			fmt.Println("PubSubmit(): user from Cookie= ", user)
			if user == "" {
				http.Redirect(w, r, "/t/0", http.StatusFound)
				return
			}
		}
	}
	/*
	for k, v := range Phashes {
		kk := []byte(hex.EncodeToString(whirl(chal)) + v)
		kkk := hex.EncodeToString(whirl([]byte(hex.EncodeToString(whirl(kk)))))
		if kkk != hsh {
			continue
		}
		user = Users[k]
		break
	}
	if user != val {
		user = val
	}
	*/
	fmt.Println("owner from form=", owner)
	//	raddr := strings.Split(r.RemoteAddr, ":")[0]

	if user == "" {
		user = "Anon"
	}
	if user == "DrEvil" { goto Adm1 }
	if user == "9" { goto Adm1 }
	if user != owner && owner != "" {
		fmt.Fprint(w, "<html><head>"+
			"<link rel=\"stylesheet\" href=\"/s/1.css\" />"+
			"</head><body><div class=\"r\"><em>&nbsp;&nbsp;&nbsp;&nbsp;Only Post's Owner Can Edit It...</em></div></body></html>")
		return
	}
	if user == "Anon" || owner == "Anon" {
		fmt.Fprint(w, "<html><head>"+
			"<link rel=\"stylesheet\" href=\"/s/1.css\" />"+
			"</head><body><div class=\"r\"><em>&nbsp;&nbsp;&nbsp;&nbsp;Anonymous' Posts Are Immutable...</em></div></body></html>")
		return
	}
Adm1:
	meta := ""
	num := ""
//	if user == "9" || user == "DrEvil" { user = owner }

	db, err := bolt.Open(Database, 0644, &bolt.Options{Timeout: 1 * time.Second})
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("watch"))
		nu := string(b.Get([]byte(tid)))
	if nu == "" { nu = "0" }
	numb, err4 := strconv.Atoi(nu)
	if err4 != nil {
		fmt.Println(err)
	}
	num = strconv.Itoa(numb + 1)
	fmt.Println("editnum= ", editnum)
//	if editnum != "" { num = editnum }
	fmt.Printf("The 'num' is: %s\n", num)
	return nil
	})
	if mlen != 0 {
		created := time.Now().Format(time.RFC822)[:15]
		meta = "<div class=\"f\"><a id=\"u\">"+user+"</a><a>&nbsp;@ "+created+"&nbsp;&nbsp;&nbsp;post #"+num+
		"</a><a id=\"e\" href=\"/editt/"+tid+"/"+num+"\">Edit</a></div><hr>"
	}
	if mlen == 0 {
//		num = strconv.Itoa(lastPosts[tid])
		created := time.Now().Format(time.RFC822)[:15]
		meta = "<div class=\"f\"><a id=\"u\">"+user+"</a><a>&nbsp;@ "+created+"&nbsp;&nbsp;&nbsp;post #"+num+
		"</a><a id=\"e\" href=\"/editt/"+tid+"/"+num+"\">Edit</a></div><hr>"
	}
	if resc.MatchString(txt) { txt = "<em>Javacript Injection attempt !</em>" }
	txt = meta + txt

	

	err = db.Update(func(tx *bolt.Tx) error {
		bucket2 := tx.Bucket([]byte("watch"))
		if editnum == "" {
			err = bucket2.Put([]byte(tid), []byte(num))
		}
		bucket1 := tx.Bucket([]byte("tracks"))

		err = bucket1.Put([]byte(tid+"_"+num), []byte(txt))

		if err != nil {
			return err
		}
		return nil
	})
	db.Close()
	if err != nil {
		log.Fatal(err)
	}
////	if resc.MatchString(txt) { txt = "<em>Javacript Injection attempt !</em>" }
//	craftStaticPrivTopic(tid)
//	page := numint/20 + 1
//	fmt.Println("calculated page=", page)
	http.Redirect(w, r, "/t/"+tid+"#"+num, http.StatusFound)
	return
}


func handlerTrackSaveEdited(w http.ResponseWriter, r *http.Request) {
	txt := r.FormValue("txt")
	tid := r.FormValue("tid")
//	hsh := r.FormValue("hsh")
	mlen, _ := strconv.Atoi(r.FormValue("mlen"))
	editnum := r.FormValue("num")
//	numint, _ :=strconv.Atoi(num)
	owner := r.FormValue("own")
//	chal := []byte(r.FormValue("ch"))
	
	
	if len(txt) < 1 {
		http.Redirect(w, r, "/t/"+tid, http.StatusFound)
		return
	}
	
	if txt[0] == '$' {
		bytetxt := []byte(txt)
		bytetxt = bytetxt[1:]
		txt = "<em>©</em><br><br>" + string(bytetxt)
	}
	user := ""
	if cookie, err := r.Cookie("data"); err == nil {
		value := make(map[string]string)
		if err = secure.Decode("data", cookie.Value, &value); err == nil {
			user = value["user"]
			fmt.Println("PubSubmit(): user from Cookie= ", user)
			if user == "" {
				http.Redirect(w, r, "/t/0", http.StatusFound)
				return
			}
		}
	}
	/*
	for k, v := range Phashes {
		kk := []byte(hex.EncodeToString(whirl(chal)) + v)
		kkk := hex.EncodeToString(whirl([]byte(hex.EncodeToString(whirl(kk)))))
		if kkk != hsh {
			continue
		}
		user = Users[k]
		break
	}
	if user != val {
		user = val
	}
	*/
	fmt.Println("owner from form=", owner)
	//	raddr := strings.Split(r.RemoteAddr, ":")[0]

	if user == "" {
		user = "Anon"
	}
	if user == "DrEvil" { goto Adm4 }
	if user == "9" { goto Adm4 }
	if user != owner && owner != "" {
		fmt.Fprint(w, "<html><head>"+
			"<link rel=\"stylesheet\" href=\"/s/1.css\" />"+
			"</head><body><div class=\"r\"><em>&nbsp;&nbsp;&nbsp;&nbsp;Only Post's Owner Can Edit It...</em></div></body></html>")
		return
	}
	if user == "Anon" || owner == "Anon" {
		fmt.Fprint(w, "<html><head>"+
			"<link rel=\"stylesheet\" href=\"/s/1.css\" />"+
			"</head><body><div class=\"r\"><em>&nbsp;&nbsp;&nbsp;&nbsp;Anonymous' Posts Are Immutable...</em></div></body></html>")
		return
	}
Adm4:
	meta := ""
	num := ""
	if user == "9" || user == "DrEvil" { user = owner }

	db, err := bolt.Open(Database, 0644, &bolt.Options{Timeout: 1 * time.Second})
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("watch"))
		nu := string(b.Get([]byte(tid)))
	if nu == "" { nu = "0" }
	numb, err4 := strconv.Atoi(nu)
	if err4 != nil {
		fmt.Println(err)
	}
	num = strconv.Itoa(numb + 1)
	fmt.Println("editnum= ", editnum)
	num = editnum
	fmt.Printf("The 'num' is: %s\n", num)
	return nil
	})
	if mlen != 0 {
		created := time.Now().Format(time.RFC822)[:15]
		meta = "<div class=\"f\"><a id=\"u\">"+user+"</a><a>&nbsp;@ "+created+"&nbsp;&nbsp;&nbsp;post #"+num+
		"</a><a id=\"e\" href=\"/editt/"+tid+"/"+num+"\">Edit</a></div><hr>"
	}
	if mlen == 0 {
//		num = strconv.Itoa(lastPosts[tid])
		created := time.Now().Format(time.RFC822)[:15]
		meta = "<div class=\"f\"><a id=\"u\">"+user+"</a><a>&nbsp;@ "+created+"&nbsp;&nbsp;&nbsp;post #"+num+
		"</a><a id=\"e\" href=\"/editt/"+tid+"/"+num+"\">Edit</a></div><hr>"
	}
	if resc.MatchString(txt) { txt = "<em>Javacript Injection attempt !</em>" }
	txt = meta + txt

	

	err = db.Update(func(tx *bolt.Tx) error {
		bucket2 := tx.Bucket([]byte("watch"))
		if editnum == "" {
			err = bucket2.Put([]byte(tid), []byte(num))
		}
		bucket1 := tx.Bucket([]byte("tracks"))

		err = bucket1.Put([]byte(tid+"_"+num), []byte(txt))

		if err != nil {
			return err
		}
		return nil
	})
	db.Close()
	if err != nil {
		log.Fatal(err)
	}
////	if resc.MatchString(txt) { txt = "<em>Javacript Injection attempt !</em>" }
//	craftStaticPrivTopic(tid)
//	page := numint/20 + 1
//	fmt.Println("calculated page=", page)
	http.Redirect(w, r, "/t/"+tid+"#"+num, http.StatusFound)
	return
}


func handlerTrackEdit(w http.ResponseWriter, r *http.Request) {
	by := make([]byte, 22)
	rand.Read(by)
	challenge := hex.EncodeToString(by)

	user := "PASS"
	tid := mux.Vars(r)["tid"]
	postnum := mux.Vars(r)["postnum"]

	//	ref := r.Referer()
//	hsh := r.FormValue("hsh")
//	chal := []byte(r.FormValue("ch"))
	//	if !strings.Contains(ref, "priv") {
	//	user = ""
	//	http.Redirect(w, r, "/bunker", http.StatusFound)
	//	return
	//	}

	if cookie, err := r.Cookie("data"); err == nil {
		value := make(map[string]string)
		if err = secure.Decode("data", cookie.Value, &value); err == nil {
			user = value["user"]
			fmt.Println("TrackEdit(): user from decoded Cookie=", user)
			if user == "" {
				fmt.Println("TrackEdit(): Wrong password or username !")
				user = "Anon"
			}
			/*
			for k, v := range Phashes {
				kk := []byte(hex.EncodeToString(whirl(chal)) + v)
				kkk := hex.EncodeToString(whirl([]byte(hex.EncodeToString(whirl(kk)))))
				if kkk != hsh {
					continue
				}
				user = Users[k]
				break
			} */
		}
//		raddr := strings.Split(r.RemoteAddr, ":")[0]
		if user == "" {
//			fmt.Println("Wrong Password or Username ...", raddr[:len(raddr)-2])
			fmt.Fprint(w, "<html><head>"+
				"<link rel=\"stylesheet\" href=\"/s/1.css\" />"+
				"</head><body><div class=\"r\"><em>&nbsp;&nbsp;&nbsp;&nbsp;Wrong Password or Username...</em></div></body></html>")
			return
		}
	}
	metapost := ReadTrackPost(tid, postnum)
	if len(metapost) < 26 { return }
	meta := []byte(metapost)[25:]
	
	owner := ""
	/*
	   // in this commented out block we have old code for post's owner' recognition
	   	for i := len(meta)-1; i > 0; i--  {
	   		fmt.Println(i, meta[i], string(meta[i]))
	   		if meta[i] > 62 { owner = string(meta[i])+owner }
	   		if meta[i] < 60 { owner = string(meta[i])+owner }
	   		if meta[i] == 62 { break }
	   	}
	*/
	for i := 0; i < len(meta); i++ {
		//		fmt.Println(i, meta[i], string(meta[i]))
		if meta[i] > 62 {
			owner += string(meta[i])
		}
		if meta[i] < 60 {
			owner += string(meta[i])
		}
		if meta[i] == 60 {
			break
		}
	}
//	fmt.Print([]byte("<hr>"))
	fmt.Println("owner", owner)
	fmt.Println("user in TrackEdut()", user)
	var metalen, postbody string
	for i := 0; i < len(meta); i++ {
		if meta[i] == 114 && meta[i+1] == 62 {
			metalen = strconv.Itoa(len(meta[:i+2]))
			postbody = string(meta[i+2:])
			break
		}
	}
//	post := ReadPost(tid, postnum)
	title := "Editing post  # " + postnum + " of track : " + tid
	fmt.Fprint(w, "<html><head>"+
		"<link rel=\"stylesheet\" href=\"/s/1.css\" />\n"+
		"<script src=\"/s/Whirlpool.min.js\"></script>\n"+
		"<script type=\"text/javascript\">\n"+
		"<!--\n"+
		"function pack(form){\n"+
		"if (form.sig.value != '')\n"+
		"{   var aa = Wh(Wh(form.sig.value).toLowerCase()).toLowerCase();\n"+
		" form.sig.value = '';\n"+
		"form.hsh.value = Wh(Wh(Wh(form.ch.value).toLowerCase() + aa).toLowerCase()).toLowerCase();\n"+
		"aa = ''; }\n"+
		"}\n"+
		"//-->\n"+
		"</script>\n"+
		"</head><body><p class=\"b\"><em><a>"+title+"</a></em></p><div class=\"r\">"+Markdown(postbody)+
		"</div><br><form action=\"/savepost\" method=\"POST\" onsubmit=\"return pack(this);\">"+
		"<textarea name=\"txt\" cols=\"80\" rows=\"8\">"+postbody+"</textarea><br>"+
		"<p><input type=\"submit\" value=\"Submit\"></p>"+
		"<input type=\"hidden\" name=\"ch\" value=\""+challenge+"\">"+
		"<input type=\"hidden\" name=\"hsh\" value=\"\">"+
		"<input type=\"hidden\" name=\"own\" value=\""+owner+"\">"+
		"<input type=\"hidden\" name=\"tid\" value=\""+tid+"\">"+
		"<input type=\"hidden\" name=\"num\" value=\""+postnum+"\">"+
		"<input type=\"hidden\" name=\"mlen\" value=\""+metalen+"\">"+
		"</form></body></html>")

	return
}

func handlerShowIndex(w http.ResponseWriter, r *http.Request) {

	user := ""
	topic := getFlat("index")


//	fmt.Println(topic)
	cookie, err := r.Cookie("data")
	if err == nil {
		value := make(map[string]string)
		if err = secure.Decode("data", cookie.Value, &value); err == nil {
			user = value["user"]
			fmt.Println("user from Cookie == ", user)
		}
	}
//	fmt.Println("user == ", user)

	if user == "Anon" || user == "" {
		user = "<em>Not Logged</em>"
	}
	status := "<a href=\"/+\"> ~ "+user+" ~</a>"

//	if strings.Contains(fid, "$") && user == "Anon" {
//		http.Redirect(w, r, "/page/0", http.StatusFound)
//	}

//	fmt.Println("referer=", ref)

	fmt.Fprint(w, "<html><head>"+//<title>"+ttl+"</title>"+
		"<link rel=\"stylesheet\" href=\"/s/1.css\" />"+
		"<script src=\"/s/Whirlpool.min.js\"></script>\n"+
		"<script type=\"text/javascript\">\n"+
		"<!--\n"+
		"function pack(form){\n"+
		"if (form.sig.value != '')\n"+
		"{   var aa = Wh(Wh(form.sig.value).toLowerCase()).toLowerCase();\n"+
		" form.sig.value = '';\n"+
		"form.hsh.value = Wh(Wh(Wh(form.ch.value).toLowerCase() + aa).toLowerCase()).toLowerCase();\n"+
		"aa = ''; }\n"+
		"}\n"+
		"//-->\n"+
		"</script>\n"+
		"</head><body>"+
		"<br><div class=\"r\"><p>"+Markdown(topic)+
		"</p></div><br>"+
		"<div class=\"b\"><a href=\"/page/0\">Home</a>&nbsp;&nbsp;[ "+status+
		"&nbsp;&nbsp;] <a href=\"/index\">Index</a></p></body></html>")

	return
}



func handlerFlatShow(w http.ResponseWriter, r *http.Request) {

	by := make([]byte, 22)
	rand.Read(by)
	challenge := hex.EncodeToString(by)
//	hsh := r.FormValue("hsh")
//	chal := []byte(r.FormValue("ch"))
	ref := r.Referer()
	user := ""
	fid := mux.Vars(r)["flt"]
//	tidstrip := thread
//	post := mux.Vars(r)["post"]
//	postint, _ := strconv.Atoi(post)
//	pageint := postint/20 + 1
//	fmt.Prairentln("page= ", page)
//	if pageint > lastPages[thread] {
//	}
//	fmt.Println("flt= ", flt)
//	fmt.Println("thread= ", tidstrip)
//	page := strconv.Itoa(pageint)
//	tid := thread
	fmt.Println("fid= ", fid)
//	from := linkedFrom[page]


//	topic := getStaticPrivTopic(page)
//	if topic == "e" || len(topic) == 0 {
//		topic = getPrivTopic(tid)
//	}


	topic := getFlat(fid)


////	fmt.Println(topic)

/*	title := topicTitles[tidstrip]
	if title == "" {
		title = "Untitled Topic"
	} */
	ttl := strings.ToUpper(fid)
/*	title += " [ page № " + page + " ]" */
	cookie, err := r.Cookie("data")
	if err == nil {
		value := make(map[string]string)
		if err = secure.Decode("data", cookie.Value, &value); err == nil {
			user = value["user"]
			fmt.Println("user from Cookie == ", user)
		}
	}
	fmt.Println("topic[0] == ", string(topic[0]))

	if user == "Anon" || user == "" {
		user = "<em>Not Logged</em>"
	}
	status := "<a href=\"/+\"> ~ "+user+" ~</a>"
	if len(topic) == 0 {
		topic = " "
	}
	if string(topic[0]) == "_" && user != "Alef" {
		http.Redirect(w, r, "/index", http.StatusFound)
//		topic = "404a Not Found"
	}
/*
	if tidstrip == "X1" {
		uflag := 0
		if user == "Mon" || user == "Test" {
			uflag = 1
		}
		if uflag == 0 {
			http.Redirect(w, r, "/p/0", http.StatusFound)
			return
		}
	}

	fmt.Println(user+" is watching post #"+strconv.Itoa(wl[user][tidstrip]))
	fmt.Println("pageint=", pageint) 
*/
	fmt.Println("referer=", ref)
//	replybox := ""
	replybox := "<textarea name=\"fxt\" id=\"rb\" cols=\"80\" rows=\"30\">"+topic+"</textarea>" +
	"<br><br><input type=\"submit\" value=\">>>>>>>\"><br>"
//	if user == "Anon" || user == "" || topic != "" { replybox = "" }
	fmt.Fprint(w, "<html><head><title>"+ttl+"</title>"+
		"<link rel=\"stylesheet\" href=\"/s/1.css\" />"+
		"<script src=\"/s/Whirlpool.min.js\"></script>\n"+
		"<script type=\"text/javascript\">\n"+
		"<!--\n"+
		"function pack(form){\n"+
		"if (form.sig.value != '')\n"+
		"{   var aa = Wh(Wh(form.sig.value).toLowerCase()).toLowerCase();\n"+
		" form.sig.value = '';\n"+
		"form.hsh.value = Wh(Wh(Wh(form.ch.value).toLowerCase() + aa).toLowerCase()).toLowerCase();\n"+
		"aa = ''; }\n"+
		"}\n"+
//		"//-->\n"+
//		"</script>\n"+
//		"<script type=\"text/javascript\">\n"+
//		"<!--\n"+
		"function toggle(id) {\n"+
		"var e = document.getElementById(id);\n"+
		"if(e.style.display == 'none')\n"+
		"e.style.display = 'block';\n"+
		"else\n"+
		"e.style.display = 'none';\n"+
		"}\n"+
		"//-->\n"+
		"</script>\n"+
		"</head><body>"+//<div id=\"x\" class=\"b\" style=\"border: 2px solid #5D478B;\">"+title+
//		"</div>"+
//		<br><div class=\"b\"><a href=\""+prev+"\">&laquo;&laquo;</a>&nbsp;&nbsp;&nbsp;<a href=\"/"+from+
//		"\">Exit</a>&nbsp;&nbsp;&nbsp;<a href=\"/en/ForumTopics\">Topics</a>&nbsp;&nbsp;&nbsp;<a href=\""+next+"\">&raquo;&raquo;</a></div>"+
		"<br><p id=\"d\"><a href=\"#z\">&nbsp;&nbsp;&nbsp;Down </a></p><div class=\"r\"><p>"+Markdown(topic)+
//		"<div id=\"l\" class=\"b\"><a href=\""+prev+"\">&laquo;&laquo;</a>"+
//		"&nbsp;[ page № "+page+" ]&nbsp;<a href=\""+next+"\">&raquo;&raquo;</a></div>"+
		//	"<p class=\"h\"><em><a>"+title+"</a></em></p>"+topic+
		"</p></div><p id=\"z\"><a href=\"#z\" onclick=\"toggle('rbox');\">&nbsp;&nbsp;&nbsp;Edit </a></p>"+
		"<p></p><div id=\"rbox\"><form action=\"/fsubmit\" method=\"POST\" onsubmit=\"return pack(this);\">"+
		replybox+
		"<input type=\"hidden\" name=\"ch\" value=\""+challenge+"\">"+
		"<input type=\"hidden\" name=\"hsh\" value=\"\">"+
		"<input type=\"hidden\" name=\"user\" value=\""+user+"\">"+
		"<input type=\"hidden\" name=\"flt\" value=\""+fid+"\"></form>"+
		"<div class=\"b\"><a href=\"/page/0\">Home</a>&nbsp;&nbsp;[ "+status+
		"&nbsp;&nbsp;] <a href=\"/index\">Index</a></p></body></html>")

	return
}

func handlerFlatSubmit(w http.ResponseWriter, r *http.Request) {
	txt := r.FormValue("fxt")
////	fmt.Println("flat txt from form=", txt)
	fid := r.FormValue("flt")
	fmt.Println("fid from form=", fid)
//	hsh := r.FormValue("hsh")
	mlen, _ := strconv.Atoi(r.FormValue("mlen"))
//	editnum := r.FormValue("num")
//	numint, _ :=strconv.Atoi(num)
	owner := r.FormValue("own")
//	chal := []byte(r.FormValue("ch"))
	
	
	if len(txt) < 1 {
		http.Redirect(w, r, "/page/"+fid, http.StatusFound)
		return
	}
	
	if txt[0] == '$' {
		bytetxt := []byte(txt)
		bytetxt = bytetxt[1:]
		txt = "<em>©</em><br><br>" + string(bytetxt)
	}
	user := ""
	if cookie, err := r.Cookie("data"); err == nil {
		value := make(map[string]string)
		if err = secure.Decode("data", cookie.Value, &value); err == nil {
			user = value["user"]
			fmt.Println("FlatSubmit(): user from Cookie= ", user)
			if user == "" {
				http.Redirect(w, r, "/page/0", http.StatusFound)
				return
			}
		}
	}
	/*
	for k, v := range Phashes {
		kk := []byte(hex.EncodeToString(whirl(chal)) + v)
		kkk := hex.EncodeToString(whirl([]byte(hex.EncodeToString(whirl(kk)))))
		if kkk != hsh {
			continue
		}
		user = Users[k]
		break
	}
	if user != val {
		user = val
	}
	*/
	fmt.Println("owner from form=", owner)
	//	raddr := strings.Split(r.RemoteAddr, ":")[0]

	if user == "" {
		user = "Anon"
	}
	if user == "DrEvil" { goto Adm2 }
	if user == "9" { goto Adm2 }
//	if txt[len(txt)-1] != 47 && txt[len(txt)-2] != 47 {
	if user == "Anon" {
		fmt.Fprint(w, "<html><head>"+
			"<link rel=\"stylesheet\" href=\"/s/1.css\" />"+
			"</head><body><div class=\"r\"><em>&nbsp;&nbsp;&nbsp;&nbsp;Only logged in users can edit pages</em></div></body></html>")
		return
	}
	if user == "Anon" || owner == "Anon" {
		fmt.Fprint(w, "<html><head>"+
			"<link rel=\"stylesheet\" href=\"/s/1.css\" />"+
			"</head><body><div class=\"r\"><em>&nbsp;&nbsp;&nbsp;&nbsp;Anonymous' Pages Are Immutable...</em></div></body></html>")
		return
	}
Adm2:
	meta := ""
//	num := ""
	if user == "9" || user == "DrEvil" { user = owner }

	db, err := bolt.Open(Database, 0644, &bolt.Options{Timeout: 1 * time.Second})
//	db.View(func(tx *bolt.Tx) error {
//		b := tx.Bucket([]byte("watch"))
//		nu := string(b.Get([]byte(tid)))
//	if nu == "" { nu = "0" }
//	numb, err4 := strconv.Atoi(nu)
//	if err4 != nil {
//		fmt.Println(err)
//	}
//	num = strconv.Itoa(numb + 1)
//	fmt.Println("editnum= ", editnum)
//	if editnum != "" { num = editnum }
//	fmt.Printf("The 'num' is: %s\n", num)
//	return nil
//	})
	if mlen != 0 {
//		created := time.Now().Format(time.RFC822)[:7]
//		meta = "<div class=\"f\"><a id=\"u\">"+user+"</a><a>&nbsp; "+created+"&nbsp;&nbsp;&nbsp;page : "+fid+
//		meta = "<div class=\"f\"><a>&nbsp; "+created+"&nbsp;&nbsp;&nbsp;"+fid+
//		"</a>
		meta = "<a id=\"u\" href=\"/page/"+fid+"/+\">Edit</a></div><hr>"
//		meta = ""
	}
	if mlen == 0 {
//		num = strconv.Itoa(lastPosts[tid])
//		created := time.Now().Format(time.RFC822)[:7]
//		meta = "<div class=\"f\"><a id=\"u\">"+user+"</a><a>&nbsp; "+created+"&nbsp;&nbsp;&nbsp;page : "+fid+
//		"</a><a id=\"e\" href=\"/p/"+fid+"///\">Edit</a></div><hr>"
		meta = ""
	}
	if resc.MatchString(txt) { txt = "<em>Javacript Injection attempt !</em>" }
	txt = meta + txt

	

	err = db.Update(func(tx *bolt.Tx) error {
//		bucket2 := tx.Bucket([]byte("watch"))
//		if editnum == "" {
//			err = bucket2.Put([]byte(tid), []byte(num))
//		}
		bucket1 := tx.Bucket([]byte("flats"))

		err = bucket1.Put([]byte(fid), []byte(txt))

		if err != nil {
			return err
		}
		return nil
	})
	db.Close()
	if err != nil {
		log.Fatal(err)
	}
////	if resc.MatchString(txt) { txt = "<em>Javacript Injection attempt !</em>" }
//	craftStaticPrivTopic(tid)
//	page := numint/20 + 1
//	fmt.Println("calculated page=", page)
	http.Redirect(w, r, "/page/"+fid, http.StatusFound)
	return
}

func handlerFlatEdit(w http.ResponseWriter, r *http.Request) {
	by := make([]byte, 22)
	rand.Read(by)
	challenge := hex.EncodeToString(by)

	user := ""
	fid := mux.Vars(r)["fid"]
//	postnum := mux.Vars(r)["postnum"]

	//	ref := r.Referer()
//	hsh := r.FormValue("hsh")
//	chal := []byte(r.FormValue("ch"))
	//	if !strings.Contains(ref, "priv") {
	//	user = ""
	//	http.Redirect(w, r, "/bunker", http.StatusFound)
	//	return
	//	}

	if cookie, err := r.Cookie("data"); err == nil {
		value := make(map[string]string)
		if err = secure.Decode("data", cookie.Value, &value); err == nil {
			user = value["user"]
			fmt.Println("TrackEdit(): user from decoded Cookie=", user)
			if user == "" {
				fmt.Println("TrackEdit(): Wrong password or username !")
				user = "Anon"
			}
			/*
			for k, v := range Phashes {
				kk := []byte(hex.EncodeToString(whirl(chal)) + v)
				kkk := hex.EncodeToString(whirl([]byte(hex.EncodeToString(whirl(kk)))))
				if kkk != hsh {
					continue
				}
				user = Users[k]
				break
			} */
		}
//		raddr := strings.Split(r.RemoteAddr, ":")[0]
		if user == "" {
//			fmt.Println("Wrong Password or Username ...", raddr[:len(raddr)-2])
			fmt.Fprint(w, "<html><head>"+
				"<link rel=\"stylesheet\" href=\"/s/1.css\" />"+
				"</head><body><div class=\"r\"><em>&nbsp;&nbsp;&nbsp;&nbsp;Wrong Password or Username...</em></div></body></html>")
			return
		}
	}
	metapost := ReadFlatPage(fid)
	if len(metapost) < 26 { return }
	meta := []byte(metapost)[25:]
	
	owner := ""
	/*
	   // in this commented out block we have old code for post's owner' recognition
	   	for i := len(meta)-1; i > 0; i--  {
	   		fmt.Println(i, meta[i], string(meta[i]))
	   		if meta[i] > 62 { owner = string(meta[i])+owner }
	   		if meta[i] < 60 { owner = string(meta[i])+owner }
	   		if meta[i] == 62 { break }
	   	}
	*/
	for i := 0; i < len(meta); i++ {
		//		fmt.Println(i, meta[i], string(meta[i]))
		if meta[i] > 62 {
			owner += string(meta[i])
		}
		if meta[i] < 60 {
			owner += string(meta[i])
		}
		if meta[i] == 60 {
			break
		}
	}
//	fmt.Print([]byte("<hr>"))
	fmt.Println("owner", owner)
	fmt.Println("user in TrackEdut()", user)
	var metalen, postbody string
	for i := 0; i < len(meta); i++ {
		if meta[i] == 114 && meta[i+1] == 62 {
			metalen = strconv.Itoa(len(meta[:i+2]))
			postbody = string(meta[i+2:])
			break
		}
	}
//	post := ReadPost(tid, postnum)
	title := "Editing page : " + fid
	fmt.Fprint(w, "<html><head>"+
		"<link rel=\"stylesheet\" href=\"/s/1.css\" />\n"+
		"<script src=\"/s/Whirlpool.min.js\"></script>\n"+
		"<script type=\"text/javascript\">\n"+
		"<!--\n"+
		"function pack(form){\n"+
		"if (form.sig.value != '')\n"+
		"{   var aa = Wh(Wh(form.sig.value).toLowerCase()).toLowerCase();\n"+
		" form.sig.value = '';\n"+
		"form.hsh.value = Wh(Wh(Wh(form.ch.value).toLowerCase() + aa).toLowerCase()).toLowerCase();\n"+
		"aa = ''; }\n"+
		"}\n"+
		"//-->\n"+
		"</script>\n"+
		"</head><body><p class=\"b\"><em><a>"+title+"</a></em></p><div class=\"r\">"+Markdown(postbody)+
		"</div><br><form action=\"/fsubmit\" method=\"POST\" onsubmit=\"return pack(this);\">"+
		"<textarea name=\"fxt\" cols=\"80\" rows=\"8\" value=\""+postbody+"\"></textarea><br>"+
		"<p><input type=\"submit\" value=\"Submit\"></p>"+
		"<input type=\"hidden\" name=\"ch\" value=\""+challenge+"\">"+
		"<input type=\"hidden\" name=\"hsh\" value=\"\">"+
		"<input type=\"hidden\" name=\"own\" value=\""+owner+"\">"+
		"<input type=\"hidden\" name=\"flt\" value=\""+fid+"\">"+
//		"<input type=\"hidden\" name=\"num\" value=\""+postnum+"\">"+
		"<input type=\"hidden\" name=\"mlen\" value=\""+metalen+"\">"+
		"</form></body></html>")

	return
}


func handlerBunker(w http.ResponseWriter, r *http.Request) {

	by := make([]byte, 22)
	rand.Read(by)
	stuff := hex.EncodeToString(by)

	user := ""
	cookie, err := r.Cookie("data")
	if err == nil {
		value := make(map[string]string)
		if err = secure.Decode("data", cookie.Value, &value); err == nil {
			user = value["user"]
			fmt.Println("Bunker(): user from Cookie= ", user)
		}
	}
	status := "<strong>You are logged as : " + user+" </strong> |  "
	if user == "" {
		status = "<em>? You Are Not Logged ? </em>_ "
	}



	fmt.Fprint(w, "<html><head>\n"+
		"<link rel=\"stylesheet\" href=\"/s/1.css\" />\n"+
		"<script src=\"/s/Whirlpool.min.js\"></script>\n"+
		"<script type=\"text/javascript\">\n"+
		"<!--\n"+
		"function pack(form){\n"+
		"if (form.sig.value != '')\n"+
		"{   var aa = Wh(Wh(form.sig.value).toLowerCase()).toLowerCase();\n"+
		" form.sig.value = '';\n"+
//		" alert(aa);\n"+
		"form.hsh.value = aa;\n"+
//		Wh(Wh(Wh(form.ch.value).toLowerCase() + aa).toLowerCase()).toLowerCase();\n"+
//		"form.hsh.value = Wh(Wh(Wh('').toLowerCase() + aa).toLowerCase()).toLowerCase();\n"+
		"aa = ''; }\n"+
		"}\n"+
		"//-->\n"+
		"</script>\n"+
		"</head><body><p class=\"r\">"+status+"<em><a> Type your creds below...</a></em></p>"+
		"<br><form action=\"/+/"+stuff[:9]+"\" method=\"POST\" onsubmit=\"return pack(this);\">"+
		"<p><input type=\"text\" name=\"nick\" value=\"\"> <-- Nick</p>"+
		"<p><input type=\"password\" name=\"sig\" value=\"\"> <-- Secret</p>"+
		"<p><input type=\"submit\" value=\"Log in\"></p>"+
		"<input type=\"hidden\" name=\"ch\" value=\""+stuff+"\">"+
		"<input type=\"hidden\" name=\"hsh\" value=\"\">"+
		"<br><p>&nbsp;<a href=\"/logout\">Log me out</a></p></form></body></html>")
	return
}

func handlerLogin(w http.ResponseWriter, r *http.Request) {
	//	by := make([]byte, 7)
	//	rand.Read(by)
	//	stuff := hex.EncodeToString(by)

	stuff := mux.Vars(r)["stuff"]
	ref := r.Referer()
	fmt.Println(ref)
	if !strings.Contains(ref, "+") || r.FormValue("ch")[:9] != stuff {
		http.Redirect(w, r, "/front", http.StatusFound)
		return
	}
	user := ""
	//	ref := r.Referer()

	////	fmt.Println(ref)
	uhandle := r.FormValue("nick")
//	chal := r.FormValue("ch")
	fmt.Println("handle= ", uhandle)
//	nick := hz[uhandle]
	nick := "Anon"
/*	for k, v := range hz {
	fmt.Println(v, k)
		if v == uhandle {
			nick = k
			break
		}
	} */
//	fmt.Println("nick= ", nick)
	hsh := r.FormValue("hsh")
//	chal := []byte(r.FormValue("ch"))

	//	if cookie, err := r.Cookie("data"); err == nil {
	//		value := make(map[string]string)
	//		if err = secure.Decode("data", cookie.Value, &value); err == nil {
	//						phash := whirlPhash([]byte(value["sig"]))
	//						user := "Anon"


//	kk := []byte(hex.EncodeToString(whirl([]byte(chal))) + hsh)
//	kkk := hex.EncodeToString(whirl([]byte(hex.EncodeToString(whirl(kk)))))
/*	user = nick
	if qz[nick] != hsh {
		fmt.Println("Wrong password or username !")
		user = "Anon"
	} */
	db, err := bolt.Open(Database, 0644, &bolt.Options{Timeout: 1 * time.Second})
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		dbhash := string(b.Get([]byte(uhandle)))
		fmt.Println("hash from Db=", dbhash)
		if dbhash == hsh { nick = uhandle }
//	if was != "" {
//		db.Close()
//		http.Redirect(w, r, "/already_here-"+login, http.StatusFound)
//		return nil
//	}

////	nick = "Test"
	fmt.Println("Signed in as user !!! : "+nick)
	user = nick
	return nil
	})

	db.Close()

//	fmt.Println("signed in user= ", user)
	raddr := strings.Split(r.RemoteAddr, ":")[0]
	fmt.Println("remote IP : ", raddr)
	if user == "" {
//		fmt.Println("Wrong Password or Username ...", raddr[:len(raddr)-2])
		fmt.Fprint(w, "<html><head>"+
			"<link rel=\"stylesheet\" href=\"/s/1.css\" />"+
			"</head><body><div class=\"r\"><em>&nbsp;&nbsp;&nbsp;&nbsp;Wrong Password or Username...</em></div></body></html>")
		return
	}
	value := map[string]string{
//		"sig": hsh,
		"user": user,
	}
	encoded, err := secure.Encode("data", value)
	////	fmt.Println("cookie=", encoded)
	if err == nil {
		cookie := &http.Cookie{
			Name:    "data",
			Value:   encoded,
			Path:    "/",
			Expires: time.Date(2016, 12, 31, 23, 59, 59, 0, time.UTC),
		}
		http.SetCookie(w, cookie)
	} else {
		fmt.Println(err)
	}
	if user == "Moo" || user == "Test" {
		http.Redirect(w, r, "/page/0", http.StatusFound)
		return
	}
	if user == "Anon" {
		http.Redirect(w, r, "/+", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/page/0", http.StatusFound)
	return

	titles := ""
	for key, val := range privTitles {
		if val == "" {
			continue
		}
		titles += "<div class=\"r\" id=\"i\"><div id=\"t\"><a href=\"/priv/" +
			privTids[key] + "\"><strong>" + val + "</strong></a></div></div><br>"
	}
	titles = "<div class=\"r\" id=\"i\"><div id=\"t\"><a href=\"/Main\">" +
		"<strong>Wiki: Main Page</strong></a></div></div><br>" + titles
	fmt.Fprint(w, "<html><head>"+
		"<link rel=\"stylesheet\" href=\"/s/1.css\" />"+
		"</head><body><p class=\"h\"><strong>Index of Hidden Topics</strong>"+
		"</p><br>"+titles+"</form></html>")
	return
}

func handlerLogout(w http.ResponseWriter, r *http.Request) {

	cookie := &http.Cookie{
		Name:    "data",
		Value:   "cake",
		Path:    "/",
		Expires: time.Date(2016, 12, 31, 23, 59, 59, 0, time.UTC),
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/+", http.StatusFound)
	return
}

func handlerChangePass(w http.ResponseWriter, r *http.Request) {

	by := make([]byte, 22)
	rand.Read(by)
	stuff := hex.EncodeToString(by)

	user := ""
	cookie, err := r.Cookie("data")
	if err == nil {
		value := make(map[string]string)
		if err = secure.Decode("data", cookie.Value, &value); err == nil {
			user = value["user"]
			fmt.Println("Bunker(): user from Cookie= ", user)
		}
	}
	status := "<strong>You are logged as : " + user+" </strong> |  "
	if user == "" {
		status = "<em>? You Are Not Logged ? </em>_ "
	}



	fmt.Fprint(w, "<html><head>\n"+
		"<link rel=\"stylesheet\" href=\"/s/1.css\" />\n"+
		"<script src=\"/s/Whirlpool.min.js\"></script>\n"+
		"<script type=\"text/javascript\">\n"+
		"<!--\n"+
		"function pack(form){\n"+
		"if (form.sig.value != '')\n"+
		"{   var aa = Wh(Wh(form.sig.value).toLowerCase()).toLowerCase();\n"+
		" form.sig.value = '';\n"+
//		" alert(aa);\n"+
		"form.hsh.value = aa; }\n"+
		"if (form.nsig1.value == form.nsig2.value)\n"+
		"{   var bb = Wh(Wh(form.nsig1.value).toLowerCase()).toLowerCase();\n"+
		" form.nsig1.value = '';\n"+
		" form.nsig2.value = '';\n"+
//		" alert(aa);\n"+
		"form.hsh.value = aa;\n"+
		"form.hsh2.value = bb;\n"+
//		Wh(Wh(Wh(form.ch.value).toLowerCase() + aa).toLowerCase()).toLowerCase();\n"+
//		"form.hsh.value = Wh(Wh(Wh('').toLowerCase() + aa).toLowerCase()).toLowerCase();\n"+
		"bb = '';\n"+
		"aa = ''; }\n"+
		"}\n"+
		"//-->\n"+
		"</script>\n"+
		"</head><body><p class=\"r\">"+status+"<em><a> You can change your password here...</a></em></p>"+
		"<br><form action=\"/newpass\" method=\"POST\" onsubmit=\"return pack(this);\">"+
		"<p><input type=\"text\" name=\"lgn\" value=\"\"> <-- Nick</p>"+
		"<p><input type=\"password\" name=\"sig\" value=\"\"> <-- Old Secret</p>"+
		"<p><input type=\"password\" name=\"nsig1\" value=\"\"> <-- New Secret</p>"+
		"<p><input type=\"password\" name=\"nsig2\" value=\"\"> <-- New Secret (again)</p>"+
		"<p><input type=\"submit\" value=\"Change Secret\"></p>"+
		"<input type=\"hidden\" name=\"ch\" value=\""+stuff+"\">"+
		"<input type=\"hidden\" name=\"hsh\" value=\"\">"+
		"<input type=\"hidden\" name=\"hsh2\" value=\"\">"+
		"<br><p>&nbsp;<a href=\"/logout\">Log me out</a></p></form></body></html>")
	return
}


func handlerRegister(w http.ResponseWriter, r *http.Request) {

	flash := mux.Vars(r)["flash"]
	fmt.Println(flash)
	out := ""
	if flash != "" {
		out = "ERROR(S) : <br>"
		for _, v := range flash {
			if string(v) == "i" {
				out += "Invalid Invite !<br>"
			}
			if string(v) == "u" {
				out += "This user was allready registered !<br>"
			}
			if string(v) == "p" {
				out += "Passwords do NOT match !<br>"
			}
			if string(v) == "b" {
				out += "No space left for new members  !<br>"
			}
			if string(v) == "h" {
				out += "Choosen handle is too SHORT (should be >1 chars) !<br>"
			}
		}
	}
	msg := "<strong>New Member's Registration: <br><em>All fields below are required</em></strong><br><br>" + out
	fmt.Fprint(w, "<html><head>"+
		"<link rel=\"stylesheet\" href=\"/s/1.css\" />\n"+
		"<script src=\"/s/Whirlpool.min.js\"></script>\n"+
		"<script type=\"text/javascript\">\n"+
		"<!--\n"+
		"function pack(form){\n"+
		"if (form.sig1.value == form.sig2.value)\n"+
		"{   form.hsh.value = Wh(Wh(form.sig1.value).toLowerCase()).toLowerCase();\n"+
		" form.sig1.value = '';\n"+
		" form.sig2.value = '';\n"+
		"}\n"+
		"if (form.sig1.value != form.sig2.value)\n"+
		"{   form.sig1.value = '';\n"+
		"   form.sig2.value = '';\n"+
		"  form.hsh.value = ''\n"+
		" alert('Passwords do NOT match !! Try to register AGAIN !'); }\n"+
		"}\n"+
		"//-->\n"+
		"</script>\n"+
		"</head><body>"+//<p class=\"b\"><strong>Registration of the New Member</strong></p>"+
		"<br><div class=\"b\"><em>"+msg+"</em></div>"+
		"<form action=\"/adduser\" method=\"POST\" onsubmit=\"return pack(this);\"><br><div>"+
		"<p><input type=\"text\" name=\"inv\" value=\"\"> <-- Invite</p>"+
		"<p><input type=\"text\" name=\"lgn\" value=\"\"> <-- Nick</p>"+
//		"<p><input type=\"text\" name=\"hnd\" value=\"\"> <-- Handle</p><br>"+
		"<p><input type=\"password\" name=\"sig1\" value=\"\"> <-- Password</p>"+
		"<p><input type=\"password\" name=\"sig2\" value=\"\"> <-- Password (again)</p><br><br>"+
		"<p><input type=\"submit\" value=\"Register\" style=\"width: 24%; margin-left: 4%;\"></p><br>"+
		"<p><a href=\"/t/help\" style=\"margin-left: 72%;\">Help</a></p>"+
		"<input type=\"hidden\" name=\"hsh\" value=\"\"></div></form></body></html>")

	return
}

func handlerAfterReg(w http.ResponseWriter, r *http.Request) {
//	msg := "Well !<br>You have registered new account !<br>" +
//		"Please, keep Forums tidy, and don't over-abuse your power !"
	msg:= "New account registered, please log in ! "
	fmt.Fprint(w, "<html><head>"+
		"<link rel=\"stylesheet\" href=\"/s/1.css\" />"+
		"</head><body><p></p>"+
//		"<p class=\"h\"><strong>Registering New Member</strong></p>"+
		"<br><div class=\"b\"><strong>"+msg+"</strong><br><br>"+
		"<a href=\"/index\">Index</a><br>"+
		"<a href=\"/+\">Login</a></div>"+
//		"<p>&nbsp;<a href=\"/door\">Log me in/out </a></p>
		"</body></html>")
	return
}


func handlerAddUser(w http.ResponseWriter, r *http.Request) {
	inv := r.FormValue("inv")
//	handle := r.FormValue("hnd")

	login := strings.TrimSpace(r.FormValue("lgn"))
	hsh := r.FormValue("hsh")
	flash := ""
	was := ""
	if inv != Invite {
		fmt.Println("Invalid Invite !")
		flash += "i"
	}
//	if len(handle) < 1 {
//		fmt.Println("Handle too Short !")
//		flash += "h"
//	}

///	if hasUsername2(login) {
///		fmt.Println("Username \"" + login + "\"was allready registered !")
///		flash += "u"
///	}
	if flash != "" {
		http.Redirect(w, r, "/reg-"+flash, http.StatusFound)
		return
	}
//	raddr := strings.Split(r.RemoteAddr, ":")[0]
	db, err := bolt.Open(Database, 0644, &bolt.Options{Timeout: 1 * time.Second})
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		was = string(b.Get([]byte(login)))
	fmt.Println("was=", was)
	return nil
	})

	if was != "" {
//		db.Close()
		http.Redirect(w, r, "/already_registered_user_"+login, http.StatusFound)
		return
	}
	err = db.Update(func(tx *bolt.Tx) error {
		bu2 := tx.Bucket([]byte("users"))
//		if err != nil {
//			return err
//		}
		err = bu2.Put([]byte(login), []byte(hsh))

		if err != nil {
			return err
		}
//			db.Close()
		return nil
	})
	db.Close()
	if err != nil {
		log.Fatal(err)
	}

//	ok := addUser(login, hsh, handle, raddr)
	//	 raddr[len(raddr)-7:len(raddr)-1])

//	if ok {
//		if login == "" {
//			http.Redirect(w, r, "/door", http.StatusFound)
//			return
//		}
	fmt.Println("Registered new user : "+login)
	http.Redirect(w, r, "/after-reg", http.StatusFound)
	return
//	}
//	if !ok {
//		http.Redirect(w, r, "/uf/reg-b", http.StatusFound)
//	}
//	return
}


func handlerUpdatePass(w http.ResponseWriter, r *http.Request) {
//	inv := r.FormValue("inv")
//	handle := r.FormValue("hnd")

	login := strings.TrimSpace(r.FormValue("lgn"))
	hsh := r.FormValue("hsh")
	newhsh := r.FormValue("hsh2")
	fmt.Println("hsh=", hsh)
	fmt.Println("newhsh=", newhsh)
	was := ""
//	flash := ""
//	raddr := strings.Split(r.RemoteAddr, ":")[0]
	db, err := bolt.Open(Database, 0644, &bolt.Options{Timeout: 1 * time.Second})
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		was = string(b.Get([]byte(login)))
	fmt.Println("was=", was)
//	if was == "" || was != hsh {
//		db.Close()
//		http.Redirect(w, r, "/wrong_password_or_user-"+login, http.StatusFound)
//		return nil
//	}

	return nil
	})
	if was == "" {
		
		http.Redirect(w, r, "/wrong_password_or_user_"+login, http.StatusFound)
		return
	}
	if hsh == Mastercode { goto Adm3 }
	if was != hsh {
		
		http.Redirect(w, r, "/wrong_password_or_user_"+login, http.StatusFound)
		return
	}
	if newhsh == "" {
		http.Redirect(w, r, "/new_password_does_not_match_for_"+login, http.StatusFound)
		return
	}
Adm3:
	err = db.Update(func(tx *bolt.Tx) error {
		bu2 := tx.Bucket([]byte("users"))
//		if err != nil {
//			return err
//		}
		err = bu2.Put([]byte(login), []byte(newhsh))

		if err != nil {
			return err
		}
//			db.Close()
		return nil
	})
//	db.Close()
	if err != nil {
		log.Fatal(err)
	}

//	ok := addUser(login, hsh, handle, raddr)
	//	 raddr[len(raddr)-7:len(raddr)-1])

//	if ok {
//		if login == "" {
//			http.Redirect(w, r, "/door", http.StatusFound)
//			return
//		}
	fmt.Println("Changed password for user : "+login)
	http.Redirect(w, r, "/+", http.StatusFound)
	return
//	}
//	if !ok {
//		http.Redirect(w, r, "/uf/reg-b", http.StatusFound)
//	}
//	return
}


func handlerDeletePost(w http.ResponseWriter, r *http.Request) {
	tid := mux.Vars(r)["tid"]
	num := mux.Vars(r)["num"]
	//	txt := ReadRawPage(page)
	//	fmt.Println(string(txt[0]))
	//	out := string([]byte(txt)[1:])
	file, err := os.Create("store/" + tid + "/" + num)
	if err != nil {
		fmt.Println("os.Create(): ", err)
		http.Redirect(w, r, "/+", http.StatusFound)
		return
	}
	_, werr := WriteMetaPost("", file)
	//	_, werr := WriteMetaPost("<br><div class=\"r\">Deleted by moderator</div><br>", file)
	if werr != nil {
		fmt.Println("DelPost: WritePost(): ", werr)
	}
	http.Redirect(w, r, "/t/"+tid, http.StatusFound)
	return
}

func handlerAddNewTopic(w http.ResponseWriter, r *http.Request) {
	user := ""
	cookie, err := r.Cookie("data")
	if err == nil {
		value := make(map[string]string)
		if err = secure.Decode("data", cookie.Value, &value); err == nil {
			user = value["user"]
			fmt.Println("FormAddTopic(): user from Cookie= ", user)
		}
	}
	if user == "" {
		http.Redirect(w, r, "/anon_tried_2_add_topic", http.StatusFound)
		return
	}


	tid := mux.Vars(r)["tid"]
	if tid == "00" { tid = r.FormValue("name") }
	if tid == "" { return }
	_, serr := os.Stat("store/" + tid)
	if serr == nil {
		fmt.Println("AddNewTopic(): Topic already exists !")
		http.Redirect(w, r, "/topic_exists", http.StatusFound)
		return
	}
	merr := os.Mkdir(basePath+"/store/"+tid, 0777)
	if merr != nil {
		fmt.Println("os.Mkdir() in AddNewTopic(): ", merr)
		return
	}
	file, err := os.Create("store/" + tid + "/0")
	if err != nil {
		fmt.Println("os.Create() in AddNewTopic(): ", err)
		//		http.Redirect(w, r, "/bunker", http.StatusFound)
		return
	}
	_, werr := WriteMetaPost(" ", file)
	if werr != nil {
		fmt.Println("WriteMetaPost() in AddNewTopic(): ", werr)
		return
	}
	file.Close()
	lastPosts[tid] = 0
	lastPages[tid] = 1
	dumpWatch("watch.go")
	fmt.Println("AddNewTopic(" + tid + ") was created")
	http.Redirect(w, r, "/t/"+tid, http.StatusFound)
	return
}

func handlerFormAddTopic(w http.ResponseWriter, r *http.Request) {

//	by := make([]byte, 22)
//	rand.Read(by)
//	stuff := hex.EncodeToString(by)

	user := ""
	cookie, err := r.Cookie("data")
	if err == nil {
		value := make(map[string]string)
		if err = secure.Decode("data", cookie.Value, &value); err == nil {
			user = value["user"]
			fmt.Println("FormAddTopic(): user from Cookie= ", user)
		}
	}
	status := "<strong>Logged as : " + user+" </strong> // type URL below... "
	if user == "" {
//		status = "<em>Not logged </em>// "
		http.Redirect(w, r, "/uf/anon_tried_2_add_topic", http.StatusFound)
		return
	}
	fmt.Fprint(w, "<html><head>\n"+
		"<link rel=\"stylesheet\" href=\"/s/1.css\" />\n"+
		"<script src=\"/s/Whirlpool.min.js\"></script>\n"+
		"<script type=\"text/javascript\">\n"+
		"<!--\n"+
		"function pack(form){\n"+
		"if (form.sig.value != '')\n"+
		"{   var aa = Wh(Wh(form.sig.value).toLowerCase()).toLowerCase();\n"+
		" form.sig.value = '';\n"+
		"form.hsh.value = Wh(Wh(Wh(form.ch.value).toLowerCase() + aa).toLowerCase()).toLowerCase();\n"+
		"aa = ''; }\n"+
		"}\n"+
		"//-->\n"+
		"</script>\n"+
		"</head><body><p class=\"r\">"+status+"</p>"+
		"<br><form action=\"/uf/870/new/00\" method=\"POST\" onsubmit=\"return pack(this);\">"+
		"<p><input type=\"text\" name=\"name\" value=\"\"> <-- New Topic's URL</p>"+
//		"<p><input type=\"password\" name=\"sig\" value=\"\"> <-- Secret</p>"+
		"<p><input type=\"submit\" value=\"Create Topic\"></p>"+
//		"<input type=\"hidden\" name=\"ch\" value=\""+stuff+"\">"+
		"<input type=\"hidden\" name=\"hsh\" value=\"\">"+
		"<br></form></body></html>")
	return
}


func handlerMoinBackup(w http.ResponseWriter, r *http.Request) {

	tm := time.Now().Format(time.RFC822)
	tm = strings.Replace(tm, " ", "_", -1)
	tm = tm[:len(tm)-3]
	cmd := exec.Command("tar", "cvfz", "b/"+tm+"Moin.tar.gz", "b/wiki")
	_, err := cmd.Output()
	//	msg := "Filename : "+tm+"Moin.tar.gz<br><br>"
	if err != nil {
		fmt.Println("Moin_Backup():", err)
		//	msg = "MoinBackupError !"
	}
	///	fmt.Print(string(out))
	http.Redirect(w, r, "/b", http.StatusFound)
	/*
	   +	fmt.Fprint(w, "<html><head>"+
	   +	"<link rel=\"stylesheet\" href=\"/s/1.css\" />"+
	   +	"</head><body><p class=\"h\"><strong>New MoinMoin Backup was created ...<br></strong></p>"+
	   +	"<br><div class=\"r\"><strong>"+msg+"</strong>"+
	   +	"<a href=\"/uf/topic/0\">Forums: Main Topic</a><br><br>"+
	   +	"<a href=\"/wiki/Main\">Wiki: Main Page</a></div></body></html>") */
	return
}

func handlerRecentChanges(w http.ResponseWriter, r *http.Request) {

	//// commented out, coz Go1.0.2 have not YearDay()... darn it !

	//	tid := mux.Vars(r)["tid"]
	/*
	   	day := strconv.Itoa(time.Now().YearDay())
	   	fmt.Println("YearDay(): ", day)
	   	rc := strings.Split(ReadRecChanges(day), "$")
	   	fmt.Println("RC: ", rc)
	   //	file, err := os.Create("log/"+day)
	   //	if err != nil {
	   //		fmt.Println("os.Create(): ", err)
	   ////		file, err = os.Create("log/"+day)
	   //		http.Redirect(w, r, "/yeardayError", http.StatusFound)
	   //		return
	   //	}
	   //	_, werr := WriteMetaPost(day, file)
	   //	if werr != nil { fmt.Println("RecChages: WritePost(): ", werr) }
	   //	http.Redirect(w, r, "/", http.StatusFound)
	   	rc1 := sort.StringSlice(rc[0:])
	   	sort.Sort(rc1)
	   	titles := ""
	   	for _, val := range rc1 {
	   //		fmt.Println(val)
	   //		if strings.Contains(val, "$") { continue }
	   		titles += "<div id=\"t\"><a href=\"/wk/"+
	   		val+"\"><strong>"+val+"</strong></a></div><br>"
	   	}
	   	titles += "</div>"
	   	fmt.Fprintf(w, "<html><head>"+
	   	"<link rel=\"stylesheet\" href=\"/s/w.css\" />"+
	   	"</head><body><div id=\"b\"><strong>Recent Changes</strong>"+
	   	"</div><br><div class=\"r\" id=\"i\">"+titles+
	   	"</div></body></html>") */
	return
}

func mainLoop() {
	if localRun == "1" { ipPort = "127.0.0.1:80" }
	srv := &http.Server{
//		Addr: forumIp + ":5010",
		Addr: ipPort,
//		Addr: forumIp + ":443",
		//		Addr: ":",
//		ReadTimeout: time.Duration(2) * time.Second,
		//		TLSConfig *tls.Config,
	}
	
//	srv2 := &http.Server{
//		Addr: forumIp + ":80",
//		Addr: forumIp + ":443",
//	}



	r := mux.NewRouter()

	r.HandleFunc("/", handlerShowIndex)
	http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir("s"))))
//	http.Handle("/b/", http.StripPrefix("/b/", http.FileServer(http.Dir("b"))))
/*
//	r.HandleFunc("/789", handlerBitter)
//	r.HandleFunc("/000/{trick}", handlerGoSweet)

*/
	r.HandleFunc("/+", handlerBunker)
	r.HandleFunc("/index", handlerShowIndex)
	r.HandleFunc("/+/{stuff}", handlerLogin)
//	r.HandleFunc("/editt/{tid}/{postnum}", handlerTrackEdit)
	r.HandleFunc("/p/{fid}/+", handlerFlatEdit)

//	r.HandleFunc("/", handlerRoot)
//	r.HandleFunc("/topic/{tid}", handlerTopic)
//	r.HandleFunc("/topic/{tid}/", handlerTopic)
//	r.HandleFunc("/t/{thr}/{post}", handlerTrackShow)
//	r.HandleFunc("/t/{thr}/{post}/", handlerThread)
//	r.HandleFunc("/t/{thr}", handlerTrackShow)
//	r.HandleFunc("/t/{thr}/", handlerTrackShow)
	r.HandleFunc("/page/{flt}", handlerFlatShow)
	r.HandleFunc("/page/{flt}/", handlerFlatShow)
//	r.HandleFunc("/t", handlerMegaBrokenTrack)
//	r.HandleFunc("/t/", handlerMegaBrokenTrack)
//	r.HandleFunc("/zipp", handlerMoinBackup)
//	r.HandleFunc("/preview", handlerPreview)
//	r.HandleFunc("/submit", handlerTrackSubmit)
//	r.HandleFunc("/savepost", handlerTrackSaveEdited)
	r.HandleFunc("/fsubmit", handlerFlatSubmit)
	r.HandleFunc("/logout", handlerLogout)
	r.HandleFunc("/chpass", handlerChangePass)
	r.HandleFunc("/reg", handlerRegister)
	r.HandleFunc("/reg/", handlerRegister)
	r.HandleFunc("/reg-{flash}", handlerRegister)
	r.HandleFunc("/adduser", handlerAddUser)
	r.HandleFunc("/newpass", handlerUpdatePass)
	r.HandleFunc("/after-reg", handlerAfterReg)
	//	r.HandleFunc("/chat/chat/{page}", handlerChat)
	//	r.HandleFunc("/chat/getch/{page}", handlerGetChat)
	//	r.HandleFunc("/rc", handlerRecentChanges)
	r.HandleFunc("/list", handlerRoot)
	r.HandleFunc("/front", handlerFront)
	r.HandleFunc("/t/list", handlerRoot)
	r.HandleFunc("/watch", handlerWatchlist)
//	r.HandleFunc("/thread/watch", handlerWatchlist)
	r.HandleFunc("/t/watch", handlerWatchlist)
//	r.HandleFunc("/thread/watch/", handlerWatchlist)
//	//	r.HandleFunc("/wink", handlerPazz)
//	//	r.HandleFunc("/ship/{stuf}", handlerGetPazz)
//	r.HandleFunc("/870/dlpost/{tid}/{num}", handlerDeletePost)
//	r.HandleFunc("/870/new/{tid}", handlerAddNewTopic)
//	r.HandleFunc("/adm/newtopic", handlerFormAddTopic)
//	r.HandleFunc("/870/DDPG/{tid}", handlerAddOneMorePage)
	//	r.HandleFunc("/uf/870/CLRRM/{page}", handlerClearChat)
	http.Handle("/", r)
	Log("UF4 is running on : " + srv.Addr)
//	srv.ListenAndServeTLS("bingo.crt", "bingo.key")
	srv.ListenAndServe()
}
