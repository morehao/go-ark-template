package retriever

var globalRetriever Retriever

func Init(r Retriever) {
	globalRetriever = r
}

func Get() Retriever {
	return globalRetriever
}

func NewRetriever() Retriever {
	return &pgRetriever{}
}
