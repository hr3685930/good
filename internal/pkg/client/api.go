package client

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/spf13/cast"
	"good/internal/pkg/errs"
	"good/pkg/tracing/http"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/url"
	"time"
)

// APIClient APIClient
type APIClient struct {
}

// NewAPIClient NewAPIClient
func NewAPIClient() *APIClient {
	return &APIClient{}
}

// GetRPCConn GetRPCConn
func (c APIClient) GetRPCConn(ctx context.Context, endpoint string) (conn *grpc.ClientConn, err error) {
	return grpc.DialContext(ctx,
		endpoint,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			MinConnectTimeout: time.Second * 5,
		}),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpc_opentracing.UnaryClientInterceptor(),
		)),
	)
}

// Get Get
func (c APIClient) Get(ctx context.Context, endpoint string, path string, options map[string]interface{}) (res []byte, err error) {
	httpCli := http.NewClient()
	baseURL, err := url.Parse("http://" + endpoint)
	if err != nil {
		return nil, err
	}

	baseURL.Path += path
	params := url.Values{}
	for key, value := range options {
		params.Add(key, cast.ToString(value))
	}
	baseURL.RawQuery = params.Encode()
	httpCli.Timeout(5 * time.Second).Get(baseURL.String())
	resp, body, errArr := httpCli.TraceEndByte(ctx)
	if len(errArr) > 0 {
		return nil, errs.InternalError(errArr[0].Error())
	}

	if resp.StatusCode >= 300 {
		return nil, errs.InternalError(string(body))
	}

	return body, nil
}

// Post Post
func (c APIClient) Post(ctx context.Context, endpoint string, path string, options map[string]interface{}) (res []byte, err error) {
	httpCli := http.NewClient()
	baseURL, err := url.Parse("http://" + endpoint)
	if err != nil {
		return nil, err
	}
	baseURL.Path += path
	httpCli.Timeout(5 * time.Second).Post(baseURL.String()).SendMap(options)
	resp, body, errArr := httpCli.TraceEndByte(ctx)
	if len(errArr) > 0 {
		return nil, errs.InternalError(errArr[0].Error())
	}

	if resp.StatusCode >= 300 {
		return nil, errs.InternalError(string(body))
	}

	return body, nil
}

// Put Put
func (c APIClient) Put(ctx context.Context, endpoint string, path string, options map[string]interface{}) (res []byte, err error) {
	httpCli := http.NewClient()
	baseURL, err := url.Parse("http://" + endpoint)
	if err != nil {
		return nil, err
	}
	baseURL.Path += path
	httpCli.Timeout(5 * time.Second).Put(baseURL.String()).SendMap(options)
	resp, body, errArr := httpCli.TraceEndByte(ctx)
	if len(errArr) > 0 {
		return nil, errs.InternalError(errArr[0].Error())
	}

	if resp.StatusCode >= 300 {
		return nil, errs.InternalError(string(body))
	}

	return body, nil
}

// Delete Delete
func (c APIClient) Delete(ctx context.Context, endpoint string, path string, options map[string]interface{}) (res []byte, err error) {
	httpCli := http.NewClient()
	baseURL, err := url.Parse("http://" + endpoint)
	if err != nil {
		return nil, err
	}
	baseURL.Path += path
	params := url.Values{}
	for key, value := range options {
		params.Add(key, cast.ToString(value))
	}
	baseURL.RawQuery = params.Encode()
	httpCli.Timeout(5 * time.Second).Delete(baseURL.String())
	resp, body, errArr := httpCli.TraceEndByte(ctx)
	if len(errArr) > 0 {
		return nil, errs.InternalError(errArr[0].Error())
	}

	if resp.StatusCode >= 300 {
		return nil, errs.InternalError(string(body))
	}

	return body, nil
}
