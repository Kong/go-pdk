# Table of Contents

- [v0.5.0](#v050)
- [v0.4.0](#v040)
- [v0.3.1](#v031)

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

[v0.5.0]: https://github.com/Kong/kong/compare/v0.4.0..v0.5.0
[v0.4.0]: https://github.com/Kong/kong/compare/v0.3.1..v0.4.0
[v0.3.1]: https://github.com/Kong/kong/compare/v0.3.0..v0.3.1

[kong.ctx]: https://docs.konghq.com/2.0.x/pdk/kong.ctx/
