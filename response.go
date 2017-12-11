package daap

// ImagePullResponsePayload ...
type ImagePullResponsePayload struct {
	ID             string `json:"id"`
	Status         string `json:"status"`
	Progress       string `json:"progress"`
	ProgressDetail struct {
		Current uint32
		Total   uint32
	} `json:"progressDetail"`
}
