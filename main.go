package main

import (
	"encoding/json"
	"fmt"
	"github.com/jD91mZM2/stdutil"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)

const URL = "https://www.cleverbot.com/getreply"

var conv1 string
var conv2 string
var stopped bool

var turn bool
var reply string

func main() {
	fmt.Print("Cleverbot key: ")
	KEY := stdutil.MustScanTrim()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		<-c
		stopped = true
	}()

	for i := 0; i < 3; i++ {
		fmt.Println()
	}

	for {
		if stopped {
			break
		}
		time.Sleep(3 * time.Second)
		turn = !turn
		if stopped {
			break
		}

		turnStr := "1"
		if !turn {
			turnStr = "2"
		}
		fmt.Print("CleverBot " + turnStr + ": ")

		v := url.Values{}
		v.Set("key", KEY)
		v.Set("input", reply)

		if turn {
			v.Set("cs", conv1)
		} else {
			v.Set("cs", conv2)
		}

		res, err := http.Get(URL + "?" + v.Encode())
		if err != nil {
			stdutil.PrintErr("Could not make request", err)
			break
		}
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			stdutil.PrintErr("Could not get data from request", err)
			break
		}

		replyMap := make(map[string]string)
		err = json.Unmarshal(body, &replyMap)
		if err != nil {
			stdutil.PrintErr("Could not parse JSON", err)
			break
		}

		if turn {
			conv1 = replyMap["cs"]
		} else {
			conv2 = replyMap["cs"]
		}

		reply = replyMap["output"]
		fmt.Println(reply)
	}
}
