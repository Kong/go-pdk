# :warning: This file is no longer maintained - use GH releases instead

## Table of Contents

- [v0.10.0](#v0100)
- [v0.9.0](#v090)
- [v0.8.0](#v080)
- [v0.7.1](#v071)
- [v0.6.0](#v060)
- [v0.5.0](#v050)
- [v0.4.0](#v040)
- [v0.3.1](#v031)

## [v0.10.0]

> Released 2023/09/07

### Fixes

- Use plugin sequence ID passed by Kong as instance ID; if one is not found,
pick a random value.
  [#158](https://github.com/Kong/go-pdk/pull/158)

## [v0.9.0]

> Released 2022/12/07

### Changes

- Fixed parameter type of `kong.service.request.set_raw_body`, return type of
  `kong.service.response.get_raw_body`,
  and body parameter type of `kong.response.exit` to `[]byte`.
  Note that old version (before 3.0.1, or commits before cd2bcf9) of kong is incompatible after this change.
  [#132](https://github.com/Kong/go-pdk/pull/132)
  [kong/kong#9526](https://github.com/Kong/kong/pull/9526)

## [v0.8.0]

> Released 2021/06/09

### Changes

- fix kong.Request.GetRawBody() with buffered content by @javierguerragiraldez in [#91](https://github.com/Kong/go-pdk/pull/91)
- avoid pass-by-value of objects that contain locks. by @javierguerragiraldez in [#79](https://github.com/Kong/go-pdk/pull/79)
- bump go version by @fffonion in [#112](https://github.com/Kong/go-pdk/pull/112)

### Addtions

- chore(*) add dependabot by @mayocream in [#98](https://github.com/Kong/go-pdk/pull/98)

## [v0.7.1]

> Released 2021/10/16

### Changes

- fix testing: don't break when the plugin Exit()s [#73](https://github.com/Kong/go-pdk/pull/73)
- Ignore unexported struct fields in config struct by @ctrox [#69](https://github.com/Kong/go-pdk/pull/69)
- Start every Headers field empty but non-null [#74](https://github.com/Kong/go-pdk/pull/74)

### Additions

- Feat/plugin testing [#64](https://github.com/Kong/go-pdk/pull/64)
- Add godoc comments [#65](https://github.com/Kong/go-pdk/pull/65)

## [v0.7.0]

> Released 2021/10/16

New ProtoBuf-based communications with Kong. Requires Kong 2.4.

## [v0.6.1]

> Released 2021/04/14

### Changes

- API bugfix: port values were given as string instead of Int

## [v0.6.0]

> Released 2021/04/14

### Additions

- New Embedded Server to replace the go-pluginserver. Requires Kong v2.3.

## [v0.5.0]

> Released 2020/05/27

### Changes

- Methods for `kong.ctx.shared` manipulation were moved from `kong.Node`
  to a new module named `kong.Ctx`, mirroring Kong's Lua PDK structure

### Additions

- Add `kong.Ctx` module, counterpart of Kong's Lua PDK [kong.ctx][kong.ctx] module.
  It contains the following methods:
  * `SetShared`: sets a value (of any type) into [`kong.ctx.shared`](https://docs.konghq.com/2.0.x/pdk/kong.ctx/#kongctxshared)
  * `GetSharedAny`: gets a value (of any type) from `kong.ctx.shared`
  * `GetSharedString`: gets a string from `kong.ctx.shared`
  * `GetSharedFloat`: gets a float from `kong.ctx.shared`
  * `GetSharedInt`: gets an integer from `kong.ctx.shared`
- Add new methods to `kong.Nginx`, allowing direct manipulation of the request context (`ngx.ctx`):
  * `SetCtx`: sets a value (of any type) into the request context
  * `GetCtxInt`: gets an integer value from the request context

## [v0.4.0]

> Released 2020/05/25

### Additions

- Add the `kong.service.response.get_raw_body` method, allowing Go plugins
  to access upstream Services response

## [v0.3.1]

> Released 2020/05/07

### Additions

- Add missing methods:
  * `kong.Nginx.GetSubsystem`
  * `kong.Node.SetCtxShared`
  * `kong.Node.GetCtxSharedAny`
  * `kong.Node.GetCtxSharedString`
  * `kong.Node.GetCtxSharedFloat`
  * `kong.Response.Exit`
  * `kong.Response.ExitStatus`

[Back to TOC](#table-of-contents)

[v0.10.0]: https://github.com/Kong/kong/compare/v0.9.0..v0.10.0
[v0.9.0]: https://github.com/Kong/kong/compare/v0.8.0..v0.9.0
[v0.8.0]: https://github.com/Kong/kong/compare/v0.7.0..v0.8.0
[v0.7.0]: https://github.com/Kong/kong/compare/v0.6.1..v0.7.0
[v0.6.0]: https://github.com/Kong/kong/compare/v0.5.0..v0.6.0
[v0.5.0]: https://github.com/Kong/kong/compare/v0.4.0..v0.5.0
[v0.4.0]: https://github.com/Kong/kong/compare/v0.3.1..v0.4.0
[v0.3.1]: https://github.com/Kong/kong/compare/v0.3.0..v0.3.1

[kong.ctx]: https://docs.konghq.com/2.0.x/pdk/kong.ctx/
