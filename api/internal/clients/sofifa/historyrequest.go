package sofifa

type HistoryRequest struct {
	Data  [][]interface{} `json:"data"`
	Start string          `json:"start"`
}
