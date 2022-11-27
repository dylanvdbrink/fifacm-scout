package sofifa

import (
	"encoding/json"
	"fifacm-scout/internal/models"
	"fmt"
	"io"
	"net/http"
)

func GetLastDBUpdate() models.DBUpdate {
	resp, err := http.Get(fmt.Sprint(baseURL, "/api/player/history?id=228702"))
	if err != nil {
		panic("could not get sofifa dbupdate history")
	} else {
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)
		body, _ := io.ReadAll(resp.Body)

		var result HistoryRequest
		if jsonErr := json.Unmarshal(body, &result); jsonErr != nil {
			panic("cannot unmarshal JSON: " + jsonErr.Error())
		}

		lastUpdate := result.Data[len(result.Data)-1]

		return models.DBUpdate{UpdateID: lastUpdate[5].(string), Name: lastUpdate[0].(string)}
	}
}
