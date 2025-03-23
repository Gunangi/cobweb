package app

import (
	"fmt"
	"strings"
)

type MsgCode uint8

const (
	downloadedPageCode MsgCode = iota
	onPageCode
	visitingPageCode
)

// Please manintan the padding when ever you add a message
var codes map[MsgCode]string = map[MsgCode]string{
	downloadedPageCode: "Downloaded Page",
	onPageCode:         "On Page        ",
	visitingPageCode:   "Visiting Page  ",
}

type ApMsg interface {
	fmt.Stringer
	Title() string
	Msg() string
}

type DownloadedPage struct {
	URL string
}

func (p DownloadedPage) Title() string { return codes[downloadedPageCode] }
func (p DownloadedPage) Msg() string {
	return fmt.Sprintf("%s", p.URL)
}
func (p DownloadedPage) String() string { return p.Title() + "    " + p.Msg() }

type OnPage struct {
	PayLoad []string
}

func (v OnPage) Title() string  { return codes[onPageCode] }
func (v OnPage) Msg() string    { return strings.Join(v.PayLoad, "  ") }
func (v OnPage) String() string { return v.Title() + "    " + v.Msg() }

type VisitingPage struct {
	PayLoad []string
}

func (v VisitingPage) Title() string  { return codes[visitingPageCode] }
func (v VisitingPage) Msg() string    { return strings.Join(v.PayLoad, "  ") }
func (v VisitingPage) String() string { return v.Title() + "    " + v.Msg() }
