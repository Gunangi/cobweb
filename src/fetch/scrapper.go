package fetch

import (
	// "fmt"
	"net/url"
	"strings"

	"github.com/cloakwiss/cobweb/app"
	"github.com/gocolly/colly"
)

type PageTable map[url.URL]Page

type Page struct {
	Data []byte
	Metadata
}

type Metadata struct {
	// TODO: make this enum later
	MediaType string
	Title     string
}

// These are header collected from some pages have not figured out how to use it effectively
// Ideally these should not be required
var header = map[string][]string{
	// Obtained from Firefox Browser
	"Accept-Encoding": {"gzip", "deflate", "br", "zstd"},
	"Accept":          {"text/html"},
	"User-Agent": {
		// Obtained from Brave Browser
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36",
		// Obtained from Firefox Browser
		"Mozilla/5.0 (X11; Linux x86_64; rv:135.0) Gecko/20100101 Firefox/135.0",
	},
}

// This function should be called for all the target URL seprately
// or we need to find out if we can safely merge the graph of page
// discovered without redoing the work.
//
// For example if the urls are pointing to same domain but are of different
// path then we can share the downloaded result to some extend atleast.
//
//	And if that is the case then the index (of the book) generation will also need to change
//
// TODO: added some headers, caching and URL filter smarter
// TODO: there is a feature to use regex to mactch urls, need to find out how to expose it to users
func Scrapper(target url.URL, argu app.Options, out chan<- app.ApMsg) PageTable {
	defer close(out)
	pagesContents := make(PageTable)
	// println("Domain Name: ", target.String())
	// println("Depth: ", argu.Depth)

	rawAllowDomains := []url.URL{target}
	if len(argu.AllowDomains) > 0 {
		rawAllowDomains = append(rawAllowDomains, argu.AllowDomains...)
	}
	allowDomains := stringOfURL(rawAllowDomains)

	// recurse limit is unused
	collector := colly.NewCollector(
		colly.AllowedDomains(allowDomains...),
		// colly.DisallowedDomains(stringOfURL(argu.BlockDomains)...),
		colly.MaxDepth(int(argu.Depth)+1),
		colly.UserAgent(header["User-Agent"][1]),
	)

	collector.Limit(&colly.LimitRule{Parallelism: 3})

	//TODO: do not change the order of these callback methods
	// https://go-colly.org/docs/introduction/start/ read this
	//-------------------------------------------------------

	collector.OnRequest(func(r *colly.Request) {
		// r.Headers = (*http.Header)(&header)
		out <- app.VisitingPage{
			PayLoad: []string{r.URL.String()},
		}
		// fmt.Println("Visiting", r.URL.String())
	})

	//TODO: this cannot be left empty so what to do here
	collector.OnError(func(r *colly.Response, err error) {})

	// If my understanding is right then OnResponse is called when the Response is received
	// and if it is HTML then `OnHTML` is also called followed by `OnScraped` otherwise
	// `OnScraped` is called directly.
	// FIND: So, What does `OnScraped` allows to do which `OnResponse` cannot ?
	// I cannot find the much in docs and even they were not using `OnResponse` & `OnScraped`
	// in tandem as I am planning to do
	collector.OnResponse(func(res *colly.Response) {
		_, found := pagesContents[*res.Request.URL]
		if !found {
			pagesContents[*res.Request.URL] = Page{
				Data: res.Body,
				// Assign this properly
				Metadata: Metadata{
					//TODO: Title is only possible for HTML so what should be title other things
					Title:     "",
					MediaType: strings.Join((*res.Headers)["Content-Type"], " "),
				},
			}
			// size, err := strconv.ParseUint(res.Headers.Get("Content-Length"), 10, 32)
			// if err != nil {
			// 	size = 0
			// }
			out <- app.DownloadedPage{
				URL: res.Request.URL.String(),
			}
			// fmt.Printf("On page: %v\n", res.Request.URL)
		}
	})

	// Need to add others too
	{
		// TODO: need to find out how can I apply the filter in this function for various resources
		// There will be 2 layers of filtering:
		// 1. Will be based on extension of file based on URI
		// 2. Will be based on http header `Content-Type`
		//
		// Type 1 will be implemented here
		htmlHandler := func(tag string) colly.HTMLCallback {
			shouldAllow := assetFilterByExtension(argu)
			return func(e *colly.HTMLElement) {
				link := e.Attr(tag)
				if shouldAllow(link) {
					e.Request.Visit(link)
				}
			}
		}

		// TODO: verify if this is all
		collector.OnHTML("a[href]", htmlHandler("href"))
		collector.OnHTML("link[href]", htmlHandler("href"))
		collector.OnHTML("base[href]", htmlHandler("href"))
		collector.OnHTML("area[href]", htmlHandler("href"))

		collector.OnHTML("data[object]", htmlHandler("object"))

		// Need to see usage of these tags in the wild
		collector.OnHTML("del[cite]", htmlHandler("cite"))
		collector.OnHTML("ins[cite]", htmlHandler("cite"))
		collector.OnHTML("blockquote[cite]", htmlHandler("cite"))
		collector.OnHTML("q[cite]", htmlHandler("cite"))

		collector.OnHTML("img[src]", htmlHandler("src"))
		collector.OnHTML("track[src]", htmlHandler("src"))
		collector.OnHTML("embed[src]", htmlHandler("src"))
		collector.OnHTML("source[src]", htmlHandler("src"))
		collector.OnHTML("script[src]", htmlHandler("src"))
		collector.OnHTML("iframe[src]", htmlHandler("src"))
	}

	// collector.OnScraped(func(r *colly.Response) {
	// 	out <- app.ApMsg{
	// 		Code: app.OnScraped,
	// 		URL:  r.Request.URL.String(),
	// 	}
	// })

	// println("Started the scrapper")
	collector.Visit(target.String())

	//-------------------------------------------------------
	return pagesContents
}

