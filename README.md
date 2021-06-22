[![Build Status][badge-travis-image]][badge-travis-url]

# Kong Plugin Development Kit - Go edition

Docs: https://pkg.go.dev/github.com/Kong/go-pdk.

## Generators

Some code in this repo such as `server/kong_plugin_protocol/pluginsocket.pb.go` is generated
from `https://raw.githubusercontent.com/Kong/kong/master/kong/pluginsocket.proto`

After making a change to this file you can run the generators with:

```shell
./scripts/update-protoc-gen.sh
```

To check if the `server/kong_plugin_protocol/pluginsocket.pb.go` is up-to-date you can run:

```shell
./scripts/verify-protoc-gen.sh
```

[badge-travis-url]: https://travis-ci.com/Kong/go-pdk/branches
[badge-travis-image]: https://travis-ci.com/Kong/go-pdk.svg?branch=master
