package kwai

import "context"

type ReportsService struct{ client *Client }

type FinanceService struct{ client *Client }

type CampaignsService struct{ client *Client }

type UnitsService struct{ client *Client }

type CreativesService struct{ client *Client }

type MaterialsService struct{ client *Client }

type OrganicService struct{ client *Client }

func (c *Client) initServices() {
	c.Business = &BusinessService{client: c}
	c.Reports = &ReportsService{client: c}
	c.Finance = &FinanceService{client: c}
	c.Campaigns = &CampaignsService{client: c}
	c.Units = &UnitsService{client: c}
	c.Creatives = &CreativesService{client: c}
	c.Materials = &MaterialsService{client: c}
	c.Organic = &OrganicService{client: c}
}

// AccountEffect — 广告主消耗.
func (s *ReportsService) AccountEffect(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/report/dspAccountEffectQuery", params)
}

// AgentEffect — 代理商消耗.
func (s *ReportsService) AgentEffect(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/report/crmAccountQuery", params)
}

// AgentBalance — 代理商余额.
func (s *ReportsService) AgentBalance(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/report/financeBalanceQuery", params)
}

// AgentMonthlyConsumption — 代理商月度消耗.
func (s *ReportsService) AgentMonthlyConsumption(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/report/agentMonthlyConsumptionReportQuery", params)
}

// AccountMonthlyConsumption — 广告主月度消耗.
func (s *ReportsService) AccountMonthlyConsumption(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/report/accountMonthlyConsumptionReportQuery", params)
}

// AccountMonthlyDetailConsumption — 广告主消耗明细.
func (s *ReportsService) AccountMonthlyDetailConsumption(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/report/accountMonthlyDetailConsumptionReportQuery", params)
}

// CampaignEffect — 广告计划报表.
func (s *ReportsService) CampaignEffect(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/report/dspCampaignEffectQuery", params)
}

// UnitEffect — 广告组报表.
func (s *ReportsService) UnitEffect(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/report/dspUnitEffectQuery", params)
}

// CreativeEffect — 广告创意报表.
func (s *ReportsService) CreativeEffect(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/report/dspCreativeEffectQuery", params)
}

// VideoEffect — 视频报表.
func (s *ReportsService) VideoEffect(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/report/dspPhotoEffectQuery", params)
}

// AppTargetEffect — App 分国家报表.
func (s *ReportsService) AppTargetEffect(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/report/dspAppTargetEffectQuery", params)
}

// PopulationEffect — 投放地/操作系统报表.
func (s *ReportsService) PopulationEffect(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/report/dspPopulationAnalysisEffectQuery", params)
}

// AccountInfos — 账户详情.
func (s *ReportsService) AccountInfos(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/report/dspAccountInfosQuery", params)
}

// Accounts — 账户列表.
func (s *ReportsService) Accounts(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/report/crmAccountQueryByAgentOrCorp", params)
}

// AdxDaily — 三方广告日报.
func (s *ReportsService) AdxDaily(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/adx/getAdxReport/day", params)
}

// AccountBalance — 广告主余额.
func (s *FinanceService) AccountBalance(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/finance/financeAccountBalanceQuery", params)
}

// List — Campaign List.
func (s *CampaignsService) List(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/campaign/dspCampaignPageQueryPerformance", params)
}

// Create — Campaign Create.
func (s *CampaignsService) Create(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/campaign/dspCampaignAddPerformance", params)
}

// Update — Campaign Update.
func (s *CampaignsService) Update(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/campaign/dspCampaignUpdatePerformance", params)
}

// UpdateStatus — Campaign UpdateStatus.
func (s *CampaignsService) UpdateStatus(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/campaign/dspCampaignUpdateOpenStatusPerformance", params)
}

// Delete — Campaign Delete.
func (s *CampaignsService) Delete(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/campaign/dspCampaignDeletePerformance", params)
}

// List — Unit List.
func (s *UnitsService) List(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/unit/dspUnitPageQueryPerformance", params)
}

// Create — Unit Create.
func (s *UnitsService) Create(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/unit/dspUnitAddPerformance", params)
}

// Update — Unit Update.
func (s *UnitsService) Update(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/unit/dspUnitUpdatePerformance", params)
}

// UpdateStatus — Unit UpdateStatus.
func (s *UnitsService) UpdateStatus(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/unit/dspUnitUpdateOpenStatusPerformance", params)
}

