package dtos

type AirtableListRecords struct {
	Offset  string `json:"offset"`
	Results []map[string]interface{}
}
