package  models
type BalanceServiceResponse struct {
AuthId    string  `json:"auth_id"`
Balance   float64 `json:"balance"`
}

type BalanceServiceRequest struct {
Balance   float64 `json:"balance"`
}
