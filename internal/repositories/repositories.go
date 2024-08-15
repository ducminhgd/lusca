package repositories

type QueryOptions struct {
	Limit  int `json:"limit" default:"100"`
	Offset int `json:"offset" default:"0"`
}
