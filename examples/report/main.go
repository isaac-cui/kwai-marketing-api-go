package main

import (
	"context"
	"fmt"
	"os"
	"time"

	kwai "github.com/isaac-cui/kwai-marketing-api-go"
)

func main() {
	client, err := kwai.New(kwai.Config{AccessToken: os.Getenv("KWAI_ACCESS_TOKEN")})
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// Replace accountId and the millisecond date range before running.
	result, err := client.Reports.AccountEffect(ctx, kwai.Params{
		"accountId": int64(1), "dataBeginTime": int64(1632873600000), "dataEndTime": int64(1633046399999),
		"timeZoneIana": "UTC+0", "granularity": 3, "pageNo": 1, "pageSize": 10, "withOutTax": false,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(string(result.Data))
}
