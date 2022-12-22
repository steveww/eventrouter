package sinks

import (
	"errors"

	"github.com/golang/glog"
	"github.com/spf13/viper"
	v1 "k8s.io/api/core/v1"
)

// EventSinkInterface is the interface used to shunt events
type EventSinkInterface interface {
	UpdateEvents(eNew *v1.Event, eOld *v1.Event)
}

// ManufactureSink will manufacture a sink according to viper configs
// TODO: Determine if it should return an array of sinks
func ManufactureSink() (e EventSinkInterface) {
	s := viper.GetString("sink")
	glog.Infof("Sink is [%v]", s)
	switch s {
	case "glog":
		e = NewGlogSink()
	case "stdout":
		viper.SetDefault("stdoutJSONNamespace", "")
		stdoutNamespace := viper.GetString("stdoutJSONNamespace")
		e = NewStdoutSink(stdoutNamespace)
	case "null":
		e = NewNullSink()
	default:
		err := errors.New("Invalid Sink Specified")
		panic(err.Error())
	}
	return e
}
