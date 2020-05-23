package awsdynamodb

// RekognitionItem :
type RekognitionItem struct {
	PrimaryKey string      `json:"pk"`           // ur => pk
	SortKey    string      `json:"sk,omitempty"` // sk => sort key
	URL        string      `json:"url"`
	Keywords   interface{} `json:"keywords"`
}
