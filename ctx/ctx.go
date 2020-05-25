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

// Called by the plugin server at initialization.
func New(ch chan interface{}) Ctx {
	return Ctx{bridge.New(ch)}
}

// kong.Ctx.SetShared() sets a value in the `kong.ctx.shared` request context table.
func (c Ctx) SetShared(k string, value interface{}) error {
	_, err := c.Ask(`kong.ctx.shared.set`, k, value)
	return err
}

// kong.Ctx.GetSharedAny() returns a value from the `kong.ctx.shared` request context table.
func (c Ctx) GetSharedAny(k string) (interface{}, error) {
	return c.Ask(`kong.ctx.shared.get`, k)
}

// kong.Ctx.GetSharedString() returns a string value from the `kong.ctx.shared` request context table.
func (c Ctx) GetSharedString(k string) (string, error) {
	return c.AskString(`kong.ctx.shared.get`, k)
}

// kong.Ctx.GetSharedFloat() returns a float value from the `kong.ctx.shared` request context table.
func (c Ctx) GetSharedFloat(k string) (float64, error) {
	return c.AskFloat(`kong.ctx.shared.get`, k)
}

// kong.Ctx.GetSharedInt() returns an integer value from the `kong.ctx.shared` request context table.
func (c Ctx) GetSharedInt(k string) (int, error) {
	return c.AskInt(`kong.ctx.shared.get`, k)
}

