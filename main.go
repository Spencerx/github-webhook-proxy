package main

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

var (
	jenkinsEndpoint string
	proxyEndpoint   string
)

var rootCmd = &cobra.Command{
	Use:          "github-proxy [sub]",
	Short:        "github-proxy",
	SilenceUsage: true,
}

var serverCmd = &cobra.Command{
	Use:          "server",
	Short:        "server",
	SilenceUsage: true,
	Run: func(c *cobra.Command, args []string) {
		listen := proxyEndpoint

		u, err := url.Parse(proxyEndpoint)
		if err == nil {
			listen = fmt.Sprintf("%s:%s", u.Hostname(), u.Port())
		}

		proxy := NewWebhookProxyServer(listen)
		proxy.ListenAndServe()
	},
}

var clientCmd = &cobra.Command{
	Use:          "client",
	Short:        "client",
	SilenceUsage: true,
	Run: func(c *cobra.Command, args []string) {
		client := NewWebhookProxyClient()
		client.Run()
	},
}

func main() {
	rootCmd.PersistentFlags().StringVar(&jenkinsEndpoint, "jenkins", "http://localhost:8080/ghprbhook/", "Jenkins ghprb endpoint")
	rootCmd.PersistentFlags().StringVar(&proxyEndpoint, "proxy", "http://127.0.0.1:8081", "address and port for the webhook proxy")
	rootCmd.AddCommand(clientCmd)
	rootCmd.AddCommand(serverCmd)
	rootCmd.Execute()
}
