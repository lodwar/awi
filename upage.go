
package main

import (
	"io"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"strconv"
	"log"
	"github.com/boltdb/bolt"
)


func WritePost(post string, writer io.Writer) (n int64, err error) {
	
//	fmt.Println(Format(post))
	nint, err := writer.Write([]byte(Format(post)))
	n = int64(nint)
	return n, err
}

func WriteMetaPost(post string, writer io.Writer) (n int64, err error) {
	
	nint, err := writer.Write([]byte(post))
	n = int64(nint)
	return n, err
}

//UF2 method
func WritePost2(post string, writer io.Writer) (n int, err error) {
	
	n, err = writer.Write([]byte(post))
//	n = int64(nint)
	return n, err
}

func WritePrivPost(post string, writer io.Writer) (n int64, err error) {

	nint, err := writer.Write([]byte(Format(post)))
	n = int64(nint)
	return n, err
}

//UF2 method
func  ReadPost(tid, pnum string) string {
	content, err := ioutil.ReadFile(filepath.Join(basePath+"/store/"+tid+"/" , pnum))
	if err != nil {
		fmt.Print(err)
	}
//	fmt.Println("+"+pnum+"+")
	return string(content)
}

func ReadPost2(p string, reader io.Reader) (n int64, err error) {

	nint, err := reader.Read([]byte(p))
	n = int64(nint)
	return n, err
}

func  ReadTrackPost(tid, pnum string) string {
	post := ""
	db, err := bolt.Open(Database, 0644, nil)
	if err != nil {
		log.Fatal(err, "333")
	}
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
	bucket := tx.Bucket([]byte("tracks"))
	if bucket == nil {
		fmt.Println("Bucket <tracks> not found in Db !")//, err)
	}
//	y, _ := strconv.Atoi(pnum)
//	pnum = strconv.Itoa(y + 1)
	fmt.Println(string(pnum))
	fmt.Println(tid+"_"+pnum)

	post = string(bucket.Get([]byte(tid+"_"+pnum)))
//        fmt.Println(string(post))
	return nil
	})
	return post
}

func  ReadFlatPage(fid string) string {
	post := ""
	db, err := bolt.Open(Database, 0644, nil)
	if err != nil {
		log.Fatal(err, "990")
	}
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
	bucket := tx.Bucket([]byte("flats"))
	if bucket == nil {
		fmt.Println("Bucket <flats> not found in Db !")//, err)
	}
//	y, _ := strconv.Atoi(pnum)
//	pnum = strconv.Itoa(y + 1)
//	fmt.Println(string(pnum))
	fmt.Println(fid)

	post = string(bucket.Get([]byte(fid)))
//        fmt.Println(string(post))
	return nil
	})
	return post
}


func getNextPnum(tid string) int {
	dlist, err := ioutil.ReadDir(basePath+"/store/"+tid)
	if err !=nil {
		fmt.Println(err)
	}
	last := 0
	for k, err1 := range dlist {
		this, err2 := strconv.Atoi(dlist[k].Name())
		if err2 !=nil {
		fmt.Println(err1, err2)
		}
		if this > last {
			last = this
		}
		
//		fmt.Println(dlist[k].Name())
	}
//	last,_ := strconv.Atoi(dlist[len(dlist)-1].Name())
	return last + 1
}

func getTopic(tid string) string {
	top := ""
	i := 0
	db, err := bolt.Open(Database, 0644, nil)
	if err != nil {
		log.Fatal(err, "444")
	}
	defer db.Close()
	
	for i = 0; i < 31; i += 1 {
		pnum := strconv.Itoa(i)
		err := db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte("tracks"))
			if bucket == nil {
				fmt.Println("Bucket <tracks> not found in Db !", err)
			}

			post := string(bucket.Get([]byte(tid+"_"+pnum)))
//        fmt.Println(string(post))
			if len(post) > 0 {
				top = top + "<div id=\""+strconv.Itoa(i)+"\" class=\"r\">"+Markdown(post)+"</div></div><br>"
			}
//			fmt.Println("getTopic("+tid+"): OK")

			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	return top
}


func getFlat(fid string) string {
	page := ""
//	i := 0
	db, err := bolt.Open(Database, 0644, nil)
	if err != nil {
		log.Fatal(err, "777")
	}
	defer db.Close()
	
//	for i = 0; i < 31; i += 1 {
//		pnum := strconv.Itoa(i)
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("flats"))
		if bucket == nil {
			fmt.Println("Bucket <flats> not found in Db !", err)
		}

		chunk := string(bucket.Get([]byte(fid)))
//		fmt.Println("page in getFlat() = ", chunk)
		if len(chunk) > 0 {
			page = chunk
//			page = page + "<div id=\""+fid+"\" class=\"r\">"+Markdown(chunk)+"</div></div><br>"
		}
		fmt.Println("getFlat("+fid+"): OK")
		

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
//	fmt.Println("page_2 in getFlat() = ", page)
	return page
}

func getStaticTopic(tid string) string {
	topic, err := ioutil.ReadFile(filepath.Join(basePath+"/store2/"+tid))
	if err != nil {
		fmt.Println("getStaticTopic(): ", err)
		return "e"
	}
	fmt.Println("getStaticTopic("+tid+"): OK")
	return string(topic)
}

func getStaticPrivTopic(tid string) string {
	topic, err := ioutil.ReadFile(filepath.Join(basePath+"/store2/"+tid))
	if err != nil {
		fmt.Println("getStaticPrivTopic(): ", err)
		return "e"
	}
	fmt.Println("getStaticPrivTopic("+tid+"): OK")
	return string(topic)
}

