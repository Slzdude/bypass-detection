package models

type OutputRequest struct {
	RawRequest  []byte `json:"request"`
	RawResponse []byte `json:"response"`
}
