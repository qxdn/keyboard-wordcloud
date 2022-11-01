package font

import "github.com/flopp/go-findfont"

/**
load font path if not exist load arial
*/
func LoadFontPath(name string) string {
	if path, err := findfont.Find(name); err == nil {
		return path
	}
	path, err := findfont.Find("arial.ttf")
	if err != nil {
		panic(err)
	}
	return path
}
