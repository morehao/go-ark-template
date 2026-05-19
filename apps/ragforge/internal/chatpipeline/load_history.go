package chatpipeline

import (
	"context"

	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/glog"
)

type PluginLoadHistory struct{}

func NewPluginLoadHistory(em *EventManager) *PluginLoadHistory {
	p := &PluginLoadHistory{}
	em.Register(p)
	return p
}

func (p *PluginLoadHistory) ActivationEvents() []EventType {
	return []EventType{EventLoadHistory}
}

func (p *PluginLoadHistory) OnEvent(ctx context.Context, _ EventType, cm *ChatManage, next func() *PluginError) *PluginError {
	var messages []model.MessageEntity
	db := dbclient.RagForgeDB(ctx)
	if err := db.Where("session_id = ?", cm.SessionID).
		Order("created_at ASC").
		Find(&messages).Error; err != nil {
		glog.Errorf(ctx, "[PluginLoadHistory] query messages fail, err:%v, sessionID:%d", err, cm.SessionID)
		return next()
	}

	cm.History = groupMessagePairs(messages)

	return next()
}

func groupMessagePairs(messages []model.MessageEntity) []History {
	var pairs []History
	var current *History
	for _, msg := range messages {
		if msg.Role == model.MessageRoleUser {
			if current != nil && current.Answer != "" {
				pairs = append(pairs, *current)
			}
			current = &History{Query: msg.Content}
		} else if msg.Role == model.MessageRoleAssistant && current != nil {
			current.Answer = msg.Content
			pairs = append(pairs, *current)
			current = nil
		}
	}
	maxRounds := 5
	if len(pairs) > maxRounds {
		pairs = pairs[len(pairs)-maxRounds:]
	}
	return pairs
}
