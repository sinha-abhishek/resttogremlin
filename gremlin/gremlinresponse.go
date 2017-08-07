package gremlin

type GremlinResponse struct {
	RequestID string                 `json:"requestId"`
	Status    map[string]interface{} `json:"status"`
	Result    map[string]interface{} `json:"result"`
}

func (gresponse *GremlinResponse) GetResultData() interface{} {
	result := gresponse.Result
	return result["data"]
}

func (gr *GremlinResponse) getStatusCode() float64 {
	status := gr.Status
	code := status["code"]
	return code.(float64)
}

func (gr *GremlinResponse) getRequestId() string {
	return gr.RequestID
}