func craftStaticTopic(tid string) string {
	topic := ""
	i := 0
	for i = 0; i < 31; i += 1 {
		post := ReadPost(tid, strconv.Itoa(i))
//		fmt.Println(post)
		topic = topic + post
	}
	file, ferr := os.Create("store2/"+tid)
	if ferr != nil {
		fmt.Println(ferr)
	}
	n, err := WritePost(Markdown(topic), file)
//	n = int64(nint)
	if err != nil {
		fmt.Println(err, n, "==bytes")
	}
	file.Close()
	return "Ok"
}

func craftStaticPrivTopic(tid string) string {
	top := ""
	i := 0
	for i = 0; i < 31; i += 1 {
		post := ReadPost(tid, strconv.Itoa(i))
		if len(post) == 0 { continue }
		metapost := ReadPost(tid, strconv.Itoa(0-i))
//		fmt.Println(post)
		top = top + "<div id=\""+strconv.Itoa(i)+"\" class=\"r\">"+metapost +"\n"+Markdown(post)+"</div></div><br>"
	}
	file, ferr := os.Create("store2/"+tid)
	if ferr != nil {
		fmt.Println(ferr)
	}
	n, err := WritePost(top, file)
//	n = int64(nint)
	if err != nil {
		fmt.Println(err, n, "==bytes")
	}
	file.Close()
	return "Ok"
}

func getPrivTopic(tid string) string {
	top := "<br>"
	i := 0
	for i = 0; i < 31; i += 1 {
		post := ReadPost(tid, strconv.Itoa(i))
		metapost := ReadPost(tid, strconv.Itoa(0-i))
//		fmt.Println(post)
		if len(post) > 0 {
			top = top + "<div id=\""+strconv.Itoa(i)+"\" class=\"r\">"+metapost +"\n"+Markdown(post)+"</div></div><br>"
		}
	}
	fmt.Println("getPrivTopic("+tid+"): OK")
//	return Markdown(top)
	return top
}

//UF2 method
func getTopic2(page string, pageint int) string {
	top := "<br>"
//pint := strconv.Atoi(page)
	i := (pageint - 1) * 20 + 1
	imax := i + 19
	for i = 0; i < imax; i += 1 {
		post := ReadPost(page, strconv.Itoa(i))
//		metapost := ReadPost(tid, strconv.Itoa(0-i))
//		fmt.Println("post===", post)
		if len(post) > 0 {
			top = top + "<div id=\""+strconv.Itoa(i)+"\" class=\"r\">"+Markdown(post)+"</div></div><br>"
		}
	}
	fmt.Println("getTopic2("+page+") OK")
//	return Markdown(top)
	return top
}


//UF2 method
func getTopic3(page, post string) string {
	top := "<br>"
	pint, _ := strconv.Atoi(post)
	i := 0
//	i := (pint - 1) * 20 + 1
////	page := pint/20 + 1
//		start := pint
//		stop := pint + 19
//	}
/*
	switch pint%20 {
	case 1:
		start := pint
	default:
		for {
			
//	for { */
		
	start := pint - pint%20 + 1
	stop := start + 20
	fmt.Println("start - stop= ", start, stop)
	for i = start; i < stop; i++ {
		post := ReadPost(page, strconv.Itoa(i))
//		metapost := ReadPost(tid, strconv.Itoa(0-i))
//		fmt.Println("post===", post)
		if len(post) > 0 {
			top = top + "<div id=\""+strconv.Itoa(i)+"\" class=\"r\">"+Markdown(post)+"</div></div><br>"
		}
	}
	fmt.Println("getTopic3("+page+") OK")
//	return Markdown(top)
	return top
}



func  ReadPage(tid string) string {
	content, err := ioutil.ReadFile(basePath+"/store/.wiki/"+tid)
	content2, err := ioutil.ReadFile(basePath+"/store/.wiki/"+tid+"$")
	if err == nil {
//		fmt.Print(".")
	}
	return Markdown(string(content) +"<br>"+ string(content2))
}

func  ReadRawPage(tid string) string {
	content, err := ioutil.ReadFile(basePath+"/store/.wiki/"+tid)
	if err == nil {
//		fmt.Print(".")
	}
	return string(content)
}

func  ReadChat(tid string) string {
	content, err := ioutil.ReadFile(basePath+"/store/.chat/"+tid)
//	content2, err := ioutil.ReadFile(basePath+"/store/.wiki/"+tid+"$")
	if err == nil {
//		fmt.Print(".")
	}
	return string(content)// +"<br>"+ string(content2))
}


func  ReadRecChanges(day string) string {
	content, err := ioutil.ReadFile(basePath+"/log/"+day)
	if err == nil {
	}
	return string(content)
}

func getLastPage(tid string) string {
	tid = strings.Replace(tid, "z", "", -1)
	cnt := 0
	for {
		_, err := os.Stat("store/"+tid+"x/1")
		if err != nil {
			fmt.Println("getLastPage(): ", err)
			break
		}
		tid += "x"
		cnt++
	}
	return tid
}

func getLastPage2(tid string) string {
//	tid = strings.Replace(tid, "z", "", -1)
	cnt := 0
	for {
		_, err := os.Stat("store/"+tid+"x/1")
		if err != nil {
			fmt.Println("getLastPage2(): ", err)
			break
		}
		tid += "x"
		cnt++
	}
	return strconv.Itoa(cnt)
}

func Format(p string) string {
//	s := string('"')

//	p = strings.Replace(p, `"`, "''", -1)

//	p = strings.Replace(p, `"`, `"`, -1)

//	pb := []byte(p)
//	for _, v := range pb {
//		if v == 0x34 { v = 0x36 }
//	}
//	return string(pb)
	return p
}


