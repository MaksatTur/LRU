package api

import (
	"encoding/json"
	"fmt"
	"log"
	"lru/memcache"
	"net/http"
	"strconv"
)

//Data is
type Data struct {
	Number1 int `json:"number1"`
	Number2 int `json:"number2"`
}

//Service is
type Service struct {
	URL string
	LRU *memcache.LRU
}

//NewService is
func NewService(url string, cap int) *Service {
	return &Service{
		URL: url,
		LRU: memcache.NewLRU(cap),
	}
}

func (s *Service) mainHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static")
}

func (s *Service) panicMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

func (s *Service) sumHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)
		var d Data
		err := decoder.Decode(&d)
		if err != nil {
			panic(err)
		}
		result := s.calculateSUM(d.Number1, d.Number2, s.LRU)
		js, err := json.Marshal(result)
		if err != nil {
			panic(err)
		}
		s.print()
		w.Write(js)
	} else {
		w.Write([]byte("Method is in valid"))
	}
}

func (s *Service) getValues(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		responce, _ := json.Marshal(s.LRU.Items)
		w.Header().Set("Content-Type", "application/json")
		w.Write(responce)
	} else {
		w.Write([]byte("Method is in valid"))
	}
}

func (s *Service) calculateSUM(num1 int, num2 int, lru *memcache.LRU) int {
	key := strconv.Itoa(num1) + "{+}" + strconv.Itoa(num2)
	sum := num1 + num2
	lru.Set(key, sum)
	return sum
}

func (s *Service) print() {
	for key := range s.LRU.Items {
		if element := s.LRU.GetElementValue(key); element != nil {
			fmt.Printf("Key = %s Value = %v\n", key, element)
		} else {
			fmt.Printf("Key = %s Value = %s\n", key, "nil")
		}
	}
	fmt.Println("-------------------------------------")
}

//StartService will started service
func (s *Service) StartService() {
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/", s.mainHandler)
	apiMux.HandleFunc("/sum", s.sumHandler)
	apiMux.HandleFunc("/getValue", s.getValues)

	mainHandler := s.panicMiddlware(apiMux)

	fmt.Printf("Service is started on %s\n", s.URL)
	log.Fatal(http.ListenAndServe(s.URL, mainHandler))
}
