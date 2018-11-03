package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"gitlab.com/irezaa/calendar/src/model"
	"bytes"
	"github.com/pkg/errors"
	"time"
)

func startCrawler(app *App) {
	client := &http.Client{}

	resultChannel := make(chan MonthCrawled)

	var allJobsCount = (1410 - 1370) * 12

	var doneJobs = 0

	print(allJobsCount)

	go func() {
		for year := 1300; year <= 1450; year++ {
			for month := 1; month <= 12; month++ {
				go getByMonth(app, client, year, month, resultChannel)
				time.Sleep(2 * time.Second)
			}
		}
	}()

	for {
		select {
		case j := <-resultChannel:
			doneJobs++
			fmt.Printf("\n")

			if j.isFailed {
				fmt.Printf("Failed Job : URL : %s - Error : %s ", j.url, j.err)
			} else {
				fmt.Printf("Done Job : Jalali Year : %d - jalaliMonth : %d - eventsCount : %d", j.jalaliYear, j.jalaliMont, j.eventsCount)
			}

		}
	}

	print("Done :D")
}

func getByMonth(app *App, client *http.Client, jalaliYear int, jalaliMonth int, resultChannel chan MonthCrawled) {
	crawlStruct := MonthCrawled{
		isFailed:    false,
		jalaliYear:  jalaliYear,
		jalaliMont:  jalaliMonth,
		eventsCount: 0,
	}

	url := fmt.Sprintf("http://www.time.ir/fa/event/api/getallevents/0/%d/%d/a1629a82eef874771f614cbb248910db246add383d1c89a14f1f8b78023b789d6ae84d83a897a25fcecff336078ad55b", jalaliYear, jalaliMonth)

	crawlStruct.url = url

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Mac; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1")
	resp, err := client.Do(req)

	if err != nil {
		fmt.Print(err)
		crawlStruct.isFailed = true
		crawlStruct.err = err
	} else {
		body, _ := ioutil.ReadAll(resp.Body)

		r := bytes.NewReader(body)
		var events []model.Event

		jsonDecodeError := json.NewDecoder(r).Decode(&events)

		if jsonDecodeError == nil {
			if events != nil && len(events) > 0 {
				for _, event := range events {
					event.EventType = model.EventTypeGeneral
					err := app.EventRepository.UpdateOrInsert(event)

					if err != nil {
						fmt.Println(err)
					}
				}

				crawlStruct.isFailed = false
				crawlStruct.eventsCount = len(events)
			} else {
				if events == nil {
					crawlStruct.err = errors.New("Events array is nil")
				} else if len(events) == 0 {
					crawlStruct.err = errors.New("Events array count is 0")
				}

				crawlStruct.isFailed = true
			}
		} else {
			fmt.Printf("Response Status : %d", resp.StatusCode)
			crawlStruct.isFailed = true
			crawlStruct.err = jsonDecodeError
		}
	}

	resultChannel <- crawlStruct

	defer resp.Body.Close()
}

type MonthCrawled struct {
	url         string
	isFailed    bool
	err         error
	eventsCount int
	jalaliMont  int
	jalaliYear  int
}