// This function will produce the pattern to match for html handler
// based on arguments given, this is desgined to be called once only
//
// The assets that matches the patter will be excluded
func assetFilterByExtension(opts app.Options) func(string) bool {
	parts := make([]string, 0, 7)
	if opts.NoCss {
		parts = append(parts, ".css")
	}
	if opts.NoJs {
		parts = append(parts, ".js")
	}
	if opts.NoFonts {
		parts = append(parts, ".woff2")
		parts = append(parts, ".woff")
		parts = append(parts, ".otf")
		parts = append(parts, ".ttf")
		// This will make thing very dicey
		parts = append(parts, ".svg")
	}
	if opts.NoImages {
		parts = append(parts, ".png")
		parts = append(parts, ".jpeg")
		parts = append(parts, ".gif")
		parts = append(parts, ".webp")
		parts = append(parts, ".svg")
	}
	// Is it feasible to enumerate all the file formats here
	// if opts.Audio { }
	// if opts.Video { }

	// This will need some thing special and maybe should be covered with
	// another function but
	// if opts.Iframe { }
	return func(uri string) bool {
		for _, ext := range parts {
			if strings.HasSuffix(uri, ext) {
				return false
			}
		}
		return true
	}
}

func assetFilterByHeader(opts app.Options) {
	if opts.NoCss {
	}
	if opts.NoJs {
	}
	if opts.NoFonts {
		// This will make thing very dicey
	}
	if opts.NoImages {
	}

}

// Stick to url struct as much as possible
func removeProtocolPrefix(url *url.URL) string {
	var hostname string
	if url.Port() == "" {
		hostname = url.Hostname()
	} else {
		hostname = url.Hostname() + ":" + url.Port()
	}
	return strings.TrimLeft(hostname, "www.")
}

// WARN: this cannot handle paths at the moment
// TODO: what should the behaviour of this in case the it has a username in the url
func stringOfURL(urls []url.URL) []string {
	parsed := make([]string, 0, len(urls))
	for _, u := range urls {
		// u.String()
		if u.Scheme == "http" || u.Scheme == "https" {
			//WARN: This url's port is stripped here
			ur := u.Hostname()
			if u.Port() != "" {
				ur += ":"
				ur += u.Port()
			}
			parsed = append(parsed, ur)
		}
	}
	return parsed
}
