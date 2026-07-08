package response

// Envelope стандартный успешный ответ API.
type Envelope struct {
	Data       interface{} `json:"data"`
	Success    bool        `json:"success" example:"true"`
	ErrMessage string      `json:"err_message" example:""`
	SpreadID   string      `json:"spread_id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// ErrorEnvelope стандартный ответ API с ошибкой.
type ErrorEnvelope struct {
	Data       interface{} `json:"data"`
	Success    bool        `json:"success" example:"false"`
	ErrMessage string      `json:"err_message" example:"invalid request body"`
	SpreadID   string      `json:"spread_id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// MessageData тело ответа для операций без доменных данных.
type MessageData struct {
	Message string `json:"message" example:"logout successful"`
}
