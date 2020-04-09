package acemon

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type ErrNanInResponse struct {
	message string
}

type ErrInRequest struct {
	message string
}

func NewErrNanInResponse(message string) *ErrNanInResponse {
	return &ErrNanInResponse{
		message: message,
	}
}

func (e *ErrNanInResponse) Error() string {
	return e.message
}

func NewErrInRequest(message string) *ErrInRequest {
	return &ErrInRequest{message: message}
}

func (e *ErrInRequest) Error() string {
	return e.message
}

type IndexStat struct {
	Status     string    `json:"status"`
	LastUpdate time.Time `json:"last_update"`
}

type IndexSchedule struct {
	sync.RWMutex
	SourceURL string
	Table     map[string]IndexStat
}

func NewIndexSchedule(sourceURL string) *IndexSchedule {
	return &IndexSchedule{
		SourceURL: sourceURL,
	}
}

// Request Makes a requests and parses the body of the response
func Request(url string) ([]byte, error) {
	response, err := http.Get(fmt.Sprintf("%s", url))

	if err != nil {
		return []byte{}, err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}
	err = response.Body.Close()
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func (idx *IndexSchedule) UpdateSchedule() error {
	var errReq *ErrInRequest
	var newRate map[string]IndexStat
	response, err := Request(idx.SourceURL)
	if err == nil {
		err = json.Unmarshal(response, &newRate)
		if err == nil {
			idx.Lock()
			// Merging maps for having Fiat + Crypto conversion rates
			if idx.Table != nil {
				idx.Table = mergeMaps(idx.Table, newRate)
			} else {
				idx.Table = newRate
			}
			idx.Unlock()
		} else {
			return NewErrNanInResponse(fmt.Sprintf("Couldn't Unmarshal message in %s err: %s", idx.SourceURL, err))
		}
	} else {
		errReq = NewErrInRequest(fmt.Sprintf("Request failed to %s, err: %s", idx.SourceURL, err))
	}

	if errReq != nil {
		return errReq
	}
	return nil
}

// mergeMaps merges "b" map into "a" map
// Returns a map
func mergeMaps(a, b map[string]IndexStat) map[string]IndexStat {
	for k, v := range b {
		a[k] = v
	}
	return a
}
