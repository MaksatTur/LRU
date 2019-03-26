package main

import (
	"encoding/json"
	"fmt"
	"log"
	"lru/memcache"
	"net/http"
	"strconv"
)

const (
	CAPACITY = 2
	PORT     = ":3000"
)

var lru = memcache.NewLRU(CAPACITY)

type Data struct {
	Number1 int `json:"number1"`
	Number2 int `json:"number2"`
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static")
}

func calculateSUM(num1, num2 int) int {
	key := strconv.Itoa(num1) + "{+}" + strconv.Itoa(num2)
	sum := num1 + num2
	lru.Set(key, sum)
	return sum
}

func sumHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("sumHandler", r.URL.Path)
	if r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)
		var d Data
		err := decoder.Decode(&d)
		if err != nil {
			panic(err)
		}
		result := calculateSUM(d.Number1, d.Number2)
		js, err := json.Marshal(result)
		if err!=nil{
			panic(err)
		}
		print()
		w.Write(js)
	}else {
		w.Write([]byte("Method is in valid"))
	}
}

func panicMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println("panicMiddlware", r.URL.Path)
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("recovered", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func print()  {
	for key, _ := range lru.Items{
		if element := lru.GetElementValue(key); element != nil{
			fmt.Printf("Key = %s Value = %v\n", key, element)
		} else {
			fmt.Printf("Key = %s Value = %s\n", key, "nil")
		}
	}
	fmt.Println("-------------------------------------")
}

func main() {
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/", mainHandler)
	apiMux.HandleFunc("/sum", sumHandler)

	mainHandler := panicMiddlware(apiMux)

	fmt.Printf("Service is started on port %s\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, mainHandler))
	//lru.Set("1",1)
	//lru.Set("2",2)
	//element := lru.Queue.Back()
	//fmt.Println("Back = ", element.Value.(*memcache.Item).Value)
	//lru.Set("3",3)
	//element = lru.Queue.Back()
	//fmt.Println("Back = ", element.Value.(*memcache.Item).Value)
	//print()
	//lru.Set("4",4)
	//element = lru.Queue.Back()
	//fmt.Println("Back = ", element.Value.(*memcache.Item).Value)
}
