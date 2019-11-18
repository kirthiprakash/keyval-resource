package models

type EmptyVersion struct {
}

type Version map[string]string

type InRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type InResponse struct {
	Version Version `json:"version"`
}

type OutRequest struct {
	Source Source `json:"source"`
	Params map[string]string
}

type OutResponse struct {
	Version Version `json:"version"`
}

type CheckRequest struct {
	Source  Source       `json:"source"`
	Version EmptyVersion `json:"version"`
}

type CheckResponse []EmptyVersion

type Source struct{}
