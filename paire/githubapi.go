package paire

import (
	githubapi "github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"os"
	"strings"
	"errors"
	"net/http"
	"io"
)

type GithubApiInterface interface {
	Release(release Release) (Release, error)
	Download(tag string, dir string) (Release, error)
	Cleanup(release Release) error
}

type GithubApi struct {
	client *githubapi.Client
	Repository Repository
}

func NewGithubApi(repo Repository, token string) GithubApi {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	return GithubApi{
		client: githubapi.NewClient(tc),
		Repository: repo,
	}
}

func (api GithubApi) Release(release Release) (Release, error) {
	api.createTag(release)
	repoRelease, err := api.getReleaseByTag(release.Tag)
	if err == nil && *repoRelease.ID > 0 {
		api.deleteRelease(*repoRelease.ID)
	}
	githubRelease, err := api.createRelease(release)
	if err != nil {
		return release, errors.New("error creating release: " + err.Error())
	}
	release.PreRelease = *githubRelease.Prerelease
	release.Pushed = true
	release.Id = *githubRelease.ID
	for i, pkg := range release.Packages {
		err = api.uploadReleasePackage(release.Id, pkg.Name)
		if err == nil {
			release.Packages[i].Pushed = true
		}
	}

	return release, err
}

func (api GithubApi) createTag(release Release) error {
	message := "Paire"
	objectType := "commit"
	_, _, err := api.client.Git.CreateTag(
		api.Repository.Owner,
		api.Repository.Repository,
		&githubapi.Tag{
			Tag: &release.Tag,
			SHA: &release.Commit,
			Message: &message,
			Object: &githubapi.GitObject{
				SHA: &release.Commit,
				Type: &objectType,
			},
		},
	)
	return err
}

func (api GithubApi) getReleaseIdByTag(tag string) int {
	repo, _ := api.getReleaseByTag(tag)
	if repo != nil && *repo.ID > 0 {
		return *repo.ID
	}
	return 0
}

func (api GithubApi) getReleaseByTag(tag string) (*githubapi.RepositoryRelease, error) {
	repoRelease, _, err := api.client.Repositories.GetReleaseByTag(
		api.Repository.Owner,
		api.Repository.Repository,
		tag,
	)
	if err != nil {
		return nil, err
	}
	return repoRelease, nil
}

func (api GithubApi) deleteRelease(release int) bool {
	_, err := api.client.Repositories.DeleteRelease(
		api.Repository.Owner,
		api.Repository.Repository,
		release,
	)
	return err == nil
}

func (api GithubApi) createRelease(release Release) (*githubapi.RepositoryRelease, error) {
	preRelease := strings.HasPrefix(release.Commit, release.Tag)
	body := ""
	repoRelease, _, err := api.client.Repositories.CreateRelease(
		api.Repository.Owner,
		api.Repository.Repository,
		&githubapi.RepositoryRelease{
			TagName: &release.Tag,
			TargetCommitish: &release.Commit,
			Name: &release.Name,
			Prerelease: &preRelease,
			Body: &body,
		},
	)
	return repoRelease, err
}

func (api GithubApi) uploadReleasePackage(releaseId int, packageName string) error {
	if _, err := os.Stat(packageName); err != nil {
		return errors.New(packageName + " - " + err.Error())
	}
	file, err := os.Open(packageName)
	if err != nil {
		return errors.New(packageName + " - " + err.Error())
	}
	_, _, err = api.client.Repositories.UploadReleaseAsset(
		api.Repository.Owner,
		api.Repository.Repository,
		releaseId,
		&githubapi.UploadOptions{
			Name: packageName,
		},
		file,
	)
	return err
}

func (api GithubApi) Download(tag string, dir string) (Release, error) {
	repoRelease, err := api.getAssetsFromRelease(tag)
	if err != nil {
		return Release{}, err
	}

	assets := []ReleasePackage{};
	for _, repoAsset := range repoRelease.Assets {
		o, err := os.Create(dir + "/" + *repoAsset.Name);
		if err != nil {
			return Release{}, err
		}
		_, url, err := api.client.Repositories.DownloadReleaseAsset(api.Repository.Owner, api.Repository.Repository, *repoAsset.ID)
		o2, err := http.Get(url)
		if err != nil {
			return Release{}, err
		}
		_, err = io.Copy(o, o2.Body);
		if err != nil {
			return Release{}, err
		}
		assets = append(assets, ReleasePackage{
			Name: *repoAsset.Name,
			Pushed: true,
			Downloaded: true,
		})
	}

	return Release{
		Commit: *repoRelease.TargetCommitish,
		Tag: *repoRelease.TagName,
		Name: *repoRelease.TagName,
		Packages: assets,
		Pushed: true,
		PreRelease: *repoRelease.Prerelease,
		Id: *repoRelease.ID,
	}, nil
}

func (api GithubApi) getAssetsFromRelease(tag string) (*githubapi.RepositoryRelease, error) {
	release, err := api.getReleaseByTag(tag)
	if err != nil {
		return nil, err
	}
	if len(release.Assets) == 0 {
		err = errors.New("Release does not have any assets")
	}
	if err != nil {
		return release, err
	}
	return release, nil
}

func (api GithubApi) Cleanup(release Release) error {
	repoRelease, err := api.getReleaseByTag(release.Tag)
	if err != nil {
		return err
	}
	api.deleteRelease(*repoRelease.ID)
	api.deleteTag(release.Tag)
	return nil
}

func (api GithubApi) deleteTag(tag string) {
	api.client.Git.DeleteRef(api.Repository.Owner, api.Repository.Repository, tag)
}
