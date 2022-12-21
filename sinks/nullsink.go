package sinks

import v1 "k8s.io/api/core/v1"

type NullSink struct {
}

func NewNullSink() EventSinkInterface {
	return &NullSink{}
}

func (ns *NullSink) UpdateEvents(eNew *v1.Event, eOld *v1.Event) {
	// nothing happens here
}
