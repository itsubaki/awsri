package reserved

type Repository struct {
	Profile  string     `json:"profile"`
	Internal RecordList `json:"internal"`
}

func (r *Repository) SelectAll() RecordList {
	return r.Internal
}
