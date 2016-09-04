package speechdb

import (
 "database/sql"
 "encoding/json"
 "os"
 "log"
 . "global"
 _ "github.com/go-sql-driver/mysql"
)

var DBSettings struct {
    Username string `json:"userName"`
    Password string `json:"passWord"`
    DB       string `json:"db"`
    Host     string `json:"host"`
}

func init() {
    log.Println("Opening db.json")
    configFile, err := os.Open("db.json")
    if err!=nil {
        log.Fatal(err)
    }
    jsonParser := json.NewDecoder(configFile)
    if err = jsonParser.Decode(&DBSettings); err != nil {
        log.Println("parsing db.json config file", err.Error())
    }
    
    log.Println("Database User: ", DBSettings.Username)

}


func LoadData(GlobalData map[string]GlobalDS) {
    con, err := sql.Open("mysql", DBSettings.Username+":"+DBSettings.Password+"@tcp("+DBSettings.Host+":3306)/"+DBSettings.DB)
    if err!= nil {
        log.Println(err)
    }

    rows, err := con.Query("SELECT * FROM speech")
    if err == nil {
        for rows.Next() {
            var sparkid string
            var mobile string
            var search string
            err = rows.Scan(&sparkid, &mobile, &search)
            if err == nil {
                GlobalData[search] = GlobalDS{append(GlobalData[search].Mobile, mobile), append(GlobalData[sparkid].Sparkid, sparkid)}
                log.Println(sparkid)
                log.Println(mobile)
                log.Println(search)
            } else {
                log.Println(err)
            }
        }
    }
    defer con.Close()

}

func Storerecord(sparkid string, mobile string, search string) (status bool) {
    log.Println("In Storerecord", sparkid, mobile, search)
        con, err := sql.Open("mysql", DBSettings.Username+":"+DBSettings.Password+"@tcp("+DBSettings.Host+":3306)/"+DBSettings.DB)
        if err != nil {
            log.Println(err)
                log.Println(err)
        }   
    stinsert, err := con.Prepare("INSERT speech SET sparkid=?,mobile=?,search=?")
        if err == nil {
            stinsert.Exec(sparkid, mobile, search)
                return true
        } else {
            log.Println(err)
        }   
    defer con.Close()
        return false
}

//To verify whether search phrase already exist
func Verifysearch(search string) (value string) {
    con, err := sql.Open("mysql", DBSettings.Username+":"+DBSettings.Password+"@tcp("+DBSettings.Host+":3306)/"+DBSettings.DB)
    if err != nil {
        log.Println(err)
    }
    rows, err := con.Query("SELECT search FROM speech where search=?", search)
    if err == nil {
        for rows.Next() {
            var search string
            err = rows.Scan(&search)
            if err == nil {
                log.Println("Verifysearch:-", search)
                return search
            }
        }
    }
    log.Println("Verifysearch:-", search)
    return ""

}


//To verify whether sparkID already exist
func Verifysparkid(sparkid string) (search string, mobile string) {
    con, err := sql.Open("mysql", DBSettings.Username+":"+DBSettings.Password+"@tcp("+DBSettings.Host+":3306)/"+DBSettings.DB)
    if err != nil {
        log.Println(err)
    }
    rows, err := con.Query("SELECT search,mobile FROM speech where sparkid=?", sparkid)
    if err == nil {
        for rows.Next() {
            var search string
            var mobile string
            err = rows.Scan(&search, &mobile)
            if err == nil {
                log.Println("Verifysparkid:-", search)
                return search, mobile
            }
        }
    }
    log.Println("Verifysparkid:-", sparkid)
    return "", ""

}
