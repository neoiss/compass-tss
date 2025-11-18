package rpc

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"

	"testing"

	"github.com/ethereum/go-ethereum/rpc"
	. "gopkg.in/check.v1"
)

func TestPackage(t *testing.T) { TestingT(t) }

type RPCSuite struct{}

var _ = Suite(&RPCSuite{})

func (s *RPCSuite) TestRetry(c *C) {
	cl := Client{maxRetries: 3}
	called := 0
	err := cl.retry(func() error {
		called++
		return nil
	})
	c.Assert(err, IsNil)
	c.Assert(called, Equals, 1)

	called = 0
	err = cl.retry(func() error {
		called++
		return errors.New("error")
	})
	c.Assert(err, NotNil)
	c.Assert(called, Equals, 1)

	called = 0
	err = cl.retry(func() error {
		called++
		return errors.New("500 Internal Server Error: work queue depth exceeded")
	})
	c.Assert(err, NotNil)
	c.Assert(called, Equals, 4)

	called = 0
	err = cl.retry(func() error {
		called++
		if called < 2 {
			return errors.New("500 Internal Server Error: work queue depth exceeded")
		}
		return nil
	})
	c.Assert(err, IsNil)
	c.Assert(called, Equals, 2)
}

func TestClient_GetRawTransactionVerboseWithFee(t *testing.T) {
	type args struct {
		txid string
	}
	tests := []struct {
		name    string
		args    args
		want    *TxRawResult
		wantErr bool
	}{
		{
			name: "t-1",
			args: args{
				txid: "aa9fe0e29497cb3cbe14e3a4376937b91e69f8c1aa74f5df518e40eb20f39d62",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			host := os.Getenv("BITCOIN_NODE_HOST")
			user := os.Getenv("BITCOIN_NODE_USER")
			password := os.Getenv("BITCOIN_NODE_PASSWORD")
			authFn := func(h http.Header) error {
				auth := base64.StdEncoding.EncodeToString([]byte(user + ":" + password))
				h.Set("Authorization", fmt.Sprintf("Basic %s", auth))
				return nil
			}

			client, err := rpc.DialOptions(context.Background(), host, rpc.WithHTTPAuth(authFn))
			if err != nil {
				t.Fatal(err)
			}

			c := &Client{
				c: client,
			}
			got, err := c.GetRawTransactionVerboseWithFee(tt.args.txid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRawTransactionVerboseWithFee() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("GetRawTransactionVerboseWithFee() got = %v, want %v", got, tt.want)
			//}
			t.Logf("%.8f", got.Fee)
		})
	}
}
