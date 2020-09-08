import (
	"fmt"
	"strings"
	"log"
	"net/http"
	"os"

	"k8s.io/klog"
)

type Client struct {
	token      string
	endpoint   string
	debug      bool
	log        ilogger
	httpclient httpClient
}

func toss(channel, combinedNote) {
	baseURL := "https://api.telegram.org/bot" + Client.token
	payload := url.Values{}
	payload.Add("chat_id", channel)
	payload.Add("text", combinedNote)

	newURL := baseURL + "/sendMessage?" + payload.Encode()
	resp1, err1 := http.Get(newURL)
	if err1 != nil {
		klog.Warningf("Failed to notify Telegram: %v", err1)
	}
	defer resp1.Body.Close()
}

func New(token string, options ...Option) *Client {
	s := &Client{
		token:      token,
		httpclient: &http.Client{},
		log:        log.New(os.Stderr, "minitelclient", log.LstdFlags|log.Lshortfile),
	}

	for _, opt := range options {
		opt(s)
	}

	return s
}

func stitch(note *CrashNotification) {
	s:= []strings{
		"**Title**",
		"```",
		note.Title,
		"```",
		"**Message**",
		"```",
		note.Message,
		"```",
		"**Logs**"
		"```",
		note.Logs,
		"```"
	}

	if note.Reason != "" {
		t := []strings{
			"**Reason**",
			"```",
			note.Reason,
			"```"
		}
		s = append(s, t)
	}
	return strings.Join(s, "\n")
}