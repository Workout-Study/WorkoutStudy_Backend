package model

type GetFitGroupDetailApiResponse struct {
	PresentFitMateCount    int      `json:"presentFitMateCount"`
	MultiMediaEndPoints    []string `json:"multiMediaEndPoints"`
	FitGroupId             int      `json:"fitGroupId"`
	FitLeaderUserId        int      `json:"fitLeaderUserId"`
	FitGroupLeaderNickname string   `json:"fitGroupLeaderUserNickname"`
	FitGroupName           string   `json:"fitGroupName"`
	PenaltyAmount          int      `json:"penaltyAmount"`
	PenaltyAccountBankCode string   `json:"penaltyAccountBankCode"`
	PenaltyAccountNumber   string   `json:"penaltyAccountNumber"`
	Category               int      `json:"category"`
	Introduction           string   `json:"introduction"`
	Cycle                  int      `json:"cycle"`
	Frequency              int      `json:"frequency"`
	CreatedAt              string   `json:"createdAt"`
	MaxFitMate             int      `json:"maxFitMate"`
	State                  bool     `json:"state"`
}

type GetFitMatesApiResponse struct {
	FitGroupId      int    `json:"fitGroupId"`
	FitLeaderDetail Leader `json:"fitLeaderDetail"`
	FitMateDetails  []Mate `json:"fitMateDetails"`
}

type Leader struct {
	FitLeaderUserId       int    `json:"fitLeaderUserId"`
	FitLeaderUserNickname string `json:"fitLeaderUserNickname"`
	CreatedAt             string `json:"createdAt"`
}

type Mate struct {
	FitMateId           int    `json:"fitMateId"`
	FitMateUserId       int    `json:"fitMateUserId"`
	FitMateUserNickname string `json:"fitMateUserNickname"`
	CreatedAt           string `json:"createdAt"`
}

type GetUserInfoApiResponse struct {
	UserID   int    `json:"userId"`
	Nickname string `json:"nickname"`
}
