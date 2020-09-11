package zk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/utils"
)

type BaseClient struct {
	Endpoints []string
	HTTP      *http.Client
	Endpoint  string
	Transport *http.Transport
}

func (c *BaseClient) Get(ctx context.Context, pathWithQuery string, out interface{}) error {
	return c.request(ctx, http.MethodGet, pathWithQuery, nil, out)
}

func (c *BaseClient) Post(ctx context.Context, pathWithQuery string, in, out interface{}) error {
	return c.request(ctx, http.MethodPost, pathWithQuery, in, out)
}

func (c *BaseClient) request(
	ctx context.Context,
	method string,
	pathWithQuery string,
	requestObj,
	responseObj interface{},
) error {
	var body io.Reader = http.NoBody
	if requestObj != nil {
		outData, err := json.Marshal(requestObj)
		if err != nil {
			return err
		}
		body = bytes.NewBuffer(outData)
	}

	request, err := http.NewRequest(method, utils.Joins(c.Endpoint, pathWithQuery), body)
	if err != nil {
		return err
	}
	resp, err := c.doRequest(ctx, request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if responseObj != nil {
		if err := json.NewDecoder(resp.Body).Decode(responseObj); err != nil {
			return err
		}
	}

	return nil
}

func (c *BaseClient) doRequest(context context.Context, request *http.Request) (*http.Response, error) {
	withContext := request.WithContext(context)
	withContext.Header.Set("Content-Type", "application/json; charset=utf-8")

	response, err := c.HTTP.Do(withContext)
	if err != nil {
		return response, err
	}
	err = checkError(response)
	return response, err
}

func checkError(response *http.Response) error {
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return errors.New("get zk agent fail pls check")
	}
	return nil
}

func (c *BaseClient) Close() {
	if c.Transport != nil {
		// When the http transport goes out of scope, the underlying goroutines responsible
		// for handling keep-alive connections are not closed automatically.
		// Since this client gets recreated frequently we would effectively be leaking goroutines.
		// Let's make sure this does not happen by closing idle connections.
		c.Transport.CloseIdleConnections()
	}
}

func (c *BaseClient) Equal(c2 *BaseClient) bool {
	// handle nil case
	if c2 == nil && c != nil {
		return false
	}

	// compare endpoint and user creds
	return c.Endpoint == c2.Endpoint
}

func (c *BaseClient) IsAlive(c2 *BaseClient) bool {
	// handle nil case
	if c2 == nil && c != nil {
		return false
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	_, err := c.GetClusterUp(timeoutCtx)
	return err == nil
}
