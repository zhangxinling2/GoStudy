package main

import (
	"GoStudy/dataStore/fatRank"
	"GoStudy/httpServer/httpPratice/frinterface"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	var rankServer frinterface.ServeInterface = NewFatRateRank()

	mux := http.NewServeMux()
	mux.Handle("/registry", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//method不能为非post
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		//注册的不能为空
		if r.Body == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		//golang中执行request.Body.close是为了释放连接资源，防止内存泄漏。如果不关闭request.Body，那么客户端可能无法重用持久的TCP连接来发送后续的请求。
		//因此，在读取完request.Body后，最好使用defer语句来关闭它。
		defer r.Body.Close()
		//读出内容后解码
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("无法读取数据,%v", err)))
			return
		}
		pi := fatRank.PersonalInformation{}
		if err = json.Unmarshal(content, &pi); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("无法解析数据,%v", err)))
			return
		}
		if err = rankServer.RegisterPersonInformation(&pi); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("无法注册信息,%v", err)))
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	mux.Handle("/personinfo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//method不能为非post
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		//注册的不能为空
		if r.Body == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		//golang中执行request.Body.close是为了释放连接资源，防止内存泄漏。如果不关闭request.Body，那么客户端可能无法重用持久的TCP连接来发送后续的请求。
		//因此，在读取完request.Body后，最好使用defer语句来关闭它。
		defer r.Body.Close()
		//读出内容后解码
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("无法读取数据,%v", err)))
			return
		}
		pi := fatRank.PersonalInformation{}
		if err = json.Unmarshal(content, &pi); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("无法解析数据,%v", err)))
			return
		}
		if fr, err := rankServer.UpdatePersonInformation(&pi); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("无法更改信息,%v", err)))
			return
		} else {
			w.WriteHeader(http.StatusOK)
			data, _ := json.Marshal(fr)
			w.Write(data)
		}

	}))
	mux.Handle("/rank", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//method不能为非Get
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		name := r.URL.Query().Get("name")
		if name == "" {
			w.Write([]byte("name未设置"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if fr, err := rankServer.GetFatrate(name); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("无法获取排行数据,%v", err)))
			return
		} else {
			w.WriteHeader(http.StatusOK)
			data, _ := json.Marshal(fr)
			w.Write(data)
		}
	}))
	mux.Handle("/rankTop", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//method不能为非Get
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if frTop, err := rankServer.GetTop(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("无法获取排行数据,%v", err)))
			return
		} else {
			w.WriteHeader(http.StatusOK)
			data, _ := json.Marshal(frTop)
			w.Write(data)
		}
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
