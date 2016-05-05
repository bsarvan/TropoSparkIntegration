/****************************************************************
* Author : Bharat Sarvan 
* 30/04/2016
* Google cloud speech API library
*****************************************************************/

package gspeech

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type ProcessedData struct {
	filename string
	data     string
	err      error
}

func readfile(filename string) (buffer []byte) {

	file, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return
	}

	var size int64 = info.Size()
	rbytes := make([]byte, size)
	bufferd := bufio.NewReader(file)
	_, err = bufferd.Read(rbytes)
	return rbytes
}

/*Main function to handle speech to text*/
func StartProcessing(file string, key string) (message string, err error) {
	// buffered data
	ch := make(chan *ProcessedData, 1)
	s2 := rand.NewSource(time.Now().UnixNano())
	r2 := rand.New(s2)

	//Create pair for the API
	pair1 := r2.Int63n(1234567891234567)
	pair := strconv.FormatInt(pair1, 10)

	//POST method
	go func(file string, key string, pair string) {
		upURL := fmt.Sprintf("https://www.google.com/speech-api/full-duplex/v1/up?key=%s&pair=%s&output=json&lang=en-US&pFilter=2&maxAlternatives=1&app=chromium&continuous&interim", key, pair)
		client := &http.Client{}
		b := bytes.NewBuffer(readfile(file))

		req, err := http.NewRequest("POST", upURL, b)
		if err != nil {
			log.Fatalln(err)
		}

		req.Header.Add("Content-Type", "audio/x-flac; rate=8000")
		req.Header.Add("Transfer-Encoding", "chunked")
		resp, err := client.Do(req)

		resp.Body.Close()
	}(file, key, pair)

	//GET method for transcription data
	go func(key string, pair string) {
		downurl := fmt.Sprintf("https://www.google.com/speech-api/full-duplex/v1/down?key=%s&pair=%s&output=json", key, pair)

		client := &http.Client{}

		req, err := http.NewRequest("GET", downurl, nil)
		if err != nil {
			log.Fatalln(err)
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		data, err := ioutil.ReadAll(resp.Body)

		//fmt.Printf("\n\n\n GETTING DATA:- %s \n\n\n",data)
		ch <- &ProcessedData{file, string(data), err}
	}(key, pair)

	result := <-ch
	return result.data, result.err
}
