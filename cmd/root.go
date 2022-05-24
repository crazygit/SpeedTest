package cmd

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
	"os"
	"speedtest/util"
	"strings"
)

var (
	url, output      string
	urlPromptContent = util.NewPromptContent(
		"The url to download?",
		"Please provide a valid url. for example: <http://212.183.159.230/5MB.zip>",
		util.OptionSetValidator(util.UrlValidator),
		util.OptionSetDefault("http://212.183.159.230/5MB.zip"),
	)
	outputPromptContent = util.NewPromptContent(
		"Where to save the download?",
		"Please provide a valid path. for example: /path/to/download/filename.ext",
		util.OptionSetDefault("5MB.zip"),
	)
	continuePromptContent = util.NewPromptContent(
		"Continue Test?",
		"Please choose Y/N",
		util.OptionSetValidator(util.YesNoValidator),
	)
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "speedTest",
	Short: "Download something from a given url to test download speed",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Args: func(cmd *cobra.Command, args []string) error {
	//	if len(args) < 1 {
	//		return errors.New("requires a url argument")
	//	}
	//	if err := IsValidUrl(args[0]); err != nil {
	//		return fmt.Errorf("invalid url specified: %s", args[0])
	//	} else {
	//		return nil
	//	}
	//},
	Run: func(cmd *cobra.Command, args []string) {
	START:
		if url == "" {
			url = util.PromptGetInput(urlPromptContent)
		}
		if output == "" {
			output = util.PromptGetInput(outputPromptContent)
		}
		util.Log.WithFields(logrus.Fields{
			"url":    url,
			"output": output,
		}).Info("Download Info")
		if err := downloadFile(url, output); err != nil {
			util.Log.Error(err)
		}
		continueTest := util.PromptGetInput(continuePromptContent)
		if strings.ToLower(continueTest) == "n" {
			os.Exit(0)
		} else {
			url = ""
			output = ""
			goto START
		}
	},
	Version: "0.01",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.speedtest.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cobra.MousetrapHelpText = "" // windows 下禁止提示This is a command line tool.
	rootCmd.Flags().StringVarP(&url, "url", "u", "", urlPromptContent.Label)
	rootCmd.Flags().StringVarP(&output, "output", "o", "", outputPromptContent.Label)
}

func downloadFile(url, output string) (err error) {
	// Create the file
	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			util.Log.Error(err)
		}
	}(f)

	//Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			util.Log.Error(err)
		}
	}(resp.Body)

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)
	// Writer the body to file
	_, err = io.Copy(io.MultiWriter(f, bar), resp.Body)
	if err != nil {
		return err
	}
	if err != nil {
		log.Fatal(err)
	}
	util.Log.WithField("costTimeInSeconds", bar.State().SecondsSince).Info("Download Success")
	return nil
}
