package brute

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

// A default brute with default values.
var DefBrute = &HttpBrute{
	BaseInviteUrl: DefaultBaseInviteUrl,
	ProxySelType:  InOrderSelection,
}

const (
	// The normal Discord invite url.
	// The only time you want to stray away from using this is for proxies.
	DefaultBaseInviteUrl = "https://discord.com/api/v9/invites/%s?with_counts=true&with_expiration=true"
)

// Allows for distinguishing between different proxy selection types.
type ProxySelectionType int64

// The seperate ProxySelectionTypes.
const (
	RandomSelection ProxySelectionType = iota
	InOrderSelection
	ReverseOrderSelection
)

type Guild struct {
	// Guild identifier.
	Id string `json:"id"`
	// The name of the guild
	Name string `json:"name"`
	// The nsfw status of the guild
	Nsfw bool `json:"nsfw"`
}

// Holds inviter (invite creator) data.
type Inviter struct {
	// User identifier.
	Id string `json:"id"`
	// The creator's username.
	Username string `json:"username"`
	// The creator's discriminator.
	Discriminator string `json:"discriminator"`
}

// Holds invite data.
type Invite struct {
	// The invite's code.
	Code string `json:"code"`
	// The date and time that the invite expires at.
	ExpiresAt string `json:"expires_at"`
	// The date and time in which the invite was made at.
	CreatedAt string `json:"created_at"`
	// The guild of the invite.
	Guild Guild `json:"guild"`
	// The creator of the invite.
	Inviter Inviter `json:"inviter"`
}

// A struct hosting methods and values for brute-testing Discord invite codes.
type HttpBrute struct {
	// The base Discord invite url.
	// Must include a %s as placement for the invite code.
	BaseInviteUrl string
	// The proxy URL(s) to be used.
	ProxyUrls []*url.URL
	// The proxy selection type.
	ProxySelType ProxySelectionType
}

/*
	Returns a sorted slice of Transport instances containing the proxies defined
	on object creation.

`selType`, if provided as nil, defaults to the brute instance's proxy
selection type.

The returned slice is of pointers to transport instances. If you wish to modify
those transport instances without the use of a for-loop post the method's
execution, you can pass a "model" transport instance. This will be used as a
base template for all generated transport instances prior to adding the proxies.
*/
func (b *HttpBrute) sortAndFormatProxies(selType *ProxySelectionType, template *http.Transport) (tPs []*http.Transport) {
	if len(b.ProxyUrls) < 1 {
		return []*http.Transport{}
	}
	if template == nil {
		template = &http.Transport{}
	}
	if selType == nil {
		selType = &b.ProxySelType
	}
	// A lambda that deep copies a http transport object - VERY EXPENSIVE
	deepCopy := func(src, dist *http.Transport) (err error) {
		buf := bytes.Buffer{}
		if err = gob.NewEncoder(&buf).Encode(src); err != nil {
			return
		}
		return gob.NewDecoder(&buf).Decode(dist)
	}
	// Just for simplification
	generateTransport := func(proxy *url.URL) (transport *http.Transport) {
		var tP *http.Transport
		_ = deepCopy(template, tP)
		tP.Proxy = http.ProxyURL(proxy)
		return tP
	}
	// Switch and case for each proxy selection type
	proxies := make([]*http.Transport, len(b.ProxyUrls))
	switch *selType {
	case RandomSelection: // Creates the illusion of randomness?
		rand.Shuffle(len(b.ProxyUrls), func(i, j int) {
			proxies[i], proxies[j] = generateTransport(b.ProxyUrls[j]), generateTransport(b.ProxyUrls[i])
		})
	case InOrderSelection:
		// Loop through all proxies and convert to transport
		for _, proxyUrl := range b.ProxyUrls {
			proxies = append(proxies, generateTransport(proxyUrl))
		}
	case ReverseOrderSelection:
		// Same as InOrderSelection, but reverse in end
		for _, proxyUrl := range b.ProxyUrls {
			proxies = append(proxies, generateTransport(proxyUrl))
		}
		for i, j := 0, len(proxies)-1; i < j; i, j = i+1, j-1 { // Reverse code!
			proxies[i], proxies[j] = proxies[j], proxies[i]
		}
	}
	return tPs
}

/*
	Attempts to fetch the invite data of all codes provided.

For each code, an Invite object is generated, and its pointer value added to
a slice. If the pointer is equal to nil, the invite does not exist.

`timeoutDelay` is the amount of time that is waited for a proxy when it is on
timeout (usually because too many requests have been sent through it).

An error is only raised by HTTP error, or when a proxy connection fails.
*/
func (b *HttpBrute) Try(timeoutDelay time.Duration, codes ...string) (invites []*Invite, err error) {
	invites = make([]*Invite, 0, len(codes))
	transports := b.sortAndFormatProxies(nil, nil)
	timeouts := make([]bool, len(transports))
	proxiesOn := len(transports) > 0
	tI := 0
	for _, code := range codes {
		if proxiesOn { // If proxies exist
			if tI > len(transports)-1 {
				tI = 0
			}
			// Check for timeout
			if timeouts[tI] {
				time.Sleep(timeoutDelay)
				timeouts[tI] = false
			}
			// Set default transport and send get request
			http.DefaultClient.Transport = transports[tI]
		}
		resp, err := http.Get(fmt.Sprintf(b.BaseInviteUrl, code))
		if err != nil {
			return invites, err
		}
		// Handle status codes
		if resp.StatusCode == 429 {
			invites = append(invites, nil)
			if proxiesOn {
				timeouts[tI] = true
			} else {
				time.Sleep(timeoutDelay) // Sleep there
			}
			continue
		} else if resp.StatusCode == 404 {
			invites = append(invites, nil)
			continue
		} else if resp.StatusCode != 200 {
			invites = append(invites, nil)
			return invites, errors.New(resp.Status)
		}
		// Use JSON to unmarshal received data
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return invites, err
		}
		invite := &Invite{}
		if err = json.Unmarshal(b, invite); err != nil {
			return invites, err
		} else {
			invites = append(invites, invite)
		}
		tI++ // Cycle through to next transport
	}
	return invites, nil
}
