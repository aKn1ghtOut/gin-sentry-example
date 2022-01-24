package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

type GetPostsResponseItem struct {
	Title       string `json:"title"`
	Description string `json:"text"`
	Type        string `json:"type"`
}

type GetPostsResponse []GetPostsResponseItem

func getPostsCall(sentryCtx context.Context) GetPostsResponse {
	// sentry span init for function, set parent to context from calling function
	defer sentry.Recover()
	span := sentry.StartSpan(sentryCtx, "[GetPostsCall]")
	defer span.Finish()

	var res GetPostsResponse

	httpClient := &http.Client{}
	span2 := sentry.StartSpan(span.Context(), "[HTTP] Get latests stories") // Start nested span to monitor first http request
	req, err := http.NewRequest("GET", "https://hacker-news.firebaseio.com/v0/newstories.json?orderBy=\"$key\"&limitToFirst=10", nil)
	if err != nil {
		fmt.Print("Error 1:", err.Error())
		sentry.CaptureException(err)
	}
	req.Close = true
	resp, err2 := httpClient.Do(req)
	if err2 != nil {
		fmt.Print("Error 2:", err2.Error())
		sentry.CaptureException(err2)
	}
	body, err3 := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	span2.Finish() // Request workflow finished; end span
	if err3 != nil {
		fmt.Print("Error 2: ", err3.Error())
		sentry.CaptureException(err3)
	}

	fmt.Printf("Body : %s", string(body))

	var ids []int
	json.Unmarshal(body, &ids)

	fmt.Print(ids[0:10])

	span3 := sentry.StartSpan(span.Context(), "[HTTP] Get data for stories") // Start span to monitor for loop execution time
	for i := range ids {
		id := ids[i]

		urlLink := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%v.json", id)
		iResp, err := http.Get(urlLink)
		if err != nil {
			sentry.CaptureException(err)
			continue
		}

		resB, _ := ioutil.ReadAll(iResp.Body)
		var temp GetPostsResponseItem
		json.Unmarshal(resB, &temp)
		res = append(res, temp)
	}
	span3.Finish()

	return res
}

func GetPostsHandler(c *gin.Context) {
	// Sentry transaction initialization
	defer sentry.Recover()
	span := sentry.StartSpan(c, "[Getposts]", sentry.TransactionName("GETPOSTS"), sentry.ContinueFromRequest(c.Request)) // Take trace id from request header
	defer span.Finish()

	response := getPostsCall(span.Context()) // Call nested function; Could be DAO methods

	c.JSON(200, response)
}
