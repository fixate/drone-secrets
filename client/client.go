package client

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackspirou/syscerts"
	"github.com/urfave/cli"
	"golang.org/x/net/proxy"
	"golang.org/x/oauth2"

	"github.com/drone/drone-go/drone"
)

// NewClient returns a new client from the CLI context.
func NewClient(c *cli.Context) (drone.Client, error) {
	var (
		skip     = c.GlobalBool("skip-verify")
		socks    = c.GlobalString("socks-proxy")
		socksoff = c.GlobalBool("socks-proxy-off")
		token    = c.GlobalString("token")
		server   = c.GlobalString("server")
	)
	server = strings.TrimRight(server, "/")

	// if no server url is provided we can default
	// to the hosted Drone service.
	if len(server) == 0 {
		return nil, fmt.Errorf("Error: you must provide the Drone server address.")
	}
	if len(token) == 0 {
		return nil, fmt.Errorf("Error: you must provide your Drone access token.")
	}

	// attempt to find system CA certs
	certs := syscerts.SystemRootsPool()
	tlsConfig := &tls.Config{
		RootCAs:            certs,
		InsecureSkipVerify: skip,
	}

	config := new(oauth2.Config)
	auther := config.Client(
		oauth2.NoContext,
		&oauth2.Token{
			AccessToken: token,
		},
	)

	trans, _ := auther.Transport.(*oauth2.Transport)

	if len(socks) != 0 && !socksoff {
		dialer, err := proxy.SOCKS5("tcp", socks, nil, proxy.Direct)
		if err != nil {
			return nil, err
		}
		trans.Base = &http.Transport{
			TLSClientConfig: tlsConfig,
			Proxy:           http.ProxyFromEnvironment,
			Dial:            dialer.Dial,
		}
	} else {
		trans.Base = &http.Transport{
			TLSClientConfig: tlsConfig,
			Proxy:           http.ProxyFromEnvironment,
		}
	}

	return drone.NewClient(server, auther), nil
}
