package paire

import "github.com/stretchr/testify/mock"

type GithubApiMock struct {
	mock.Mock
}

func (api *GithubApiMock) Release(release Release) (Release, error) {
	args := api.Called(release)
	return args.Get(0).(Release), args.Error(1)
}

func (api *GithubApiMock) Download(tag string, dir string) (Release, error) {
	args := api.Called(tag)
	return args.Get(0).(Release), args.Error(1)
}

func (api *GithubApiMock) Cleanup(release Release) error {
	args := api.Called(release)
	return args.Error(0)
}
