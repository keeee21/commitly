package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ISlackGateway Slackゲートウェイのインターフェース
type ISlackGateway interface {
	SendMessage(ctx context.Context, webhookURL string, message *SlackMessage) error
}

// SlackMessage Slackメッセージ構造体
type SlackMessage struct {
	Text        string            `json:"text,omitempty"`
	Blocks      []SlackBlock      `json:"blocks,omitempty"`
	Attachments []SlackAttachment `json:"attachments,omitempty"`
}

// SlackBlock Slackブロック構造体
type SlackBlock struct {
	Type     string       `json:"type"`
	Text     *SlackText   `json:"text,omitempty"`
	Fields   []SlackText  `json:"fields,omitempty"`
	Elements []SlackText  `json:"elements,omitempty"` // context block用
}

// SlackText Slackテキスト構造体
type SlackText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// SlackAttachment Slack添付構造体
type SlackAttachment struct {
	Color  string `json:"color,omitempty"`
	Text   string `json:"text,omitempty"`
	Footer string `json:"footer,omitempty"`
}

type slackGateway struct {
	httpClient *http.Client
}

// NewSlackGateway コンストラクタ
func NewSlackGateway() ISlackGateway {
	return &slackGateway{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (g *slackGateway) SendMessage(ctx context.Context, webhookURL string, message *SlackMessage) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal slack message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send slack message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// エラーレスポンスを読み取る
		body := make([]byte, 1024)
		n, _ := resp.Body.Read(body)
		return fmt.Errorf("slack webhook returned non-200 status: %d, body: %s", resp.StatusCode, string(body[:n]))
	}

	return nil
}

// RivalCommitSummary ライバルのコミットサマリー
type RivalCommitSummary struct {
	Username string
	Commits  int
}

// MonthlyComparison 月次比較データ
type MonthlyComparison struct {
	CurrentMonth  int
	PreviousMonth int
}

// BuildWeeklyReportMessage 週次レポートメッセージを構築
func BuildWeeklyReportMessage(
	username string,
	userCommits int,
	rivals []RivalCommitSummary,
	startDate, endDate string,
) *SlackMessage {
	headerText := ":chart_with_upwards_trend: *Commitly 週次レポート*"

	// 自分のコミット数
	userStatsText := fmt.Sprintf(
		":bust_in_silhouette: *%s* の今週のコミット数: *%d*\n_（%s 〜 %s）_",
		username, userCommits, startDate, endDate,
	)

	blocks := []SlackBlock{
		{
			Type: "header",
			Text: &SlackText{
				Type: "plain_text",
				Text: "Commitly 週次レポート",
			},
		},
		{
			Type: "section",
			Text: &SlackText{
				Type: "mrkdwn",
				Text: userStatsText,
			},
		},
	}

	// ライバルの情報
	if len(rivals) > 0 {
		blocks = append(blocks, SlackBlock{
			Type: "divider",
		})

		rivalText := ":crossed_swords: *ライバルの今週のコミット数*\n"
		for i, rival := range rivals {
			emoji := getComparisonEmoji(userCommits, rival.Commits)
			rivalText += fmt.Sprintf("%d. %s %s: *%d* コミット\n", i+1, emoji, rival.Username, rival.Commits)
		}

		blocks = append(blocks, SlackBlock{
			Type: "section",
			Text: &SlackText{
				Type: "mrkdwn",
				Text: rivalText,
			},
		})
	}

	// フッター
	blocks = append(blocks, SlackBlock{
		Type: "context",
		Elements: []SlackText{
			{
				Type: "mrkdwn",
				Text: fmt.Sprintf("_Sent by Commitly at %s_", time.Now().Format("2006-01-02 15:04")),
			},
		},
	})

	return &SlackMessage{
		Text:   headerText,
		Blocks: blocks,
	}
}

// BuildMonthlyReportMessage 月次レポートメッセージを構築
func BuildMonthlyReportMessage(
	username string,
	comparison MonthlyComparison,
	rivals []RivalCommitSummary,
	monthLabel string,
) *SlackMessage {
	headerText := ":calendar: *Commitly 月次レポート*"

	// 前月との比較
	diff := comparison.CurrentMonth - comparison.PreviousMonth
	var diffText string
	var trendEmoji string

	if diff > 0 {
		trendEmoji = ":arrow_up:"
		diffText = fmt.Sprintf("+%d", diff)
	} else if diff < 0 {
		trendEmoji = ":arrow_down:"
		diffText = fmt.Sprintf("%d", diff)
	} else {
		trendEmoji = ":arrow_right:"
		diffText = "±0"
	}

	// 成長率
	var growthRate string
	if comparison.PreviousMonth > 0 {
		rate := float64(diff) / float64(comparison.PreviousMonth) * 100
		growthRate = fmt.Sprintf("（前月比 %+.1f%%）", rate)
	} else if comparison.CurrentMonth > 0 {
		growthRate = "（前月: 0コミット）"
	} else {
		growthRate = ""
	}

	userStatsText := fmt.Sprintf(
		":bust_in_silhouette: *%s* の%sのコミット数\n\n"+
			"• 今月: *%d* コミット\n"+
			"• 先月: *%d* コミット\n"+
			"• %s 差分: *%s* %s",
		username, monthLabel,
		comparison.CurrentMonth,
		comparison.PreviousMonth,
		trendEmoji, diffText, growthRate,
	)

	blocks := []SlackBlock{
		{
			Type: "header",
			Text: &SlackText{
				Type: "plain_text",
				Text: fmt.Sprintf("Commitly 月次レポート（%s）", monthLabel),
			},
		},
		{
			Type: "section",
			Text: &SlackText{
				Type: "mrkdwn",
				Text: userStatsText,
			},
		},
	}

	// ライバルの情報
	if len(rivals) > 0 {
		blocks = append(blocks, SlackBlock{
			Type: "divider",
		})

		rivalText := ":crossed_swords: *ライバルの今月のコミット数*\n"
		for i, rival := range rivals {
			emoji := getComparisonEmoji(comparison.CurrentMonth, rival.Commits)
			rivalText += fmt.Sprintf("%d. %s %s: *%d* コミット\n", i+1, emoji, rival.Username, rival.Commits)
		}

		blocks = append(blocks, SlackBlock{
			Type: "section",
			Text: &SlackText{
				Type: "mrkdwn",
				Text: rivalText,
			},
		})
	}

	// フッター
	blocks = append(blocks, SlackBlock{
		Type: "context",
		Elements: []SlackText{
			{
				Type: "mrkdwn",
				Text: fmt.Sprintf("_Sent by Commitly at %s_", time.Now().Format("2006-01-02 15:04")),
			},
		},
	})

	return &SlackMessage{
		Text:   headerText,
		Blocks: blocks,
	}
}

// getComparisonEmoji 比較結果に応じた絵文字を返す
func getComparisonEmoji(userCommits, rivalCommits int) string {
	if rivalCommits > userCommits {
		return ":fire:" // ライバルがリード
	} else if rivalCommits < userCommits {
		return ":muscle:" // 自分がリード
	}
	return ":handshake:" // 同点
}
