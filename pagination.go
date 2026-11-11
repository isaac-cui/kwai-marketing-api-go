package kwai

import (
	"context"
	"encoding/json"
	"errors"
)

type Page struct {
	Total int64             `json:"total"`
	Data  []json.RawMessage `json:"data"`
}

// EachPage visits the standard data.total/data.data envelope. Returning false
// stops fetching. maxPages bounds requests even if the API reports a bad total.
// params is copied and never mutated; pageNo defaults to 1, pageSize to 20.
func (c *Client) EachPage(ctx context.Context, path string, params Params, maxPages int, visit func(*Page) bool) error {
	if maxPages <= 0 || visit == nil {
		return errors.New("kwai: positive maxPages and visitor are required")
	}
	query := Params{}
	for k, v := range params {
		query[k] = v
	}
	pageNo := 1
	pageSize := 20
	if v, ok := query["pageNo"]; ok {
		n, ok := v.(int)
		if !ok || n < 1 {
			return errors.New("kwai: pageNo must be a positive int")
		}
		pageNo = n
	}
	if v, ok := query["pageSize"]; ok {
		n, ok := v.(int)
		if !ok || n < 1 {
			return errors.New("kwai: pageSize must be a positive int")
		}
		pageSize = n
	}
	query["pageSize"] = pageSize
	for i := 0; i < maxPages; i++ {
		query["pageNo"] = pageNo
		response, err := c.Call(ctx, path, query)
		if err != nil {
			return err
		}
		var page Page
		if err = response.DecodeData(&page); err != nil {
			return err
		}
		if page.Data == nil || page.Total < 0 {
			return errors.New("kwai: invalid pagination envelope")
		}
		if !visit(&page) || len(page.Data) == 0 || int64(pageNo)*int64(pageSize) >= page.Total {
			return nil
		}
		pageNo++
	}
	return errors.New("kwai: pagination reached maxPages before completion")
}
