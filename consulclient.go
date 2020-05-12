package consulclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// Client topic client
type Client struct {
	endpoint string
	cli      *http.Client
}

// HealthCheckService :
type HealthCheckService struct {
	Node HealthCheckNode `json:Node`
}

// HealthCheckNode :
type HealthCheckNode struct {
	Address string `json:Address`
}

func makeTransport(timeout time.Duration) http.RoundTripper {
	return &http.Transport{
		// inherited http.DefaultTransport
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: timeout,
	}
}

// New new client
func New(endpoint string, timeout time.Duration) *Client {
	cli := &http.Client{Transport: makeTransport(timeout)}
	return &Client{endpoint, cli}
}

// PassingNodes : get Passing nodes
func (c *Client) PassingNodes(ctx context.Context, service string) ([]HealthCheckNode, error) {
	var nodes []HealthCheckNode
	uri := fmt.Sprintf("%s/v1/health/service/%s?passing", c.endpoint, service)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nodes, err
	}
	req = req.WithContext(ctx)
	res, err := c.cli.Do(req)
	if err != nil {
		return nodes, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nodes, fmt.Errorf("Failed to publish msg:%v", res.Status)
	}

	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nodes, err
	}

	var srvs []HealthCheckService
	err = json.Unmarshal(d, &srvs)
	if err != nil {
		return nodes, err
	}

	// log.Printf("%v", srvs)
	for _, srv := range srvs {
		nodes = append(nodes, srv.Node)
	}

	return nodes, nil
}

// PassingHasIP service has ip
func (c *Client) PassingHasIP(ctx context.Context, service, ip string) (bool, error) {
	nodes, err := c.PassingNodes(ctx, service)
	if err != nil {
		return false, err
	}

	for _, node := range nodes {
		if node.Address == ip {
			return true, nil
		}
	}

	return false, nil
}

// PassingNodeLen service has ip
func (c *Client) PassingNodeLen(ctx context.Context, service string) (int, error) {
	nodes, err := c.PassingNodes(ctx, service)
	if err != nil {
		return 0, err
	}
	return len(nodes), nil
}
