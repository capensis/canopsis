package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	ini "github.com/vaughan0/go-ini"
)

type Args struct {
	Token   *string
	APIUrl  *string
	Tag     *string
	INI     *string
	Cleanup *bool
}

type GitLab struct {
	Token  string
	APIUrl string
	Client *http.Client
}

func (g *GitLab) NewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("PRIVATE-TOKEN", g.Token)
	return req, nil
}

func (g *GitLab) URLEncode(project_path string) string {
	return url.PathEscape(project_path)
}

func (g *GitLab) DelTag(tag string, project_path string) error {
	pp := g.URLEncode(project_path)
	tag = g.URLEncode(tag)
	urlpath := fmt.Sprintf(
		"%s/projects/%s/repository/tags/%s", g.APIUrl, pp, tag,
	)
	req, err := g.NewRequest("DELETE", urlpath, nil)
	if err != nil {
		return err
	}
	r, err := g.Client.Do(req)
	if err != nil {
		return err
	}

	if r.StatusCode == 404 || r.StatusCode == 204 {
		return nil
	}

	return errors.New(fmt.Sprintf("cannot delete tag: %s", r.Status))
}

func (g *GitLab) AddTag(tag string, ref string, project_path string) error {
	pp := g.URLEncode(project_path)
	tag = g.URLEncode(tag)
	urlpath := fmt.Sprintf(
		"%s/projects/%s/repository/tags", g.APIUrl, pp,
	)

	data := url.Values{}
	data.Add("id", pp)
	data.Add("tag_name", tag)
	data.Add("ref", ref)

	req, err := g.NewRequest("POST", urlpath, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	r, err := g.Client.Do(req)
	if err != nil {
		return err
	}

	if r.StatusCode != 201 {
		return errors.New(fmt.Sprintf("cannot create tag: %s", r.Status))
	}
	return nil
}

func args() Args {
	a := Args{
		Token:   flag.String("token", "", "GitLab private token"),
		APIUrl:  flag.String("apiurl", "https://git.canopsis.net/api/v4", "gitlab api url"),
		Tag:     flag.String("tag", "", "tag name"),
		INI:     flag.String("ini", "catag.ini", "path to ini file"),
		Cleanup: flag.Bool("cleanup", false, "do not add tags, only remove the specified one"),
	}

	flag.Parse()

	return a
}

func work(g *GitLab, args Args, project_path string, ref string) {
	tag := *args.Tag
	fmt.Printf("working on %s:%s: ", project_path, ref)
	if err := g.DelTag(tag, project_path); err != nil {
		fmt.Printf("%s\n", err)
	} else if !*args.Cleanup {
		if err := g.AddTag(tag, ref, project_path); err != nil {
			fmt.Printf("%s\n", err)
		}
	}
	fmt.Println("ok")
}

func main() {
	args := args()

	g := GitLab{
		APIUrl: *args.APIUrl,
		Token:  *args.Token,
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
	}

	inicfg, err := ini.LoadFile(*args.INI)
	if err != nil {
		panic(err)
	}

	for project_path, ref := range inicfg["projects"] {
		work(&g, args, project_path, ref)
	}

}
