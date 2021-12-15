package repository

type DatabaseRepository interface {
	GetArticlesList()
	GetArticleById()
	GetCommentsByArticleId()
}
