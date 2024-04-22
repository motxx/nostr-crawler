/*
Copyright © 2024 moti <motohiko.ave@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"nostr-crawler/pubsub"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:  "nostr-crawler",
	Long: "crawls nostr notes and publish them to a relay",
	RunE: func(cmd *cobra.Command, args []string) error {
		froms, err := cmd.Flags().GetStringArray("from")
		if err != nil {
			return err
		}
		if len(froms) == 0 {
			return fmt.Errorf("no relay 'from' specified")
		}
		if len(froms) > 1 {
			return fmt.Errorf("currently only one relay can be specified in 'from'")
		}

		tos, err := cmd.Flags().GetStringArray("to")
		if err != nil {
			return err
		}
		if len(tos) == 0 {
			return fmt.Errorf("no relay 'to' specified")
		}

		return pubsub.Loop(froms, tos)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringArrayP("from", "f", []string{}, "relay to crawl")
	rootCmd.Flags().StringArrayP("to", "t", []string{}, "relay to publish")
}
