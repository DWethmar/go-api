package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/dwethmar/go-api/pkg/api/input"
)

var (
	// wg is used to force the application wait for all goroutines to finish before exiting.
	wg sync.WaitGroup
	// batchSize defines the batch size where each job batch will contain maximum certain number of jobs.
	batchSize = 3
	// batchInChan is a unbuffered channel that has the capacity of 1 resource(batch) slot.
	batchInChan = make(chan []input.AddContent)
)

func main() {
	host := os.Args[1]
	url, err := url.Parse(host + "/content")
	if err != nil {
		panic(err)
	}

	fmt.Println("BEGIN")
	amount := 100

	items := make([]input.AddContent, 0)
	for i := 0; i < amount; i++ {
		items = append(items, input.AddContent{
			Name: "miauw",
			Fields: input.FieldTranslations{
				"nl": {
					"attr1": "yes",
					"attr2": "yes",
				},
			},
		})
	}

	// Create 2 goroutines.
	wg.Add(2)

	go processor(1, url, batchInChan, &wg)
	go processor(2, url, batchInChan, &wg)

	// Create job batch and push them into `batchInChan` channel.
	for i := 0; i < len(items); i += batchSize {
		j := i + batchSize

		if j > len(items) {
			j = len(items)
		}

		batchInChan <- items[i:j]
	}

	// Close channel to remove the lock.
	close(batchInChan)

	// Block exiting until all the goroutines are finished.
	wg.Wait()

	fmt.Println("END")
}

func processor(id int, url *url.URL, batchInChan <-chan []input.AddContent, wg *sync.WaitGroup) {
	// As soon as the current goroutine finishes (job done!), notify back WaitGroup.
	defer wg.Done()

	b := 1

	// Listen on `batchInChan` to see if there is any resource pending in it.
	for batch := range batchInChan {
		for _, entry := range batch {
			fmt.Println("processor:", id, "batch:", b, "job:", entry.Name, "- started")

			// wait := rand.Intn(len(waiters))
			// time.Sleep(time.Duration(wait) * time.Second)
			postContent(url, entry)

			fmt.Println("processor:", id, "batch:", b, "job:", entry.Name, "- finished in")
		}

		b++
	}
}

func postContent(url *url.URL, entry input.AddContent) {
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // ignore expired SSL certificates
		},
	}
	client := &http.Client{Transport: transCfg}

	// create a request object
	jsonReq, err := json.Marshal(entry)
	if err != nil {
		log.Fatal("Error:", err)
	}

	req := &http.Request{
		Method: "POST",
		URL:    url,
		Header: map[string][]string{
			"Content-Type": {
				"application/json; charset=UTF-8",
			},
		},
		Body: ioutil.NopCloser(bytes.NewBuffer(jsonReq)),
	}

	res, err := client.Do(req)
	// check for response error
	if err != nil {
		log.Fatal("Error:", err)
	}
	// close response body
	res.Body.Close()
}
