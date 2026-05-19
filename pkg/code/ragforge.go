package code

import "github.com/morehao/golib/gerror"

const (
	KBCreateError           = 101100
	KBDeleteError           = 101101
	KBUpdateError           = 101102
	KBGetDetailError        = 101103
	KBGetPageListError      = 101104
	KBNotExistError         = 101105
	KnowledgeCreateError    = 101106
	KnowledgeDeleteError    = 101107
	KnowledgeUpdateError    = 101108
	KnowledgeGetDetailError = 101109
	KnowledgeGetPageListError = 101110
	KnowledgeNotExistError  = 101111
	ChunkDeleteError        = 101112
	ChunkSearchError        = 101113
	FAQCreateError          = 101114
	FAQDeleteError          = 101115
	FAQUpdateError          = 101116
	FAQSearchError          = 101117
	ChunkUpdateError        = 101118
	ChunkGetPageListError     = 101119
	FAQGetDetailError         = 101120
	FAQGetPageListError       = 101121
	SessionCreateError        = 101122
	SessionDeleteError        = 101123
	SessionUpdateError        = 101124
	SessionGetDetailError     = 101125
	SessionGetPageListError   = 101126
	MessageListError          = 101127
	MessageDeleteError        = 101128
	QAChatError               = 101129
	ModelCreateError           = 101130
	ModelDeleteError           = 101131
	ModelUpdateError           = 101132
	ModelGetDetailError        = 101133
	ModelGetPageListError      = 101134
	ModelTestError             = 101135
	VectorStoreCreateError     = 101136
	VectorStoreDeleteError     = 101137
	VectorStoreUpdateError     = 101138
	VectorStoreGetDetailError  = 101139
	VectorStoreGetPageListError = 101140
	VectorStoreTestError       = 101141
	TagCreateError             = 101142
	TagDeleteError             = 101143
	TagUpdateError             = 101144
	TagListError               = 101145
	KBCopyError                = 101146
	ChunkGetDetailError        = 101147
	SessionGenerateTitleError  = 101148
	MessageSearchError         = 101149
	KnowledgeDownloadError     = 101150
	FAQImportError             = 101151
	SessionStopError           = 101152
)

var ragforgeErrorMsgMap = gerror.CodeMsgMap{
	KBCreateError:              "创建知识库失败",
	KBDeleteError:              "删除知识库失败",
	KBUpdateError:              "修改知识库失败",
	KBGetDetailError:           "查看知识库详情失败",
	KBGetPageListError:         "查看知识库列表失败",
	KBNotExistError:            "知识库不存在",
	KnowledgeCreateError:       "创建知识失败",
	KnowledgeDeleteError:       "删除知识失败",
	KnowledgeUpdateError:       "修改知识失败",
	KnowledgeGetDetailError:    "查看知识详情失败",
	KnowledgeGetPageListError:  "查看知识列表失败",
	KnowledgeNotExistError:     "知识不存在",
	ChunkDeleteError:           "删除块失败",
	ChunkSearchError:           "搜索块失败",
	FAQCreateError:             "创建FAQ失败",
	FAQDeleteError:             "删除FAQ失败",
	FAQUpdateError:             "更新FAQ失败",
	FAQSearchError:             "搜索FAQ失败",
	ChunkUpdateError:           "更新块失败",
	ChunkGetPageListError:      "查看块列表失败",
	FAQGetDetailError:          "查看FAQ详情失败",
	FAQGetPageListError:        "查看FAQ列表失败",
	SessionCreateError:         "创建会话失败",
	SessionDeleteError:         "删除会话失败",
	SessionUpdateError:         "修改会话失败",
	SessionGetDetailError:      "查看会话详情失败",
	SessionGetPageListError:    "查看会话列表失败",
	MessageListError:           "查看消息列表失败",
	MessageDeleteError:         "删除消息失败",
	QAChatError:                   "问答对话失败",
	ModelCreateError:              "创建模型失败",
	ModelDeleteError:              "删除模型失败",
	ModelUpdateError:              "修改模型失败",
	ModelGetDetailError:           "查看模型详情失败",
	ModelGetPageListError:         "查看模型列表失败",
	ModelTestError:                "测试模型失败",
	VectorStoreCreateError:        "创建向量存储失败",
	VectorStoreDeleteError:        "删除向量存储失败",
	VectorStoreUpdateError:        "修改向量存储失败",
	VectorStoreGetDetailError:     "查看向量存储详情失败",
	VectorStoreGetPageListError:   "查看向量存储列表失败",
	VectorStoreTestError:          "测试向量存储失败",
	TagCreateError:                "创建标签失败",
	TagDeleteError:                "删除标签失败",
	TagUpdateError:                "修改标签失败",
	TagListError:                  "查看标签列表失败",
	KBCopyError:                   "复制知识库失败",
	ChunkGetDetailError:           "查看块详情失败",
	SessionGenerateTitleError:     "生成会话标题失败",
	MessageSearchError:            "搜索消息失败",
	KnowledgeDownloadError:        "下载知识失败",
	FAQImportError:                "导入FAQ失败",
	SessionStopError:              "停止会话失败",
}
