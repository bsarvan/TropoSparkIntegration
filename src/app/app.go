package main

import (
  "log"
  . "global"
  . "speechdb"
  "net/http"
  "html/template"
  "os"
  "os/exec"	
  "fmt"
  "io"
  . "gspeechimpl"
)

type templatedata struct {
    Sparkid string
    Mobile  string
    Search  string
}

func init() {
    log.Println("Initialising Database...")
    LoadData(GlobalData)  
}

func handler(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("f1.html")
    t.Execute(w,nil)
}

func dummy() (status bool) {
    return false
}

func sparkadd(w http.ResponseWriter, r *http.Request) {
    log.Println("SparkID to be added", r.FormValue("sparkid"))
    log.Println("Mobile to be added", r.FormValue("mobile"))
    var A1 templatedata 
    returnvalue := dummy() 
    if returnvalue == false {
        log.Println("In Here")
        t, _ := template.ParseFiles("f2.html")
        A1.Sparkid = r.FormValue("sparkid")
        A1.Mobile = r.FormValue("mobile")
        A1.Search = "1234"
        t.Execute(w, A1) 
    }
    
}

func Uploadfile(w http.ResponseWriter, r *http.Request) {
    file, handler, err := r.FormFile("filename")
    if err != nil {
        log.Println(err)
        return
    } 
    
    log.Println("In function Uploadfile")  
    defer file.Close()
    fmt.Fprintf(w, "%v", handler.Header)
    log.Println("Saving file ", handler.Filename)
    f, err := os.OpenFile("/home/ec2-user/GoLang/wav/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        log.Println(err)
        return
    }   
    defer f.Close()
    io.Copy(f, file)
    //Convert wav to flac format
    SoxCmd := exec.Command("sox", "/home/ec2-user/GoLang/wav/"+handler.Filename, "/home/ec2-user/GoLang/wav/"+handler.Filename+".flac")

    SoxOut, err := SoxCmd.Output()
    if err != nil {
        panic(err)
	return
    }   
    log.Println("Processed file to flac")
    log.Println(string(SoxOut))
    output := Processflac("/home/ec2-user/GoLang/wav/" + handler.Filename + ".flac")
    log.Println("Processed string ", output)
    
return
}

func addmatch(w http.ResponseWriter, r *http.Request) {
    log.Println("NEW:- addmatch to be added", r.FormValue("search"))
    log.Println("NEW:- mobilenumber to be added", r.FormValue("Mobile"))

}

func main() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/addsparkid", sparkadd)
    http.HandleFunc("/addmatch", addmatch) 
    /* Recording file from Tropo */
    http.HandleFunc("/uploadmedia", Uploadfile)	
    http.ListenAndServe(":8080", nil)
}


