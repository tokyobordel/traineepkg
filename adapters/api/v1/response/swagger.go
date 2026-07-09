package response

// Envelope описывает стандартный успешный ответ API.
type Envelope struct {
	Data       interface{} `json:"data"`
	Success    bool        `json:"success" example:"true"`
	ErrMessage string      `json:"err_message" example:""`
}

// ErrorEnvelope описывает стандартный ответ API с ошибкой.
type ErrorEnvelope struct {
	Data       interface{} `json:"data"`
	Success    bool        `json:"success" example:"false"`
	ErrMessage string      `json:"err_message" example:"invalid request body"`
}

// MessageData описывает тело ответа для операций без доменных данных.
type MessageData struct {
	Message string `json:"message" example:"logout successful"`
}
