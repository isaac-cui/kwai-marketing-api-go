# Kwai Marketing API Go SDK

快手海外 Marketing API 的非官方 Go SDK。覆盖主手册及商务管理附属手册的 **61 个业务接口**，另支持 OAuth 刷新、分页、视频分片上传、资质文件上传。仅依赖 Go 标准库，Go 1.22+。

[完整接口索引](docs/endpoints.md) · [机器可读接口清单](docs/endpoints.json)

## 接入

仓库位于 https://github.com/isaac-cui/kwai-marketing-api-go。配置好工作账号的仓库访问权限后，其他项目可安装：

```sh
go get github.com/isaac-cui/kwai-marketing-api-go@main
```

私有仓库需设置 `GOPRIVATE=github.com/isaac-cui/*`，并配置 GitHub 认证。本地联调也可通过 go.mod 的 replace 指向源码目录。

```go
import (
    "context"
    "fmt"
    "os"
    kwai "github.com/isaac-cui/kwai-marketing-api-go"
)

client, err := kwai.New(kwai.Config{AccessToken: os.Getenv("KWAI_ACCESS_TOKEN")})
if err != nil { panic(err) }
response, err := client.Reports.AccountEffect(context.Background(), kwai.Params{
    "accountId": int64(1),
    "dataBeginTime": int64(1632873600000),
    "dataEndTime": int64(1633046399999),
    "granularity": 3,
    "timeZoneIana": "UTC+0",
    "pageNo": 1,
    "pageSize": 10,
    "withOutTax": false,
})
if err != nil { panic(err) }
fmt.Println(string(response.Data))
```

以上为调用片段；可执行示例位于 [examples/report/main.go](examples/report/main.go)。替换账户和日期后，在环境中设置 KWAI_ACCESS_TOKEN，执行 `go run ./examples/report`。

## 请求与响应

服务分组：Reports、Campaigns、Units、Creatives、Materials、Organic、Business、Finance。所有具名方法均接收请求结构体或 Params，并返回 *Response。ReportRequest、CampaignPageRequest 提供常见字段类型，其余接口按文档原名传 Params；嵌套列表使用 []kwai.Params，ID 与金额使用 int64。false、0、空数组不被 Params 丢弃。结构体的可选布尔/状态字段使用指针以区别未传与零值。

`client.Call(ctx, "/rest/n/mapi/...", params)` 可调用新增接口。HTTP 2xx 且业务 status=200 才算成功。Response 保留 Status、Message、RequestID、Data 和原始 Raw。使用 `response.DecodeData(&target)` 解码到自定义结构体或 map；map 中数字为 json.Number，避免大整数 ID 经 float64 丢失精度。

```go
var apiErr *kwai.APIError
if errors.As(err, &apiErr) {
    fmt.Println(apiErr.HTTPStatus, apiErr.Status, apiErr.RequestID)
}
```

此片段需要额外 import errors。APIError.Body 保留服务端错误响应；Message 可能包含文档所示嵌套 JSON 错误详情。

## OAuth

```go
token, err := client.RefreshAccessToken(ctx, clientID, clientSecret, refreshToken)
if err != nil { return err }
// 持久化 token.RefreshToken；按 ExpiresIn 提前刷新。
client = client.WithAccessToken(token.AccessToken)
```

手册规定 access_token 约 1 小时，refresh_token 默认 360 天。刷新显式进行，不在失败的写入请求后自动刷新重放。WithAccessToken 返回新的客户端，原客户端不变，可并发调用。配置 BaseURL、TokenURL、HTTPClient 可接入本地测试或定制超时；默认超时 30 秒。每次调用支持 context 取消/超时。客户端拒绝跟随重定向。

## 分页

```go
err := client.EachPage(ctx, "/rest/n/mapi/report/dspAccountEffectQuery",
    kwai.Params{"accountId": int64(1), "pageSize": 100}, 50,
    func(page *kwai.Page) bool {
        // page.Data 是 []json.RawMessage；消费本页后返回 true 继续。
        return true
    })
```

示例需补齐该报表的日期和时区。仅适用于 data.total/data.data 结构。pageNo/pageSize 为正 int，默认 1/20；到达总数、空页或回调返回 false 时停止。超过 maxPages 返回错误，防止无界请求。不会修改输入 Params。

## 视频上传

1. 调用 Materials.UploadToken，传 accountId、fileType（1 竖版/2 横版）、fileFormat="mp4" 等字段。
2. 将响应 data.token 和 data.httpEndpointList 中的节点传给 NewUploader。
3. 打开本地文件，以 `uploader.Upload(ctx, reader, kwai.MaxFragmentSize)` 按不超过 4 MiB 的分片上传并完成。调用方负责关闭 reader。
4. 调用 Materials.TransferVideo，传 accountId、token、fileType、name、footageLanguageTag、audioLanguageTag、aigc。该接口首次调用发起转码，后续调用查询。
5. 由调用方按 500ms～1s 间隔、有截止时间地查询，直到 step=4 且 checkFailedList 为空。校验失败时应停止处理。

同时提供 Fragment、Complete、Resume 以支持自行管理断点。上传 token 只属于单个文件，不能复用。裸域名默认使用 HTTPS；手册的示例使用 HTTP，若节点确需 HTTP 可传显式 URL。上传请求只使用 upload_token，不携带 Access-Token。

## 资质与商务

Business.LicenseUploadToken 获取资质 token；NewUploader 后用 `uploader.License(ctx, content)` 发送 JPG/JPEG/PNG 原始字节，再用 Business.LicenseFileKey 获取 fileBlobKey，用于开户的 blobstoreKey。Business.CreateAccount 与 Business.Account 支持开户与查询；Finance.Charge、Finance.Refund 支持资金请求。参数及金额单位见接口索引和附属手册。

## 验证与已知边界

```sh
go test ./...
go vet ./...
```

测试使用本地 HTTP 服务，覆盖字段保持、大整数、错误响应、OAuth、分页、取消、重定向隔离及分片上传，不需要真实凭据。不自动重试业务调用，也不自动调整时间、金额、分页上限或账户权限。

应用查询路径在手册中存在冲突，应用新建已暂停；PWA 等部分章节缺少明确方法标注。具体差异见接口索引。当前未使用真实账号联调，接口实际权限、文档不一致之处及转码结果需在目标账号确认。

