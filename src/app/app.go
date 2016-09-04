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
  . "spark"
  . "gspeechimpl"
  "strings"
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
    t,err := template.ParseFiles("src/f1.html")
    if err!=nil {
	http.Error(w,err.Error(),http.StatusInternalServerError)
	return
    }
    t.Execute(w,nil)
}

func dummy() (status bool) {
    return false
}

func sparkadd(w http.ResponseWriter, r *http.Request) {
    log.Println("SparkID to be added", r.FormValue("sparkid"))
    log.Println("Mobile to be added", r.FormValue("mobile"))
    search, mobile := Verifysparkid(r.FormValue("sparkid"))
    var A1 templatedata 
   
    if search != "" {
        log.Println("sparkadd:- mobile, form", mobile, r.FormValue("mobile"))
        //Spark ID already exist - Verify mobile number matches
        if mobile != r.FormValue("mobile") {
            t,_ := template.ParseFiles("src/invalidmobile.html")
            t.Execute(w, A1)
        } else {
            t, _ := template.ParseFiles("src/f2.html")
            A1.Sparkid = r.FormValue("sparkid")
            A1.Mobile = r.FormValue("mobile")
            A1.Search = "1234"
            t.Execute(w, A1)
        }
    } else {
        //New Spark ID to be added
        t, _ := template.ParseFiles("src/f2a.html")
        A1.Sparkid = r.FormValue("sparkid")
        A1.Mobile = r.FormValue("mobile")
        A1.Search = "1234"
        t.Execute(w, A1)
    }

    log.Println("Selected spardid to be added", r.FormValue("sparkid")) 
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
   
    //Verify search string matches
    //tospark := GlobalData[strings.ToLower(output)].Sparkid
    //Send spark message
    Sendspark("Voice Message - "+output, "bsarvan@cisco.com")
     
    return
}

func addmatch(w http.ResponseWriter, r *http.Request) {
    log.Println("NEW:- addmatch to be added", r.FormValue("search"))
    log.Println("NEW:- mobilenumber to be added", r.FormValue("Mobile"))
	
    if Verifysearch(r.FormValue("search")) == "" {
        returnvalue := Storerecord(r.FormValue("sparkid"), r.FormValue("mobile"), r.FormValue("search"))
        if returnvalue == true {
            GlobalData[strings.ToLower(r.FormValue("search"))] = GlobalDS{append(GlobalData[r.FormValue("search")].Mobile, r.FormValue("mobile")), append(GlobalData[r.FormValue("search")].Sparkid, r.FormValue("sparkid"))}
            t, _ := template.ParseFiles("src/f3.html")
            t.Execute(w, nil)
        } else {
            t, _ := template.ParseFiles("src/error.html")
            t.Execute(w, nil)
        }
        log.Println("selected room to be added", r.FormValue("search"))
    } else {
        t, _ := template.ParseFiles("src/duplicatesearch.html")
        t.Execute(w, nil)
    }

}

func main() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/addsparkid", sparkadd)
    http.HandleFunc("/addmatch", addmatch) 
    /* Recording file from Tropo */
    http.HandleFunc("/uploadmedia", Uploadfile)	
    http.ListenAndServe(":8080", nil)
}


