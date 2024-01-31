package runner

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"

	"module-wget/parser"
)

// Sitemap represents data needed to build sitemap using BuildSitemap
type Sitemap struct {
	rootLink     string
	visitedLinks map[string]struct{}
	directory    string
}

// NewSitemap creates instance of Sitemap
func NewSitemap(rootLink string, directory string) Sitemap {
	v := map[string]struct{}{
		rootLink: {},
	}

	s := Sitemap{
		rootLink:     rootLink,
		visitedLinks: v,
		directory:    directory,
	}

	return s
}

// DownloadSite recursively visits links in queue and downloads them.
// If depth is greater than zero it restricts number of recursive calls.
func (s *Sitemap) DownloadSite(queue []string, depth int) error {
	if depth == 0 {
		return nil
	}

	discoveredLinks, err := s.processQueue(queue)
	if err != nil {
		return err
	}

	if len(discoveredLinks) > 0 {
		depth--
		return s.DownloadSite(discoveredLinks, depth)
	}

	return nil
}

func (s *Sitemap) processQueue(queue []string) ([]string, error) {
	discoveredLinks := make([]string, 0)

	for _, v := range queue {
		links, err := s.processLink(v)
		if err != nil {
			return nil, err
		}
		discoveredLinks = append(discoveredLinks, links...)
	}

	return discoveredLinks, nil
}

func (s *Sitemap) processLink(link string) ([]string, error) {
	resp, err := s.getResponse(link)
	if err != nil {
		return nil, err
	}

	body, err := s.getBody(resp)
	if err != nil {
		return nil, err
	}

	fileName, err := s.createFile(resp, body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("\nDownloaded a file %s\n", fileName)

	return s.parseLinks(body)
}

func (s *Sitemap) getResponse(link string) (*http.Response, error) {
	resp, err := http.Get(link)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}
	return resp, nil
}

func (s *Sitemap) getBody(resp *http.Response) (*bytes.Reader, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	return bytes.NewReader(body), nil
}

func (s *Sitemap) createFile(resp *http.Response, body *bytes.Reader) (string, error) {
	fileName, err := s.createFileFromResponse(resp, body)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	return fileName, nil
}

func (s *Sitemap) createFileFromResponse(resp *http.Response, r *bytes.Reader) (string, error) {
	mediatype, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return "", err
	}
	ext, err := mime.ExtensionsByType(mediatype)
	if err != nil || len(ext) == 0 {
		return "", err
	}

	fileName := path.Join(s.directory, path.Base(resp.Request.URL.Path)+ext[0])
	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

// parseLinks reads html data from io.Reader and creates array of links.
// It only parses links with Sitemap.rootLink domain
func (s *Sitemap) parseLinks(r io.Reader) ([]string, error) {
	res, err := parser.ParseHTML(r)
	if err != nil {
		return nil, err
	}

	links := make([]string, 0)

	for _, v := range res {
		href := v.Href.Host + v.Href.Path
		visited := true

		switch {
		case strings.HasPrefix(href, s.rootLink):
			visited = s.isVisited(href)
		case strings.HasPrefix(href, "/"):
			href = s.rootLink + href
			visited = s.isVisited(href)
		}

		if !visited {
			links = append(links, href)
		}
	}

	return links, nil
}

// isVisited checks if url visited
func (s *Sitemap) isVisited(href string) (visited bool) {
	if _, visited = s.visitedLinks[href]; !visited {
		s.visitedLinks[href] = struct{}{}
	}

	return visited
}
