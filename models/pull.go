package models

import "encoding/json"

type Result struct {
	Id              int              `json:"id" bson:"_id"`
	Title           string           `json:"title" bson:"title"`
	Slug            string           `json:"slug" bson:"slug"`
	PublicationDate int64            `json:"publication_date" bson:"publication_date"`
	Place           *Place           `json:"place" bson:"place"`
	Description     string           `json:"description" bson:"description"`
	Dates           []Date           `json:"dates" bson:"dates"`
	BodyText        string           `json:"body_text" bson:"body_text"`
	Location        Location         `json:"location" bson:"location"`
	Categories      []string         `json:"categories" bson:"categories"`
	TagLine         string           `json:"tagline" bson:"tagline"`
	AgeRestriction  *json.RawMessage `json:"age_restriction" bson:"age_restriction"`
	Price           string           `json:"price" bson:"price"`
	IsFree          bool             `json:"is_free" bson:"is_free"`
	Images          []Image          `json:"images" bson:"images"`
	FavouritesCount int              `json:"favourites_count" bson:"favourites_count"`
	CommentsCount   int              `json:"comments_count" bson:"comments_count"`
	SiteUrl         string           `json:"site_url" bson:"site_url"`
	ShortTitle      string           `json:"short_title" bson:"short_title"`
	Tags            []string         `json:"tags" bson:"tags"`
	Participants    []Participant    `json:"participants" bson:"participants"`
}

type Participant struct {
	Role  Role  `json:"role" bson:"role"`
	Agent Agent `json:"agent" bson:"agent"`
}

type Role struct {
	Id         int    `json:"id" bson:"_id"`
	Slug       string `json:"slug" bson:"slug"`
	Name       string `json:"name" bson:"name"`
	NamePlural string `json:"name_plural" bson:"name_plural"`
}

type Agent struct {
	Id              int     `json:"id" bson:"_id"`
	Title           string  `json:"title" bson:"title"`
	Slug            string  `json:"slug" bson:"slug"`
	Description     string  `json:"description" bson:"description"`
	BodyText        string  `json:"body_text" bson:"body_text"`
	Rank            float64 `json:"rank" bson:"rank"`
	AgentType       string  `json:"agent_type" bson:"agent_type"`
	Images          []Image `json:"images" bson:"images"`
	FavoritesCount  int     `json:"favorites_count" bson:"favorites_count"`
	CommentsCount   int     `json:"comments_count" bson:"comments_count"`
	SiteUrl         string  `json:"site_url" bson:"site_url"`
	DisableComments bool    `json:"disable_comments" bson:"disable_comments"`
	IsStub          bool    `json:"is_stub" bson:"is_stub"`
}

type Image struct {
	Image  string `json:"image" bson:"image"`
	Source Source `json:"source" bson:"source"`
}

type Source struct {
	Link   string `json:"link" bson:"link"`
	Source string `json:"source" bson:"source"`
}

type Place struct {
	Id          int         `json:"id" bson:"_id"`
	Title       string      `json:"title" bson:"title"`
	Slug        string      `json:"slug" bson:"slug"`
	Address     string      `json:"address" bson:"address"`
	Phone       string      `json:"phone" bson:"phone"`
	IsStub      bool        `json:"is_stub" bson:"is_stub"`
	SiteUrl     string      `json:"site_url" bson:"site_url"`
	Coordinates Coordinates `json:"coords" bson:"coords"`
	Subway      string      `json:"subway" bson:"subway"`
	IsClosed    bool        `json:"is_closed" bson:"is_closed"`
	Location    string      `json:"location" bson:"location"`
}

type Location struct {
	Slug        string      `json:"slug" bson:"slug"`
	Name        string      `json:"name" bson:"name"`
	Timezone    string      `json:"timezone" bson:"timezone"`
	Coordinates Coordinates `json:"coords" bson:"coords"`
	Language    string      `json:"language" bson:"language"`
	Currency    string      `json:"currency" bson:"currency"`
}

type Coordinates struct {
	Latitude  float64 `json:"lat" bson:"lat"`
	Longitude float64 `json:"lon" bson:"lon"`
}

type Date struct {
	StartDate    *string `json:"start_date" bson:"start_date"`
	StartTime    *string `json:"start_time" bson:"start_time"`
	Start        int64   `json:"start" bson:"start"`
	EndDate      *string `json:"end_date" bson:"end_date"`
	EndTime      *string `json:"end_time" bson:"end_time"`
	End          int64   `json:"end" bson:"end"`
	IsContinuous bool    `json:"is_continuous" bson:"is_continuous"`
	IsEndless    bool    `json:"is_endless" bson:"is_endless"`
	IsStartless  bool    `json:"is_startless" bson:"is_startless"`
	//Schedules []interface{}
	UsePlaceSchedule bool `json:"use_place_schedule" bson:"use_place_schedule"`
}
