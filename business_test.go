package kwai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBusinessFinanceContract(t *testing.T) {
	calls := 0
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls == 1 {
			if r.URL.Path != "/rest/n/mapi/bc/bcLicenseFileUploadToken" {
				t.Error("wrong license token path")
			}
			var params Params
			_ = json.NewDecoder(r.Body).Decode(&params)
			if params["fileFormat"] != "png" {
				t.Error("wrong license format")
			}
			fmt.Fprint(w, `{"status":200,"data":{"token":"single-file"}}`)
		} else {
			if r.URL.Path != "/rest/n/mapi/finance/chargeFinance" {
				t.Error("wrong finance path")
			}
			var params struct {
				ChargeMoney struct {
					Units  int64 `json:"units"`
					Micros int64 `json:"micros"`
				} `json:"chargeMoney"`
			}
			_ = json.NewDecoder(r.Body).Decode(&params)
			if params.ChargeMoney.Units != 10 || params.ChargeMoney.Micros != 500000 {
				t.Error("money values changed")
			}
			fmt.Fprint(w, `{"status":500,"message":"declined","requestId":"finance-trace"}`)
		}
	})
	ctx := context.Background()
	if _, err := client.WithAccessToken("updated").Business.LicenseUploadToken(ctx, Params{"fileFormat": "png"}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Finance.Charge(ctx, Params{"accountId": int64(1), "chargeType": 1, "chargeMoney": Params{"units": int64(10), "micros": int64(500000)}}); err == nil {
		t.Fatal("ignored rejected charge")
	}
	if calls != 2 {
		t.Fatal("finance request retried")
	}
}

func TestLicenseUploadCompatibility(t *testing.T) {
	for _, body := range []string{`{"result":1}`, `{"status":1}`} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/upload" || r.Method != "POST" || r.URL.Query().Get("upload_token") != "license" {
				t.Error("wrong license protocol")
			}
			fmt.Fprint(w, body)
		}))
		uploader, _ := NewUploader(server.URL, "license", nil)
		if _, err := uploader.License(context.Background(), []byte("png-content")); err != nil {
			t.Fatal(err)
		}
		server.Close()
	}
}
