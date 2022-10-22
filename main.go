package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Data struct {
	Status Status`json:"status"`
}

type Status struct {
	Water int64 `json:"water"`
	Wind  int64 `json:"wind"`
}

func GetFile() []byte {
	file, err := os.Open("data.json")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	return data
}

func Changedata(data []byte) []byte {
	var d Data
	err := json.Unmarshal(data, &d)
	if err != nil {
		fmt.Println(err)
	}
	rand.Seed(time.Now().UnixNano())
	d.Status.Water = rand.Int63n(100)
	d.Status.Wind = rand.Int63n(100)
	data, err = json.Marshal(d)
	if err != nil {
		fmt.Println(err)
	}
	return data
}

func changeValueinFile(data []byte) {
	err := ioutil.WriteFile("data.json", data, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func updateFile(data []byte) {
	data = Changedata(data)
	err := ioutil.WriteFile("data.json", data, 0644)
	if err != nil {
		fmt.Println(err)
	}
}


func autoChange(){
	data := GetFile()
	for range time.Tick(15 * time.Second) {
		updateFile(data)
	}
}

func main() {
	go autoChange()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := GetFile()
		w.Write(data)
	})
	http.ListenAndServe(":8080", nil)
}