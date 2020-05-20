package models

// QueryFilter :
type QueryFilter struct {
	Offset int64       `json:"offset"`
	Limit  int64       `json:"limit"`
	Max    int64       `json:"max"`
	Order  string      `json:"order"`
	Sort   string      `json:"sort"`
	Search string      `json:"search"`
	Extras interface{} `json:"extras"`
}

// QueryFilters :
type QueryFilters []QueryFilter

// QueryFilterAlias :
type QueryFilterAlias QueryFilter
