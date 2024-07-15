package response

type RatingSummaryResponse struct {
	OneStarCount   int     `json:"one_star_count"`
	TwoStarCount   int     `json:"two_star_count"`
	ThreeStarCount int     `json:"three_star_count"`
	FourStarCount  int     `json:"four_star_count"`
	FiveStarCount  int     `json:"five_star_count"`
	Average        float64 `json:"average"`
}