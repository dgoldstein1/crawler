package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// posts possible new edges to GRAPH_DB_ENDPOINT
func addNeighbors(curr int, neighborIds []int) (resp GraphResponseSuccess, err error) {
	// POST new neighbors to db
	jsonValue, _ := json.Marshal(map[string][]int{
		"neighbors": neighborIds,
	})
	url := os.Getenv("GRAPH_DB_ENDPOINT") + "/edges"
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	q.Add("node", strconv.Itoa(curr))
	req.URL.RawQuery = q.Encode()

	// return the result of the POST request
	client := http.Client{
		Timeout: timeout,
	}
	res, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	// assert response is 200
	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return resp, err
		}
		errResp := GraphResponseError{}
		err = json.Unmarshal(body, &errResp)
		if err != nil {
			return resp, err
		}
		// fails with error
		return resp, errors.New(errResp.Error)
	}

	// 200 level response, continue as normal
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return resp, err
	}
	resp = GraphResponseSuccess{}
	err = json.Unmarshal(body, &resp)
	return resp, err
}

// gets wikipedia int id from article url
func getArticleIds(articles []string) (resp TwoWayResponse, err error) {
	// create array of entries
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(articles)
	// post to endpoint
	url := os.Getenv("TWO_WAY_KV_ENDPOINT") + "/entries"
	req, _ := http.NewRequest("POST", url, b)
	req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	q.Add("muteAlreadyExistsError", "true")
	req.URL.RawQuery = q.Encode()
	client := http.Client{
		Timeout: timeout,
	}
	res, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	// read out response
	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return resp, err
		}
		errResp := GraphResponseError{}
		err = json.Unmarshal(body, &errResp)
		if err != nil {
			return resp, err
		}
		// fails with error
		return resp, errors.New(errResp.Error)
	}
	// succesful request
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return resp, err
	}
	resp = TwoWayResponse{}
	err = json.Unmarshal(body, &resp)
	return resp, err
}

// connects to given databse and initializes scraper
func ConnectToDB() error {
	resp, err := http.Get(os.Getenv("GRAPH_DB_ENDPOINT"))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	return err
}
