package app

import (
	"net/url"
	"time"

	"github.com/spf13/cobra"
)

type Mode uint8

// FIX: make this some thing like enum
const (
	Silent Mode = iota
	Default
	Debug
)

type Options struct {
	// exclude the following or not
	Audio, Css, Iframe, Fonts, Js, Images, Video, Metadata bool
	Targets, AllowDomains, BlockDomains                    []url.URL
	Output, Cookie                                         string
	Depth                                                  uint8
	Mode                                                   Mode
	Timeout                                                time.Duration
}

func addURLs(raw []string) []url.URL {
	// log.Printf("Input: %+v", *to)
	// defer log.Printf("Output: %+v", *to)
	parsed := make([]url.URL, 0, len(raw))
	for _, ur := range raw {
		u, err := url.Parse(ur)
		if err != nil {
			return nil
		}
		parsed = append(parsed, *u)
	}
	return parsed
}

func Args() Options {
	var args Options
	// This should accept more than 1 urls
	rootCmd := &cobra.Command{
		Use:           "cobweb",
		Short:         "Download Webpages as ePUB",
		Long:          "Mimimal, friendly tool to download webpage as ePUB",
		SilenceErrors: false,
		SilenceUsage:  false,
		Args:          cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, arg []string) error {
			// fmt.Printf("%+v\n", arg)
			// println("Arg: ", len(arg))
			if len(arg) >= 1 {
				// log.Println("Entered")
				ur, err := url.Parse(arg[0])
				// fmt.Printf("Parsed: %+v\n", ur)
				if err == nil {
					args.Targets = append(args.Targets, *ur)
				} else {
					panic(err)
				}
				return nil
			} else {
				return cmd.Help()
			}
		},
	}

	rootCmd.Flags().StringVarP(&args.Output, "output", "O", "", "Name of output epub file (Same as the page's title)")
	rootCmd.Flags().Uint8VarP(&args.Depth, "depth", "d", 0, "Maximum recursion")

	// Copied from monolith => Maybe reorder it a bit
	// Exclusion flags for various elements
	rootCmd.Flags().BoolVarP(&args.Js, "no-js", "j", false, "Remove JavaScript")
	rootCmd.Flags().BoolVarP(&args.Css, "no-css", "c", false, "Remove CSS")
	rootCmd.Flags().BoolVarP(&args.Images, "no-images", "i", false, "Remove images")
	rootCmd.Flags().BoolVarP(&args.Fonts, "no-fonts", "f", false, "Remove fonts")
	rootCmd.Flags().BoolVarP(&args.Audio, "no-audio", "a", false, "Remove audio sources")
	rootCmd.Flags().BoolVarP(&args.Video, "no-video", "V", false, "Remove video sources")
	rootCmd.Flags().BoolVarP(&args.Metadata, "no-metadata", "m", false, "Exclude timestamp and source information")
	//TODO: Check if this is possible in EPUB
	// rootCmd.Flags().BoolP("no-frames", "F", false, "Remove frames and iframes")

	//TODO: Will need Need to lookup these flags' behaviour
	// rootCmd.Flags().BoolP("ignore-errors", "e", false, "Ignore network errors")
	// rootCmd.Flags().BoolP("isolate", "I", false, "Cut off document from the Internet")
	// rootCmd.Flags().BoolP("insecure", "k", false, "Allow invalid X.509 (TLS) certificates")
	// rootCmd.Flags().BoolP("unwrap-noscript", "n", false, "Replace NOSCRIPT elements with their contents")

	// Remaing flags from monolith
	// rootCmd.Flags().BoolP("base-url" <http://localhost/> "b",   "Set custom base URL")
	// rootCmd.Flags().BoolP("user-agent",           "u",   "Set custom User-Agent string")
	allow := rootCmd.Flags().StringArrayP("allow-domain", "A", []string{}, "Specify domains to use for white/black-listing")
	// if err := addURLs(allow, &args.AllowDomains); err != nil {
	// 	panic("Allow Domain addition failed")
	// }
	if ans := addURLs(*allow); ans != nil {
		args.AllowDomains = ans
	}
	block := rootCmd.Flags().StringArrayP("block-domains", "D", []string{}, "Treat list of specified domains as blacklist")
	if ans := addURLs(*block); ans != nil {
		args.BlockDomains = ans
	}

	rootCmd.Flags().StringVarP(&args.Cookie, "cookies", "C", "", "Specify cookie file")
	rootCmd.Flags().DurationVarP(&args.Timeout, "timeout", "T", 60*time.Second, "Adjust network request timeout")

	rootCmd.Flags().BoolP("silent", "s", false, "Suppress verbosity")
	rootCmd.Flags().BoolP("version", "v", false, "Print version information")

	rootCmd.Flags().SortFlags = false

	rootCmd.Execute()

	// Filling the remaining arguments
	//TODO: validate output file name

	return args
}
