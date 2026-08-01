package router

import "github.com/morehao/golib/biz/gserver/ginserver"

func RegisterRouter(groups *ginserver.RouterGroups, appName string) {
	kbRouter(groups)
	knowledgeRouter(groups)
	chunkRouter(groups)
	faqRouter(groups)
	sessionRouter(groups)
	messageRouter(groups)
	qaRouter(groups)
	modelRouter(groups)
	vectorStoreRouter(groups)
	tagRouter(groups)
}
