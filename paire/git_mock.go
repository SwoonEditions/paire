package paire

import (
	"github.com/stretchr/testify/mock"
)

type GitMock struct {
	mock.Mock
}

func (git GitMock) CurrentCommit() string {
	return git.Called().String(0)

}

func (git GitMock) CurrentTag() string {
	return git.Called().String(0)
}

func (git GitMock) CurrentRemotes() string {
	return git.Called().String(0)
}

func (git GitMock) Shorten(sha string) string {
	return git.Called(sha).String(0)
}