// Delete — Unit Delete.
func (s *UnitsService) Delete(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/unit/dspUnitDeletePerformance", params)
}

// List — Creative List.
func (s *CreativesService) List(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/creative/dspCreativePageQueryPerformance", params)
}

// Create — Creative Create.
func (s *CreativesService) Create(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/creative/dspCreativeAddPerformance", params)
}

// Update — Creative Update.
func (s *CreativesService) Update(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/creative/dspCreativeUpdatePerformance", params)
}

// UpdateStatus — Creative UpdateStatus.
func (s *CreativesService) UpdateStatus(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/creative/dspCreativeUpdateOpenStatusPerformance", params)
}

// Delete — Creative Delete.
func (s *CreativesService) Delete(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/creative/dspCreativeDeletePerformance", params)
}

// AuditResult — 广告组审核建议.
func (s *UnitsService) AuditResult(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/unit/dspUnitAuditResultPerformance", params)
}

// UpdateBudget — 批量修改预算.
func (s *UnitsService) UpdateBudget(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/unit/dspUnitUpdateBudgetPerformance", params)
}

// UpdateBid — 批量修改出价.
func (s *UnitsService) UpdateBid(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/unit/dspUnitUpdateBidPerformance", params)
}

// UpdateCboBid — CBO 出价.
func (s *UnitsService) UpdateCboBid(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/unit/dspUnitUpdateCboBidPerformance", params)
}

// DerivedActionTypes — 衍生行为.
func (s *UnitsService) DerivedActionTypes(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/unit/derivedActionType", params)
}

// AuditResult — 创意审核建议.
func (s *CreativesService) AuditResult(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/creative/dspCreativeAuditResultPerformance", params)
}

// PostAuditResult — 社区审核建议.
func (s *CreativesService) PostAuditResult(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/post/profilePostAuditResultPerformance", params)
}

// Apps — 应用查询（接口清单路径，详见文档差异说明）.
func (s *MaterialsService) Apps(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/material/dspMaterialAppQuery", params)
}

// Pixels — Pixel 查询.
func (s *MaterialsService) Pixels(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/material/dspMaterialPixelQuery", params)
}

// Videos — 视频查询.
func (s *MaterialsService) Videos(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/material/dspMaterialVideoQuery", params)
}

// UploadToken — 获取单文件上传令牌.
func (s *MaterialsService) UploadToken(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/material/dspMaterialFileUploadToken", params)
}

// TransferVideo — 发起/查询视频转码.
func (s *MaterialsService) TransferVideo(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/material/dspMaterialVideoTransfer", params)
}

// CreatePwa — 创建 PWA.
func (s *MaterialsService) CreatePwa(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/material/createPWA", params)
}

// Pwas — 查询 PWA.
func (s *MaterialsService) Pwas(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/material/pageQueryPWA", params)
}

// RefreshPwa — 重试 PWA.
func (s *MaterialsService) RefreshPwa(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/material/refreshPWA", params)
}

// AppealPwa — 申诉 PWA.
func (s *MaterialsService) AppealPwa(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/material/appealPWA", params)
}

// Avatars — 头像.
func (s *MaterialsService) Avatars(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/material/dspMaterialAvatarQuery", params)
}

// Populations — 人群包.
func (s *MaterialsService) Populations(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/material/dspMaterialPopulationQuery", params)
}

// Playables — 试玩素材.
func (s *MaterialsService) Playables(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/material/dspMaterialPlayableQuery", params)
}

// Interactives — 互动样式.
func (s *MaterialsService) Interactives(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/material/dspMaterialInteractiveQuery", params)
}

// Accounts — 已授权原生账户.
func (s *OrganicService) Accounts(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/organic/authKwaiAccountQuery", params)
}

// Posts — 原生作品（已弃用）.
func (s *OrganicService) Posts(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/organic/profilePostQuery", params)
}

// PostsV2 — 原生作品 V2.
func (s *OrganicService) PostsV2(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/organic/profilePostQueryV2", params)
}

// PushVideo — 视频同步至原生账户.
func (s *OrganicService) PushVideo(ctx context.Context, params any) (*Response, error) {
	return s.client.Call(ctx, "/rest/n/mapi/organic/profileMaterialVideoPush", params)
}
