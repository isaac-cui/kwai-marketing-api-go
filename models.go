package kwai

// ReportRequest covers common account-effect report fields. Use Params for
// endpoint-specific filters and future additions. All timestamps are milliseconds.
type ReportRequest struct {
	AccountID     int64  `json:"accountId"`
	DataBeginTime int64  `json:"dataBeginTime"`
	DataEndTime   int64  `json:"dataEndTime"`
	Granularity   int    `json:"granularity,omitempty"`
	TimeZoneIana  string `json:"timeZoneIana"`
	PageNo        int    `json:"pageNo,omitempty"`
	PageSize      int    `json:"pageSize,omitempty"`
	WithOutTax    *bool  `json:"withOutTax,omitempty"`
}

type CampaignPageRequest struct {
	AccountID      int64   `json:"accountId"`
	AdCategory     int     `json:"adCategory"`
	CampaignIDList []int64 `json:"campaignIdList,omitempty"`
	Status         *int    `json:"status,omitempty"`
	PageNo         int     `json:"pageNo,omitempty"`
	PageSize       int     `json:"pageSize,omitempty"`
}
