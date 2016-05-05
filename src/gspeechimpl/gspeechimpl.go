/****************************************************************
* Author : Bharat Sarvan 
* 30/04/2016
* GSpeech implementation
*****************************************************************/

package gspeechimpl

import (
	"encoding/json"
	. "gspeech"
	"log"
	"os"
	"strings"
)

//Result JSON
type Trans struct {
	Result       []Alternative
	Result_index int
}
type Alternative struct {
	Alternative []transcript
	Stability   float64
	Final       bool
}

type transcript struct {
	Transcript string
}

//API key
var Key string

var Googlesettings struct {
	Key string `json:"Key"`
}

//Load key from google.json
func init() {
	configFile, err := os.Open("google.json")
	if err != nil {
		log.Println("auth.json file error")
	}
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&Googlesettings); err != nil {
		log.Println("parsing google.json config file", err.Error())
	}
	Key = Googlesettings.Key
	log.Println("Google Key", Key)
}

//Process flac files
func Processflac(filename string) (result string) {
	log.Println("Start processing the filename= \n", filename)

	//Invoke google voice library
	messages, err := StartProcessing(filename, Key)
	if err == nil {

		// JSON contains multiple results
		lines := strings.Split(messages, "\n")
		//fmt.Printf("Length is %d",len(lines))

		for icnt := range lines {
			// Ignore empty result
			if (strings.Contains(lines[icnt], `{"result":[]}`) == false) && (len(lines[icnt]) > 0) {
				// Decode the json object
				u := &Trans{}

				err := json.Unmarshal([]byte(lines[icnt]), &u)
				if err != nil {
					panic(err)
				} else {
					// look for Final value
					if u.Result[0].Final == true {
						//fmt.Printf("now processing : %s \n",lines[icnt])
						log.Println("%s \n", u.Result[0].Alternative[0].Transcript)
						log.Println("%t \n", u.Result[0].Final)
						return u.Result[0].Alternative[0].Transcript
					}
				}
			}
		}
	}
	return ""
}
