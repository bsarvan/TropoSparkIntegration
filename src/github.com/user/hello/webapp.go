package main

import (
  "net/http"
  "html/template"
  "log"  
)


func handler(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("f1.html")
    t.Execute(w,nil)
}

func sparkadd(w http.ResponseWriter, r *http.Request) {

        log.Println("SparkID to be added", r.FormValue("sparkid"))
        log.Println("Mobile to be added", r.FormValue("mobile"))
    
}

func main() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/addsparkid", sparkadd)
    
    http.ListenAndServe(":8080", nil)
}
