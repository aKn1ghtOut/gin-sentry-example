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
	defer sentry.Recover()
	span := sentry.StartSpan(sentryCtx, "[GetPostsCall]")
	defer span.Finish()

	var res GetPostsResponse

	httpClient := &http.Client{}
	span2 := sentry.StartSpan(span.Context(), "[HTTP] Get latests stories")
	req, err := http.NewRequest("GET", "https://hacker-news.firebaseio.com/v0/newstories.json", nil)
	if err != nil {
		print("Error 1")
	}
	req.Close = true
	resp, err2 := httpClient.Do(req)
	if err2 != nil {
		print("Error 2")
	}
	span2.Finish()
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	fmt.Printf("Body : %s", string(body))

	var ids []int
	json.Unmarshal(body, &ids)

	fmt.Print(ids[0:10])

	incr := 0

	span3 := sentry.StartSpan(span.Context(), "[HTTP] Get data for stories")
	for i := range ids {
		id := ids[i]
		if incr > 10 {
			break
		}
		incr += 1

		urlLink := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%v.json", id)
		iResp, err := http.Get(urlLink)
		if err != nil {
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
	defer sentry.Recover()
	span := sentry.StartSpan(c, "[Getposts]", sentry.TransactionName("GETPOSTS"), sentry.ContinueFromRequest(c.Request))
	defer span.Finish()

	response := getPostsCall(span.Context())

	c.JSON(200, response)
}
