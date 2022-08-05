package main

import (
	"encoding/json"
	"flag"
	"invite-brute/brute"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	nToGenerate  int
	proxyUrls    []*url.URL
	inviteUrl    string
	proxySelType brute.ProxySelectionType
	timeoutDelay time.Duration
	outPath      string
)

func main() {
	// Take in flags
	flag.Parse()
	// Create brute and generator
	b := brute.HttpBrute{
		BaseInviteUrl: inviteUrl,
		ProxyUrls:     proxyUrls,
		ProxySelType:  proxySelType,
	}
	g := brute.CodeGenerator{
		Chars: brute.DefaultChars[:],
	}
	// Generate codes and test all
	codes := g.GenerateCodes(nToGenerate, 5, 10)
	(*codes)[9] = "paevz3qZ"
	invites, err := b.Try(timeoutDelay, *codes...)
	if err != nil {
		panic(err)
	}
	// Pick out any failed invites
	parsedInvs := make([]*brute.Invite, 0, len(invites))
	for _, invite := range invites {
		if invite != nil {
			parsedInvs = append(parsedInvs, invite)
		}
	}
	// Convert to JSON and write to file
	cvb, err := json.MarshalIndent(parsedInvs, "", "\t") // Indent for readability
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile(outPath, os.O_WRONLY, 0664)
	if err != nil {
		panic(err)
	} else {
		defer f.Close()
	}
	_, err = f.Write(cvb)
	if err != nil {
		panic(err)
	}
}

func init() {
	// Parse flags
	nToGenerate = *flag.Int("codes", 10, "Number of codes to generate")
	for _, pU := range strings.Split((*flag.String("proxies", "", "The proxies to be used (seperate with ',')")), ",") {
		parsed, err := url.Parse(pU) // Parse all proxies to URL objects
		if err != nil || pU == "" {
			// Ignore if URL is invalid
			continue
		} else {
			proxyUrls = append(proxyUrls, parsed)
		}
	}
	inviteUrl = *flag.String("url", brute.DefaultBaseInviteUrl, "The base Discord invite url to be used: must include %s as code")
	// Switch case and determine selection type
	switch *flag.String("proxy_selection", "in_order", "The proxy selection order/type. 'random', 'in_order' or 'reverse'.") {
	case "in_order":
		proxySelType = brute.InOrderSelection
	case "random":
		proxySelType = brute.RandomSelection
	case "reverse":
		proxySelType = brute.ReverseOrderSelection
	default:
		proxySelType = brute.InOrderSelection
	}
	timeoutDelay = time.Duration(int(time.Second) * (*flag.Int("timeout_delay", 5, "The time that the program waits when it recieves a 429 status code. In seconds.")))
	outPath = *flag.String("out_path", ".\\out.json", "The JSON file path to write all invites to.")
}
