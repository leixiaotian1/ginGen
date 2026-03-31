package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/leixiaotian1/ginGen/internal/web"
	"github.com/spf13/cobra"
)

var (
	webPort int
	webHost string
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the web interface for ginGen",
	Long:  `Start a web server that provides a graphical interface for generating Gin projects with module selection.`,
	Run: func(cmd *cobra.Command, args []string) {
		router := web.SetupRouter()
		address := fmt.Sprintf("%s:%d", webHost, webPort)
		fmt.Printf("🚀 ginGen Web Interface is starting...\n")
		fmt.Printf("📱 Open your browser and visit: http://%s\n", address)
		fmt.Printf("Press Ctrl+C to stop the server\n\n")

		if err := http.ListenAndServe(address, router); err != nil {
			log.Fatalf("Failed to start web server: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.Flags().IntVarP(&webPort, "port", "p", 8080, "Port to run the web server on")
	webCmd.Flags().StringVar(&webHost, "host", "localhost", "Host to bind the web server to")
}
