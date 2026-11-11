package kwai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const MaxFragmentSize = 4 << 20

// Uploader uses a single-file upload token, never the Marketing API access token.
type Uploader struct {
	endpoint string
	token    string
	http     *http.Client
}
type UploadResponse struct {
	Status   int             `json:"status"`
	Result   int             `json:"result"`
	Checksum string          `json:"checksum"`
	Size     int64           `json:"size"`
	Raw      json.RawMessage `json:"-"`
}

// NewUploader accepts a host from httpEndpointList or an explicit http(s) URL.
// Bare hosts use HTTPS; pass an explicit HTTP URL only if your upload node needs it.
func NewUploader(endpoint, token string, client *http.Client) (*Uploader, error) {
	if !strings.Contains(endpoint, "://") {
		endpoint = "https://" + endpoint
	}
	if err := validateURL(endpoint); err != nil {
		return nil, err
	}
	u, _ := url.Parse(endpoint)
	if u.Path != "" && u.Path != "/" {
		return nil, errors.New("kwai: upload endpoint must be an origin")
	}
	if token == "" {
		return nil, errors.New("kwai: upload token is required")
	}
	return &Uploader{strings.TrimRight(endpoint, "/"), token, newHTTPClient(client)}, nil
}

func (u *Uploader) request(ctx context.Context, method, path string, query url.Values, body []byte) (*UploadResponse, error) {
	query.Set("upload_token", u.token)
	req, err := http.NewRequestWithContext(ctx, method, u.endpoint+path+"?"+query.Encode(), bytes.NewReader(body))
	if err != nil {
		return nil, errors.New("kwai: invalid upload request")
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	raw, status, err := execute(u.http, req)
	if err != nil {
		return nil, err
	}
	var result UploadResponse
	if status != 200 || decodeJSON(raw, &result) != nil || (result.Result != 1 && !(path == "/api/upload" && result.Status == 1)) {
		return nil, errors.New("kwai: upload service rejected request")
	}
	result.Raw = raw
	return &result, nil
}

// License uploads the raw JPG/JPEG/PNG bytes in a single request. Obtain the
// token via Business.LicenseUploadToken, then resolve Business.LicenseFileKey.
func (u *Uploader) License(ctx context.Context, content []byte) (*UploadResponse, error) {
	if len(content) == 0 {
		return nil, errors.New("kwai: empty license file")
	}
	return u.request(ctx, http.MethodPost, "/api/upload", url.Values{}, content)
}

func (u *Uploader) Fragment(ctx context.Context, index int, fragment []byte) (*UploadResponse, error) {
	if index < 0 || len(fragment) == 0 || len(fragment) > MaxFragmentSize {
		return nil, errors.New("kwai: fragment index must be nonnegative and size within 1 byte to 4 MiB")
	}
	return u.request(ctx, http.MethodPost, "/api/upload/fragment", url.Values{"fragment_id": {strconv.Itoa(index)}}, fragment)
}
func (u *Uploader) Complete(ctx context.Context, count int) (*UploadResponse, error) {
	if count <= 0 {
		return nil, errors.New("kwai: fragment count must be positive")
	}
	return u.request(ctx, http.MethodPost, "/api/upload/complete", url.Values{"fragment_count": {strconv.Itoa(count)}}, nil)
}
func (u *Uploader) Resume(ctx context.Context) (*UploadResponse, error) {
	return u.request(ctx, http.MethodGet, "/api/upload/resume", url.Values{}, nil)
}

// Upload streams bounded chunks then completes the file. Transcoding is a
// separate Materials.TransferVideo call. It does not close the reader.
func (u *Uploader) Upload(ctx context.Context, reader io.Reader, chunkSize int) (*UploadResponse, error) {
	if reader == nil || chunkSize < 1 || chunkSize > MaxFragmentSize {
		return nil, errors.New("kwai: reader and chunk size within 1 to 4 MiB are required")
	}
	buffer := make([]byte, chunkSize)
	count := 0
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		n, err := io.ReadFull(reader, buffer)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return nil, err
		}
		if n > 0 {
			if _, uploadErr := u.Fragment(ctx, count, buffer[:n]); uploadErr != nil {
				return nil, uploadErr
			}
			count++
		}
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
	}
	if count == 0 {
		return nil, errors.New("kwai: cannot upload an empty file")
	}
	return u.Complete(ctx, count)
}
