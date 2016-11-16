
package main

import (
	"fmt"
	bf "github.com/russross/blackfriday"
	"os"
	"strings"
	"strconv"
	"encoding/hex"
//	"bitbucket.org/ede/
	"sha3"
//	"github.com/jzelinskie/
	"whirlpool"
	"hash"
	"crypto/rand"
	"math/big"
)

const ShortLen = 27

type short_index uint

func Markdown(source string) string {
	// set up the HTML renderer
	flags := 0
	flags |= bf.HTML_USE_SMARTYPANTS
	flags |= bf.HTML_SMARTYPANTS_FRACTIONS
	renderer := bf.HtmlRenderer(flags, "", "")

	// set up the parser
	ext := 0
	ext |= bf.EXTENSION_NO_INTRA_EMPHASIS
	ext |= bf.EXTENSION_TABLES
	ext |= bf.EXTENSION_HARD_LINE_BREAK
	ext |= bf.EXTENSION_LAX_HTML_BLOCKS
	ext |= bf.EXTENSION_FENCED_CODE
	ext |= bf.EXTENSION_AUTOLINK
	ext |= bf.EXTENSION_STRIKETHROUGH
	ext |= bf.EXTENSION_SPACE_HEADERS

	return string(bf.Markdown([]byte(source), renderer, ext))
}

func keccak224(b []byte) []byte {
	var x hash.Hash = sha3.New224()
	x.Write(b)
	return x.Sum(nil)
}

func keccak512(b []byte) []byte {
	var x hash.Hash = sha3.New512()
	x.Write(b)
	return x.Sum(nil)
}

func whirl(g []byte) []byte {
	var u hash.Hash = whirlpool.New()
	u.Write(g)
	return u.Sum(nil)
}

//func scryptPhash(d []byte) string {
//	return hex.EncodeToString(
//	q, err := scrypt.Key(keccak512(d), keccak224(d), 16384, 9, 2, 32)
//	if err != nil {
//		fmt.Println(err)
//		panic("Scrypt error !")
//	}
//	return hex.EncodeToString(q)
//}

func whirlPhash(d []byte) string {

	return hex.EncodeToString(whirl([]byte(hex.EncodeToString(whirl(d)))))
}

func hasUsername2(name string) bool {
	for k, _ := range qz {
		if k == name { return true }
	}
	return false
}

func privIndex(priv string) int {
	for k, v := range privTids {
		if v == priv { return k }
	}
	return -1
}


func addUser(login, phash, handl, raddr string) bool {
	addon := NewUpid()[:14]
	dumpZcreds("zcreds.go_"+addon)
	qz[login] = phash
	hz[login] = handl
	file, err := os.Create(basePath+"/zreg0/"+login+"_"+addon)
	if err != nil {
		fmt.Println("addUser(): os.Create(): ", err)
		return true
	}
	_, werr := WritePost(phash+"\n"+raddr, file)
	if werr != nil {
		fmt.Println("WritePost(): ", err)
		file.Close()
		return true
	}
	file.Close()
	dumpZcreds("zcreds.go")
	return true
}

func dumpZcreds(dest string) {
	head := "\npackage main\n\nfunc init() {\n\n"
	body := ""
	tail := "\n}"
	for login, phash := range qz {
		body += strings.Replace("\tqz[*"+login+"*] = *"+phash+"*\n", "*", string('"'), -1)
		body += strings.Replace("\thz[*"+login+"*] = *"+hz[login]+"*\n", "*", string('"'), -1)
	}
	creds := head + body + tail+"\n\n"
	file, err := os.Create("/0/uf2/"+dest)
	_, werr := WriteMetaPost(creds, file)
	if err != nil || werr != nil {
		fmt.Println(err, werr)
	}
//	fmt.Print(creds)
	file.Close()
}

func dumpWatch(dest string) {
	head := "\npackage main\n\n\nfunc init() {\n\n"//var lastPosts = make(map[string]int, 100000)\n\nfunc init() {\n\n"
	body := ""
	tail := "\n}"
	for td, lst := range lastPages {
		body += strings.Replace("\tlastPages[*"+td+"*] = "+strconv.Itoa(lst)+"\n", "*", string('"'), -1)
	}
	body += "\n"
	for tid, last := range lastPosts {
		body += strings.Replace("\tlastPosts[*"+tid+"*] = "+strconv.Itoa(last)+"\n", "*", string('"'), -1)
	}
	body += "\n"
	for k, v := range wl {
//		fmt.Println("k=", k)
		for g, f := range v {
//			fmt.Println(g, f)
			body += strings.Replace("\twl[*"+k+"*][*"+g+"*] = "+strconv.Itoa(f)+"\n", "*", string('"'), -1)
		}
	}
	watch := head + body + tail+"\n\n"
	file, err := os.Create("/0/uf2/"+dest)
	_, werr := WriteMetaPost(watch, file)
	if err != nil || werr != nil {
		fmt.Println(err, werr)
	}
	fmt.Print(watch)
	file.Close()
}


func NewUpid() (shortid string) {
	const symbols = "ABCDEFGHJKLMNoPqRSTUVWXYZ123456789"
	for i := 0; i < ShortLen; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(symbols))))
		if err != nil {
			shortid = ""
			return
		} else {
			index := n.Int64()
			shortid += string(symbols[index])
		}
	}
	return
}

