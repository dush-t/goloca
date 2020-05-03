package api

type responseMessage map[string]interface{}

func failMessage(errorMessage string) responseMessage {
	return responseMessage{
		"status": "failed",
		"error":  errorMessage,
		"id":     nil,
	}
}

func successMessage(id string) responseMessage {
	return responseMessage{
		"status": "success",
		"error":  nil,
		"id":     id,
	}
}
