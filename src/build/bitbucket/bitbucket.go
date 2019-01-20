package bitbucket

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/shibukawa/configdir"
)

var cachedRelease map[string]string
var cache = configdir.New("tanema", "ludus").QueryCacheFolder()

// Get will fetch the release for that arch
func Get(version, arch string) (zipReader *zip.Reader, err error) {
	if file, err := cache.Open(arch + "-" + version + ".zip"); err == nil {
		info, err := file.Stat()
		if err != nil {
			return nil, err
		}
		zipReader, err = zip.NewReader(file, int64(info.Size()))
	} else {
		zipData, err := fetchZipFileFromBitbucket(version, arch)
		if err != nil {
			return nil, err
		}
		zipReader, err = zip.NewReader(zipData, int64(zipData.Len()))
	}
	return zipReader, err
}

func fetchZipFileFromBitbucket(version, arch string) (*bytes.Reader, error) {
	archMap, err := getRelease(version)
	if err != nil {
		return nil, err
	}

	zipURL, ok := archMap[arch]
	if !ok {
		return nil, fmt.Errorf("Arch %v not found in the love realease %v", arch, version)
	}

	resp, err := http.Get(zipURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	cache.WriteFile(arch+"-"+version+".zip", body)

	return bytes.NewReader(body), nil
}

func getRelease(version string) (map[string]string, error) {
	if len(cachedRelease) > 0 {
		return cachedRelease, nil
	}
	archMap := map[string]string{}
	downloads := []download{}
	currentPage := downloadsPage{Next: downloadsURL}
	var err error
	for currentPage.Next != "" {
		currentPage, err = getDownloadPage(currentPage.Next)
		if err != nil {
			return archMap, err
		}
		downloads = append(downloads, currentPage.Downloads...)
		if selected := selectDownloadVersions(version, downloads); len(selected) >= 3 {
			for _, download := range selected {
				if strings.Contains(download.Name, "macos") {
					archMap["macos"] = download.Links.Self.Href
				} else if strings.Contains(download.Name, "win32") {
					archMap["win32"] = download.Links.Self.Href
				} else if strings.Contains(download.Name, "win64") {
					archMap["win64"] = download.Links.Self.Href
				}
			}
			cachedRelease = archMap
			return archMap, nil
		}
	}
	return archMap, fmt.Errorf("Cannot find an appropriate release for version %v", version)
}

func selectDownloadVersions(version string, downloads []download) []download {
	selected := []download{}
	for _, download := range downloads {
		if strings.HasPrefix(download.Name, "love-"+version) &&
			!strings.Contains(download.Name, "src") && !strings.Contains(download.Name, "source") &&
			!strings.Contains(download.Name, "libraries") && strings.HasSuffix(download.Name, ".zip") {
			selected = append(selected, download)
		}
	}
	return selected
}
