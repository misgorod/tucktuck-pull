package pull

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"net/http"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
)

type pullResponse struct {
	Count    int `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []result `json:"results"`
}

type result struct {
	Id              int `json:"id" bson:"_id"`
	Title           string `json:"title" bson:"title"`
	Slug            string `json:"slug" bson:"slug"`
	PublicationDate int64 `json:"publication_date" bson:"publication_date"`
	Place           place `json:"place" bson:"place"`
	Description     string `json:"description" bson:"description"`
	Dates           []date `json:"dates" bson:"dates"`
	BodyText        string `json:"body_text" bson:"body_text"`
	Location        location `json:"location" bson:"location"`
	Categories      []string `json:"categories" bson:"categories"`
	TagLine         string `json:"tagline" bson:"tagline"`
	AgeRestriction  int    `json:"age_restriction" bson:"age_restriction"`
	Price           string `json:"price" bson:"price"`
	SsFree          bool `json:"is_free" bson:"is_free"`
	Images          []image `json:"images" bson:"images"`
	FavouritesCount int    `json:"favourites_count" bson:"favourites_count"`
	CommentsCount   int    `json:"comments_count" bson:"comments_count"`
	SiteUrl         string `json:"site_url" bson:"site_url"`
	ShortTitle      string `json:"short_title" bson:"short_title"`
	Tags            []string `json:"tags" bson:"tags"`
	Participants    []participant `json:"participants" bson:"participants"`
}

type participant struct {
	Role role `json:"role" bson:"role"`
	Agent agent `json:"agent" bson:"agent"`
}

type role struct {
	Id int `json:"id" bson:"_id"`
	Slug string `json:"slug" bson:"slug"`
	Name string `json:"name" bson:"name"`
	NamePlural string `json:"name_plural" bson:"name_plural"`
}

type agent struct {
	Id int `json:"id" bson:"_id"`
	Title string `json:"title" bson:"title"`
	Slug string `json:"slug" bson:"slug"`
	Description string `json:"description" bson:"description"`
	BodyText string `json:"body_text" bson:"body_text"`
	Rank int `json:"rank" bson:"rank"`
	AgentType string `json:"agent_type" bson:"agent_type"`
	Images []image `json:"images" bson:"images"`
	FavoritesCount int `json:"favorites_count" bson:"favorites_count"`
	CommentsCount int `json:"comments_count" bson:"comments_count"`
	SiteUrl string `json:"site_url" bson:"site_url"`
	DisableComments bool `json:"disable_comments" bson:"disable_comments"`
	IsStub bool `json:"is_stub" bson:"is_stub"`
}

type image struct {
	Image string `json:"image" bson:"image"`
	Source source `json:"source" bson:"source"`
}

type source struct {
	Link string `json:"link" bson:"link"`
	Source string `json:"source" bson:"source"`
}

type place struct {
	Id int `json:"id" bson:"_id"`
	Title string `json:"title" bson:"title"`
	Slug string `json:"slug" bson:"slug"`
	Address string `json:"address" bson:"address"`
	Phone string `json:"phone" bson:"phone"`
	IsStub bool `json:"is_stub" bson:"is_stub"`
	SiteUrl string `json:"site_url" bson:"site_url"`
	Coordinates coordinates `json:"coords" bson:"coords"`
	Subway string `json:"subway" bson:"subway"`
	IsClosed bool `json:"is_closed" bson:"is_closed"`
	Location string `json:"location" bson:"location"`
}

type location struct {
	Slug string `json:"slug" bson:"slug"`
	Name string `json:"name" bson:"name"`
	Timezone string `json:"timezone" bson:"timezone"`
	Coordinates coordinates `json:"coords" bson:"coords"`
	Language string `json:"language" bson:"language"`
	Currency string `json:"currency" bson:"currency"`
}

type coordinates struct {
	Latitude  float64 `json:"lat" bson:"lat"`
	Longitude float64 `json:"lon" bson:"lon"`
}

type date struct {
	StartDate string `json:"start_date" bson:"start_date"`
	StartTime string `json:"start_time" bson:"start_time"`
	Start int64 `json:"start" bson:"start"`
	EndDate string `json:"end_date" bson:"end_date"`
	EndTime string `json:"end_time" bson:"end_time"`
	End int64 `json:"end" bson:"end"`
	IsContinuous bool `json:"is_continuous" bson:"is_continuous"`
	IsEndless bool `json:"is_endless" bson:"is_endless"`
	IsStartless bool `json:"is_startless" bson:"is_startless"`
	//Schedules []interface{}
	UsePlaceSchedule bool `json:"use_place_schedule" bson:"use_place_schedule"`
}

type Handler struct {
	client *mongo.Client `json:"-"`
	insertResult InsertResult `json:"insertResult"`
}

type InsertResult struct {
	MatchedCount int64
	ModifiedCount int64
	UpsertedCount int64
	UpsertedID interface{}
}

func New(ctx context.Context, client *mongo.Client) *Handler {
	handler := &Handler{client:client}
	collection := client.Database("tucktuck").Collection("events")
	go func() {
		select {
		case <-ctx.Done():
			log.WithError(ctx.Err()).Error(" received cancel from external context")
			return
		default:
		}
		logger := log.WithFields(log.Fields{
			"handler": "pullHandler",
			"method":  "new",
		})
		actualSince := time.Now().Unix()
		url := fmt.Sprintf("https://kudago.com/public-api/v1.4/events/?fields=id,publication_date,dates,title,short_title,slug,place,description,body_text,location,categories,tagline,age_restriction,price,is_free,images,favorites_count,comments_count,site_url,tags,participants&expand=images,place,location,dates,participants&text_format=text&location=msk&actual_since=%v", actualSince)
		response, err := http.Get(url)
		if err != nil {
			logger.WithError(err).Error("couldn't get response from kudago")
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			logger.WithError(err).Error("couldn't read response from kudago")
		}
		var pullResponse pullResponse
		err = json.Unmarshal(body, &pullResponse)
		if err != nil {
			logger.WithError(err).Error("couldn't unmarshal response from kudago")
		}
		tmp := make([]interface{}, len(pullResponse.Results))
		for _, result := range pullResponse.Results {
			tmp = append(tmp, result)
		}
		updateResult, err := collection.UpdateMany(ctx, bson.D{}, tmp, options.Update().SetUpsert(true))
		if err != nil {
			logger.WithError(err).Error()
		}
		handler.insertResult := InsertResult{
			MatchedCount:  updateResult.MatchedCount,
			ModifiedCount: updateResult.ModifiedCount,
			UpsertedCount: updateResult.UpsertedCount,
			UpsertedID:    updateResult.UpsertedID,
		}
		logger.WithFields(log.Fields{
			"matched count": updateResult.MatchedCount,
			"modified count": updateResult.ModifiedCount,
			"upserted count": updateResult.UpsertedCount,
			"upserted id": updateResult.UpsertedID,
		})
		time.Sleep(time.Hour)
	}()
	return handler
}

func (p *Handler) Get(w http.ResponseWriter, r *http.Request) {

}
