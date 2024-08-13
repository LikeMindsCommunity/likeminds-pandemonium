package statuschk

// Response defines statuschkresposne GET API response body
type Response struct {
	Status int    `json:"status_code" binding:"required"`
	Data   string `json:"data"`
	Error  string `json:"error"`
}
