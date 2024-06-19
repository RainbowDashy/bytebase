package wecom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/bytebase/bytebase/backend/plugin/webhook"
)

// WebhookResponse is the API message for WeCom webhook response.
type WebhookResponse struct {
	ErrorCode    int    `json:"errcode"`
	ErrorMessage string `json:"errmsg"`
}

// WebhookMarkdown is the API message for WeCom webhook markdown.
type WebhookMarkdown struct {
	Content string `json:"content"`
}

// Webhook is the API message for WeCom webhook.
type Webhook struct {
	MessageType string          `json:"msgtype"`
	Markdown    WebhookMarkdown `json:"markdown"`
}

func init() {
	webhook.Register("bb.plugin.webhook.wecom", &Receiver{})
}

// Receiver is the receiver for WeCom.
type Receiver struct {
}

func (*Receiver) Post(context webhook.Context) error {
	metaStrList := []string{}
	for _, meta := range context.GetMetaList() {
		metaStrList = append(metaStrList, fmt.Sprintf("%s: <font color=\"comment\">%s</font>", meta.Name, meta.Value))
	}
	metaStrList = append(metaStrList, fmt.Sprintf("By: <font color=\"comment\">%s (%s)</font>", context.CreatorName, context.CreatorEmail))

	status := ""
	switch context.Level {
	case webhook.WebhookSuccess:
		status = "<font color=\"green\">Success</font> "
	case webhook.WebhookWarn:
		status = "<font color=\"yellow\">Warn</font> "
	case webhook.WebhookError:
		status = "<font color=\"red\">Error</font> "
	}
	content := fmt.Sprintf("# %s%s\n\n%s\n[View in Bytebase](%s)", status, context.Title, strings.Join(metaStrList, "\n"), context.Link)
	if context.Description != "" {
		content = fmt.Sprintf("# %s%s\n> %s\n\n%s\n[View in Bytebase](%s)", status, context.Title, context.Description, strings.Join(metaStrList, "\n"), context.Link)
	}

	post := Webhook{
		MessageType: "markdown",
		Markdown: WebhookMarkdown{
			Content: content,
		},
	}
	body, err := json.Marshal(post)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal webhook POST request to %s", context.URL)
	}
	req, err := http.NewRequest("POST",
		context.URL, bytes.NewBuffer(body))
	if err != nil {
		return errors.Wrapf(err, "failed to construct webhook POST request to %s", context.URL)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: webhook.Timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "failed to POST webhook to %s", context.URL)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "failed to read POST webhook response from %s", context.URL)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("failed to POST webhook to %s, status code: %d, response body: %s", context.URL, resp.StatusCode, b)
	}

	webhookResponse := &WebhookResponse{}
	if err := json.Unmarshal(b, webhookResponse); err != nil {
		return errors.Wrapf(err, "malformed webhook response from %s", context.URL)
	}

	if webhookResponse.ErrorCode != 0 {
		return errors.Errorf("%s", webhookResponse.ErrorMessage)
	}

	return nil
}
