package cmd

import (
	"errors"
	"log"

	"github.com/bakito/request-logger/pkg/handler"
	"github.com/spf13/cobra"
)

// forwardCmd represents the forward command
var forwardCmd = &cobra.Command{
	Use:   "forward <target URL>",
	Short: "Forward requests to a different URL, logging all requests and responses",
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a target url argument")
		}
		return nil
	},
	RunE: func(_ *cobra.Command, args []string) error {
		r := router()
		r.HandleFunc("/{path:.*}", handler.ForwardFor(args[0], disableLogger, withTLS()))

		log.Printf("Forwarding requests to %s", args[0])
		start(r)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(forwardCmd)
}
