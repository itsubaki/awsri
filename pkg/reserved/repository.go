package reserved

type Repository struct {
	Profile  string     `json:"profile"`
	Region   []string   `json:"region"`
	Internal RecordList `json:"internal"`
}

func (r *Repository) SelectAll() RecordList {
	return r.Internal
}
