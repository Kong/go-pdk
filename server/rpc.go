package server

import (
	"log"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/ugorji/go/codec"
)

type rpcHandler struct {
	constructor       func() interface{}
	configType        reflect.Type
	version           string // version number
	priority          int    // priority info
	lock              sync.RWMutex
	instances         map[int]*instanceData
	events            map[int]*eventData
	lastCloseInstance time.Time
}

var methodNames = [...]string{
	"Certificate",
	"Rewrite",
	"Access",
	"Response",
	"Preread",
	"Log",
}

func getHandlerNames(t reflect.Type) []string {
	handlers := []string{}
	for _, name := range methodNames {
		_, hasIt := t.MethodByName(name)
		if hasIt {
			handlers = append(handlers, strings.ToLower(name))
		}
	}
	return handlers
}

func newRpcHandler(constructor func() interface{}, version string, priority int) *rpcHandler {

	constructorType := reflect.TypeOf(constructor)
	if constructorType == nil {
		log.Printf("nil constructor")
		return nil
	}

	if constructorType.Kind() != reflect.Func {
		log.Printf("Constructor is not a function")
		return nil
	}

	if constructorType.NumIn() != 0 || constructorType.NumOut() != 1 {
		log.Printf("Wrong constructor signature")
		return nil
	}

	return &rpcHandler{
		constructor: constructor,
		configType:  reflect.TypeOf(constructor()),
		version:     version,
		priority:    priority,
		instances:   map[int]*instanceData{},
		events:      map[int]*eventData{},
	}
}

type schemaDict map[string]interface{}

func decodeSchema(s string) schemaDict {
	var (
		sd  schemaDict
		h   codec.JsonHandle
		ret strings.Builder
	)

	// allows for strings enclosed in single quotes
	for i := 0; i < len(s); i++ {
		if s[i] == '\'' && i > 0 && s[i-1] != '\\' {
			ret.WriteByte('"')
		} else {
			// allows for escaped single quotes as \'
			ret.WriteByte(s[i])
		}
	}

	enc := codec.NewDecoder(strings.NewReader(ret.String()), &h)
	enc.MustDecode(&sd)
	return sd
}

func updateSchemaDict(a, b schemaDict) {
	for k, v := range b {
		a[k] = v
	}
}

func getSchemaDict(t reflect.Type) schemaDict {
	switch t.Kind() {
	case reflect.String:
		return schemaDict{"type": "string"}

	case reflect.Bool:
		return schemaDict{"type": "boolean"}

	case reflect.Int, reflect.Int32:
		return schemaDict{"type": "integer"}

	case reflect.Uint, reflect.Uint32:
		return schemaDict{
			"type":    "integer",
			"between": []int{0, 2147483648},
		}

	case reflect.Float32, reflect.Float64:
		return schemaDict{"type": "number"}

	case reflect.Ptr:
		return getSchemaDict(t.Elem())

	case reflect.Slice:
		elemType := getSchemaDict(t.Elem())
		if elemType == nil {
			break
		}
		return schemaDict{
			"type":     "array",
			"elements": elemType,
		}

	case reflect.Map:
		kType := getSchemaDict(t.Key())
		vType := getSchemaDict(t.Elem())
		if kType == nil || vType == nil {
			break
		}
		return schemaDict{
			"type":   "map",
			"keys":   kType,
			"values": vType,
		}

	case reflect.Struct:
		fieldsArray := []schemaDict{}
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			// ignore unexported fields
			if len(field.PkgPath) != 0 {
				continue
			}
			typeDecl := getSchemaDict(field.Type)
			if typeDecl == nil {
				// ignore unrepresentable types
				continue
			}
			name := field.Tag.Get("json")
			if name == "" {
				name = strings.ToLower(field.Name)
			}

			if s, ok := field.Tag.Lookup("schema"); ok {
				if customSchema := decodeSchema(s); customSchema != nil {
					if field.Type.Kind() == reflect.Map {
						updateSchemaDict(typeDecl["keys"].(schemaDict), customSchema)
					} else {
						updateSchemaDict(typeDecl, customSchema)
					}
				}
			}

			fieldsArray = append(fieldsArray, schemaDict{name: typeDecl})
		}
		return schemaDict{
			"type":   "record",
			"fields": fieldsArray,
		}
	}

	return nil
}

type pluginInfo struct {
	Name     string     // plugin name
	ModTime  time.Time  `codec:",omitempty"` // plugin file modification time
	LoadTime time.Time  `codec:",omitempty"` // plugin load time
	Phases   []string   // events it can handle
	Version  string     // version number
	Priority int        // priority info
	Schema   schemaDict // representation of the config schema
}

func (rh *rpcHandler) getInfo() (info pluginInfo, err error) {
	name, err := getName()
	if err != nil {
		return
	}

	info = pluginInfo{
		Name:   name,
		Phases: getHandlerNames(rh.configType),
		Schema: schemaDict{
			"name": name,
			"fields": []schemaDict{
				{"config": getSchemaDict(rh.configType)},
			},
		},
		Version:  rh.version,
		Priority: rh.priority,
	}

	return
}
