package utils

type Link struct {
	Href string `json:"href"`
}

func CreateLinks(paths map[string]string) map[string]Link {
	links := make(map[string]Link, len(paths))
	for rel, path := range paths {
		links[rel] = Link{Href: path}
	}
	return links
}
