package api

import "APP4/api/repository"

type OauthTwitterHandlers struct {
	Repo repository.RepoInterfaces
}

func NewOAuthTwitterHandlers(repo repository.RepoInterfaces) *OauthTwitterHandlers {
	return &OauthTwitterHandlers{
		Repo: repo}
}
