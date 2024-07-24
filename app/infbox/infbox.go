package infbox

import (
	"github.com/YumaFuu/ssm-tui/app/pubsub"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type InfoBox struct {
	*tview.TextView
	pubsub *pubsub.PubSub
}

func NewInfoBox(ps *pubsub.PubSub) *InfoBox {
	v := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true)

	v.
		SetBackgroundColor(tcell.ColorDefault).
		SetBorder(true)

	return &InfoBox{v, ps}
}

func (v InfoBox) WaitTopic() {
	chUpdate := v.pubsub.Sub(pubsub.TopicUpdateInfoBox)

	for {
		select {
		case msg := <-chUpdate:
			if s, ok := msg.(string); ok {
				v.SetText(s)
			}
		}
	}
}
