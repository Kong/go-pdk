module github.com/Kong/gp-pdk/examples

go 1.18

require github.com/Kong/go-pdk v0.10.2

require (
	github.com/ugorji/go/codec v1.2.12 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)

replace github.com/Kong/go-pdk => ../
