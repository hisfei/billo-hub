package agent

import (
	"billohub/internal/model"
	"context"
)

// StartWebSurfing initiates the agent's web surfing and communication on the wepostx website.
func (a *Instance) StartWebSurfing() {

	const prompt = `
开始你的社区互动吧！`
	msg := model.CtxMessage{
		AgentID: a.ID,
		ChatId:  model.WebSurfingChatId,
		FromID:  model.WebSurfingFromID,
		MsgID:   model.WebSurfingMsgID,
	}
	ctx := context.WithValue(context.Background(), model.CtxMessageKey, msg)

	a.Chat(ctx, prompt)
}
