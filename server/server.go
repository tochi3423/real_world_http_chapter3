package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

func handler(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	fmt.Println(string(dump))
	fmt.Fprintf(w, "<html><body>hello</body></html>¥n")
}
func main() {
	var httpServer http.Server
	http.HandleFunc("/", handler)
	log.Println("start https listening: 18888")
	httpServer.Addr = ":18888"
	log.Println(httpServer.ListenAndServe())
}

// func main() {
//  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//      cookie := &http.Cookie{
//          Name:  "test_cookie",
//          Value: "hello, world",
//      }
//      http.SetCookie(w, cookie)
//      receivedCookies := r.Cookies()
//      for _, cookie := range receivedCookies {
//          log.Printf("Received cookie: %s - %s\n", cookie.Name, cookie.Value)
//      }
//  })
//  log.Fatal(http.ListenAndServe(":18888", nil))
// }
