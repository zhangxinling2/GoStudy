package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/user", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello user"))
	}))
	mux.Handle("/rank", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello rank"))
	}))
	mux.Handle("/history/xiaoqiang", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello xiaoqiang"))
	}))
	mux.Handle("/history", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mq := r.URL.Query()
		name := mq.Get("name")
		w.Write([]byte(fmt.Sprintf("%s 的rank为", name)))
	}))
	http.ListenAndServe(":8080", mux)
	// http.ListenAndServe(":8088", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	// time.Sleep(2 * time.Second)
	// 	// w.Write([]byte("hello 你好"))

	// 	//request读的时候要使用POST
	// 	if r.Body == nil {
	// 		w.Write([]byte("no body"))
	// 		return
	// 	}
	// 	data, _ := ioutil.ReadAll(r.Body)
	// 	defer r.Body.Close()
	// 	encoded := base64.StdEncoding.EncodeToString(data)
	// 	w.Write(append(data, []byte(encoded)...))

	// 	// qp := r.URL.Query()
	// 	// data, _ = json.Marshal(qp)
	// 	// w.Write([]byte("hello 你好" + string(data)))
	// }))
}
