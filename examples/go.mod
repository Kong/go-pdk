module github.com/Kong/gp-pdk/examples

go 1.18

require github.com/Kong/go-pdk v0.8.0

replace github.com/Kong/go-pdk => ../

require (
	github.com/ugorji/go/codec v1.2.12 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
)
