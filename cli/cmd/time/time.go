package time

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "Returns the time from the central Time API",
	Long: `The time command gets the current time that is used by 
	all of the other services`,
	RunE: Run,
}

var timeApiURL string

type TimeResponse struct {
	Time string `json:time`
}

func TimeCmd() *cobra.Command {
	return timeCmd
}

func init() {
	timeApiURL = os.Getenv("TIME_API_URL")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Run(cmd *cobra.Command, args []string) error {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(timeApiURL)
	if err != nil {
		return fmt.Errorf("failed to get the time: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to get read response body: %w", err)
	}

	var tBody TimeResponse
	err = json.Unmarshal(body, &tBody)
	if err != nil {
		return fmt.Errorf("failed to get parse response body: %w", err)
	}

	fmt.Printf("%s\n", tBody.Time)
	return nil
}
