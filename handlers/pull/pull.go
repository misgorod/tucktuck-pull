package pull

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/misgorod/tucktuck-pull/common"
	"github.com/misgorod/tucktuck-pull/models"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type pullResponse struct {
	Count    int             `json:"count"`
	Next     string          `json:"next"`
	Previous *string         `json:"previous"`
	Results  []models.Result `json:"results"`
}

type Result struct {
	Count      int           `json:"count"`
	UpsertedId []string      `json:"upserted_id"`
	Duration   time.Duration `json:"duration"`
	Error      error         `json:"-"`
}

type Handler struct {
	client     *models.Client
	lastResult Result
	mutex      *sync.Mutex
}

func (h *Handler) makeRequest(ctx context.Context, logger *log.Logger) Result {
	startTime := time.Now()
	actualSince := time.Now().Unix()
	url := fmt.Sprintf("https://kudago.com/public-api/v1.4/events/?fields=id,publication_date,dates,title,short_title,slug,place,description,body_text,location,categories,tagline,age_restriction,price,is_free,images,favorites_count,comments_count,site_url,tags,participants&expand=images,place,location,dates,participants&text_format=text&location=msk&actual_since=%v", actualSince)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	lastResult := Result{
		Count:      0,
		UpsertedId: make([]string, 0),
		Error:      nil,
	}
	for url != "" {
		response, err := client.Get(url)
		if err != nil {
			lastResult.Error = err
			return lastResult
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			lastResult.Error = err
			return lastResult
		}
		var pullResponse pullResponse
		err = json.Unmarshal(body, &pullResponse)
		if err != nil {
			lastResult.Error = err
			return lastResult
		}
		upsertResult, err := h.client.UpsertMany(ctx, pullResponse.Results)
		if err != nil {
			lastResult.Error = err
			return lastResult
		}
		lastResult.UpsertedId = append(lastResult.UpsertedId, upsertResult.UpsertedID...)
		logger.WithField("upsert Result", fmt.Sprintf("%+v", upsertResult)).Info("got events")
		url = pullResponse.Next
	}
	lastResult.Error = nil
	lastResult.Count = len(lastResult.UpsertedId)
	lastResult.Duration = time.Since(startTime)
	return lastResult
}

func New(client *models.Client) *Handler {
	handler := &Handler{
		client: client,
		lastResult: Result{
			Count:      0,
			UpsertedId: make([]string, 0),
			Error:      nil,
		},
		mutex: &sync.Mutex{},
	}
	return handler
}

func (h *Handler) Start(period time.Duration, logger *log.Logger) {
	ticker := time.NewTicker(period)
	defer ticker.Stop()
	go func() {
		for ; ; <-ticker.C {
			result := h.makeRequest(context.Background(), logger)
			if result.Error != nil {
				logger.WithError(result.Error).Error()
			} else {
				logger.WithField("result", fmt.Sprintf("%+v", result)).Info("got result")
			}
			h.mutex.Lock()
			h.lastResult = result
			h.mutex.Unlock()
		}
	}()
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	if h.lastResult.Error != nil {
		common.RespondError(r.Context(), w, http.StatusInternalServerError, h.lastResult.Error)
		return
	}
	common.RespondJSON(r.Context(), w, http.StatusOK, h.lastResult)
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	result := h.makeRequest(r.Context(), log.StandardLogger())
	if result.Error != nil {
		common.RespondError(r.Context(), w, http.StatusInternalServerError, result.Error)
		return
	}
	h.mutex.Lock()
	h.lastResult = result
	h.mutex.Unlock()
	common.RespondJSON(r.Context(), w, http.StatusOK, result)
}
