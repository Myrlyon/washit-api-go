package utils

type Link struct {
	Href   string `json:"href"`
	Method string `json:"method"`
}

// func CreateLinks(paths map[string]string) map[string]string {

// }

// links := make(map[string]Link, len(paths))
// 	for rel, path := range paths {
// 		if href, ok := path.(string); ok {
// 			links[rel] = Link{
// 				Href:   href,
// 				Method: "GET",
// 			}
// 		}
// 	}
// 	return links
