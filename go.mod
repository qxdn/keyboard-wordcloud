module github.com/qxdn/keyboard-wordcloud

go 1.18

require (
	github.com/flopp/go-findfont v0.1.0
	github.com/psykhi/wordclouds v0.0.0-20220728072901-2d77dabdd4fd
	github.com/sirupsen/logrus v1.9.0
)

require (
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	golang.org/x/image v0.0.0-20191009234506-e7c1f5e7dbb8 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)

replace github.com/psykhi/wordclouds => github.com/qxdn/wordclouds v0.0.0-20221102135804-f8d3dd8234fd
