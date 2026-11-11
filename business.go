package kwai

import "context"

type BusinessService struct{ client *Client }

// LicenseUploadToken — 资质文件上传令牌.
func (s *BusinessService) LicenseUploadToken(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/bc/bcLicenseFileUploadToken", params)
}

// LicenseFileKey — 资质文件唯一标识.
func (s *BusinessService) LicenseFileKey(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/bc/bcLicenseFileUploadTokenToFileKey", params)
}

// CreateAccount — 新建广告主账户.
func (s *BusinessService) CreateAccount(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/bc/bcAdAccountCreate", params)
}

// Account — 查询广告主账户.
func (s *BusinessService) Account(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/bc/bcAccountQueryByAdAccount", params)
}

// Charge — 广告主账户充值.
func (s *FinanceService) Charge(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/finance/chargeFinance", params)
}

// Refund — 广告主账户退款.
func (s *FinanceService) Refund(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/finance/refundMoney", params)
}
