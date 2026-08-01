package code

import "github.com/morehao/golib/gerror"

const (
	KBCreateError               = 101300
	KBDeleteError               = 101301
	KBUpdateError               = 101302
	KBGetDetailError            = 101303
	KBGetPageListError          = 101304
	KBNotExistError             = 101305
	KnowledgeCreateError        = 101306
	KnowledgeDeleteError        = 101307
	KnowledgeUpdateError        = 101308
	KnowledgeGetDetailError     = 101309
	KnowledgeGetPageListError   = 101310
	KnowledgeNotExistError      = 101311
	ChunkDeleteError            = 101312
	ChunkSearchError            = 101313
	FAQCreateError              = 101314
	FAQDeleteError              = 101315
	FAQUpdateError              = 101316
	FAQSearchError              = 101317
	ChunkUpdateError            = 101318
	ChunkGetPageListError       = 101319
	FAQGetDetailError           = 101320
	FAQGetPageListError         = 101321
	SessionCreateError          = 101322
	SessionDeleteError          = 101323
	SessionUpdateError          = 101324
	SessionGetDetailError       = 101325
	SessionGetPageListError     = 101326
	MessageListError            = 101327
	MessageDeleteError          = 101328
	QAChatError                 = 101329
	ModelCreateError            = 101330
	ModelDeleteError            = 101331
	ModelUpdateError            = 101332
	ModelGetDetailError         = 101333
	ModelGetPageListError       = 101334
	ModelTestError              = 101335
	VectorStoreCreateError      = 101336
	VectorStoreDeleteError      = 101337
	VectorStoreUpdateError      = 101338
	VectorStoreGetDetailError   = 101339
	VectorStoreGetPageListError = 101340
	VectorStoreTestError        = 101341
	TagCreateError              = 101342
	TagDeleteError              = 101343
	TagUpdateError              = 101344
	TagListError                = 101345
	KBCopyError                 = 101346
	ChunkGetDetailError         = 101347
	SessionGenerateTitleError   = 101348
	MessageSearchError          = 101349
	KnowledgeDownloadError      = 101350
	FAQImportError              = 101351
	SessionStopError            = 101352
)

var ragforgeErrorMsgMap = gerror.CodeMsgMap{
	KBCreateError:               "创建知识库失败",
	KBDeleteError:               "删除知识库失败",
	KBUpdateError:               "修改知识库失败",
	KBGetDetailError:            "查看知识库详情失败",
	KBGetPageListError:          "查看知识库列表失败",
	KBNotExistError:             "知识库不存在",
	KnowledgeCreateError:        "创建知识失败",
	KnowledgeDeleteError:        "删除知识失败",
	KnowledgeUpdateError:        "修改知识失败",
	KnowledgeGetDetailError:     "查看知识详情失败",
	KnowledgeGetPageListError:   "查看知识列表失败",
	KnowledgeNotExistError:      "知识不存在",
	ChunkDeleteError:            "删除块失败",
	ChunkSearchError:            "搜索块失败",
	FAQCreateError:              "创建FAQ失败",
	FAQDeleteError:              "删除FAQ失败",
	FAQUpdateError:              "更新FAQ失败",
	FAQSearchError:              "搜索FAQ失败",
	ChunkUpdateError:            "更新块失败",
	ChunkGetPageListError:       "查看块列表失败",
	FAQGetDetailError:           "查看FAQ详情失败",
	FAQGetPageListError:         "查看FAQ列表失败",
	SessionCreateError:          "创建会话失败",
	SessionDeleteError:          "删除会话失败",
	SessionUpdateError:          "修改会话失败",
	SessionGetDetailError:       "查看会话详情失败",
	SessionGetPageListError:     "查看会话列表失败",
	MessageListError:            "查看消息列表失败",
	MessageDeleteError:          "删除消息失败",
	QAChatError:                 "问答对话失败",
	ModelCreateError:            "创建模型失败",
	ModelDeleteError:            "删除模型失败",
	ModelUpdateError:            "修改模型失败",
	ModelGetDetailError:         "查看模型详情失败",
	ModelGetPageListError:       "查看模型列表失败",
	ModelTestError:              "测试模型失败",
	VectorStoreCreateError:      "创建向量存储失败",
	VectorStoreDeleteError:      "删除向量存储失败",
	VectorStoreUpdateError:      "修改向量存储失败",
	VectorStoreGetDetailError:   "查看向量存储详情失败",
	VectorStoreGetPageListError: "查看向量存储列表失败",
	VectorStoreTestError:        "测试向量存储失败",
	TagCreateError:              "创建标签失败",
	TagDeleteError:              "删除标签失败",
	TagUpdateError:              "修改标签失败",
	TagListError:                "查看标签列表失败",
	KBCopyError:                 "复制知识库失败",
	ChunkGetDetailError:         "查看块详情失败",
	SessionGenerateTitleError:   "生成会话标题失败",
	MessageSearchError:          "搜索消息失败",
	KnowledgeDownloadError:      "下载知识失败",
	FAQImportError:              "导入FAQ失败",
	SessionStopError:            "停止会话失败",
}
