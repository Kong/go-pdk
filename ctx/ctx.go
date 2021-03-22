/*
Current request context data.

`kong.ctx.shared`:
A table that has the lifetime of the current request and is shared between
all plugins. It can be used to share data between several plugins in a given
request.

Values inserted in this table by a plugin will be visible by all other
plugins.  One must use caution when interacting with its values, as a naming
conflict could result in the overwrite of data.

Usage:
// Two plugins A and B, if plugin A has a higher priority than B's
// (it executes before B), and plugin A is a Go plugin:

// plugin A PluginA.go
func (conf Config) Access(kong *pdk.PDK) {
	err := kong.Ctx.SetShared("hello world")
	if err != nil {
		kong.Log.Err(err)
		return
	}
}

// plugin B handler.lua
function plugin_b_handler:access(conf)
  kong.log(kong.ctx.shared.foo) // "hello world"
end
*/
package ctx

import (
	"github.com/Kong/go-pdk/bridge"
)

// Holds this module's functions.  Accessible as `kong.Ctx`
type Ctx struct {
	bridge.PdkBridge
}

// kong.Ctx.SetShared() sets a value in the `kong.ctx.shared` request context table.
func (c Ctx) SetShared(k string, value interface{}) error {
	err := c.Ask(`kong.ctx.shared.set`, bridge.WrapString(k), nil)
	return err
}

// kong.Ctx.GetSharedAny() returns a value from the `kong.ctx.shared` request context table.
func (c Ctx) GetSharedAny(k string) (interface{}, error) {
	return c.AskValue(`kong.ctx.shared.get`, bridge.WrapString(k))
}

// kong.Ctx.GetSharedString() returns a string value from the `kong.ctx.shared` request context table.
func (c Ctx) GetSharedString(k string) (string, error) {
	v, err := c.GetSharedAny(k)
	if err != nil {
		return "", err
	}

	s, ok := v.(string)
	if ok {
		return s, nil
	}

	return "", bridge.ReturnTypeError("string")
}

// kong.Ctx.GetSharedFloat() returns a float value from the `kong.ctx.shared` request context table.
func (c Ctx) GetSharedFloat(k string) (float64, error) {
	v, err := c.GetSharedAny(k)
	if err != nil {
		return 0, err
	}

	f, ok := v.(float64)
	if ok {
		return f, nil
	}

	return 0, bridge.ReturnTypeError("number")
}

