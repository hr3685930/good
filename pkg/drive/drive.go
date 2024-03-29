package drive

import (
	"github.com/creasty/defaults"
	"github.com/fatih/structs"
	"github.com/jeremywohl/flatten"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"reflect"
	"strings"
)


// IgnoreErr ignore error
var IgnoreErr = false


// Conf Conf
type Conf interface {
	Connect(key string) error
	Register(key string)
}

// Drive Drive
func Drive(driveEnv interface{}) error {
	var typeInfo = reflect.TypeOf(driveEnv)
	var valInfo = reflect.ValueOf(driveEnv)
	num := typeInfo.NumField()
	var defaultDrive string
	for i := 0; i < num; i++ {
		params := make([]reflect.Value, 1)
		if typeInfo.Field(i).Name == "Default" {
			defaultDrive = valInfo.Field(i).String()
			continue
		}
		params[0] = reflect.ValueOf(strings.ToLower(typeInfo.Field(i).Name))
		item := valInfo.Field(i).MethodByName("Connect")
		res := item.Call(params)
		err := res[0].Interface()
		if err != nil {
			return err.(error)
		}
	}

	var defaultDriveIndex int
	for i := 0; i < num; i++ {
		if strings.ToLower(typeInfo.Field(i).Name) == defaultDrive {
			defaultDriveIndex = i
		}
	}

	dOption := make([]reflect.Value, 1)
	dOption[0] = reflect.ValueOf(strings.ToLower(defaultDrive))
	d := valInfo.Field(defaultDriveIndex).MethodByName("Register")
	d.Call(dOption)
	return nil
}

// Load priority env > yaml > default
func Load(e interface{}) error {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Transform config struct to map
	confMap := structs.Map(e)
	if err := mapstructure.Decode(e, &confMap); err != nil {
		return errors.Wrap(err, "Unable to Decode config")
	}

	// Flatten nested conf map
	flat, err := flatten.Flatten(confMap, "", flatten.DotStyle)
	if err != nil {
		return errors.Wrap(err, "Unable to flatten config")
	}

	// Bind each conf fields to environment vars
	for key := range flat {
		err := v.BindEnv(key)
		if err != nil {
			return errors.Wrapf(err, "Unable to bind env var: %s", key)
		}
	}

	_ = v.ReadInConfig()
	if err := defaults.Set(e); err != nil {
		return err
	}
	if err := v.Unmarshal(e); err != nil {
		return err
	}
	return nil
}
