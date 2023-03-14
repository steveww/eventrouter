package main

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"github.com/steveww/eventrouter/sinks"

	v1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

var (
	kubernetesWarningEventCounterVec = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "eventrouter_warnings_total",
		Help: "Total number of warning events in the kubernetes cluster",
	}, []string{
		"involved_object_kind",
		"involved_object_name",
		"involved_object_namespace",
		"reason",
		"source",
		"message",
	})
)

// EventRouter is responsible for maintaining a stream of kubernetes
// system Events and pushing them to another channel for storage
type EventRouter struct {
	// kubeclient is the main kubernetes interface
	kubeClient kubernetes.Interface

	// store of events populated by the shared informer
	eLister corelisters.EventLister

	// returns true if the event store has been synced
	eListerSynched cache.InformerSynced

	// event sink
	// TODO: Determine if we want to support multiple sinks.
	eSink sinks.EventSinkInterface
}

// NewEventRouter will create a new event router using the input params
func NewEventRouter(kubeClient kubernetes.Interface, eventsInformer coreinformers.EventInformer) *EventRouter {
	if viper.GetBool("enable-prometheus") {
		prometheus.MustRegister(kubernetesWarningEventCounterVec)

		g := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "version_info",
			Help: "Version information",
			ConstLabels: prometheus.Labels{
				"version": Version,
			},
		})
		prometheus.MustRegister(g)
		g.Set(1)
	}

	er := &EventRouter{
		kubeClient: kubeClient,
		eSink:      sinks.ManufactureSink(),
	}
	eventsInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    er.addEvent,
		UpdateFunc: er.updateEvent,
		DeleteFunc: er.deleteEvent,
	})
	er.eLister = eventsInformer.Lister()
	er.eListerSynched = eventsInformer.Informer().HasSynced
	return er
}

// Run starts the EventRouter/Controller.
func (er *EventRouter) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer glog.Infof("Shutting down EventRouter")

	glog.Infof("Starting EventRouter")

	// here is where we kick the caches into gear
	if !cache.WaitForCacheSync(stopCh, er.eListerSynched) {
		utilruntime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}
	<-stopCh
}

// addEvent is called when an event is created, or during the initial list
func (er *EventRouter) addEvent(obj interface{}) {
	e := obj.(*v1.Event)
	if e.Type == "Warning" {
		prometheusEvent(e)
		er.eSink.UpdateEvents(e, nil)
	}
}

// updateEvent is called any time there is an update to an existing event
func (er *EventRouter) updateEvent(objOld interface{}, objNew interface{}) {
	eOld := objOld.(*v1.Event)
	eNew := objNew.(*v1.Event)
	if eNew.Type == "Warning" {
		prometheusEvent(eNew)
		er.eSink.UpdateEvents(eNew, eOld)
	}
}

// deleteEvent should only occur when the system garbage collects events via TTL expiration
// NOTE: This should *only* happen on TTL expiration there
func (er *EventRouter) deleteEvent(obj interface{}) {
	e := obj.(*v1.Event)
	if e.Type == "Warning" {
		unregisterEvent(e)
	}
}

// prometheusEvent is called when an event is added or updated
func prometheusEvent(event *v1.Event) {
	if !viper.GetBool("enable-prometheus") {
		return
	}

	// limit the length of messages
	message := substr(event.Message, 0, 50)
	counter, err := kubernetesWarningEventCounterVec.GetMetricWithLabelValues(
		event.InvolvedObject.Kind,
		event.InvolvedObject.Name,
		event.InvolvedObject.Namespace,
		event.Reason,
		event.Source.Host,
		message,
	)

	if err != nil {
		// Not sure this is the right place to log this error?
		glog.Warning(err)
	} else {
		counter.Add(1)
	}
}

func unregisterEvent(event *v1.Event) {
	if !viper.GetBool("enable-prometheus") {
		return
	}

	message := substr(event.Message, 0, 50)
	counter, err := kubernetesWarningEventCounterVec.GetMetricWithLabelValues(
		event.InvolvedObject.Kind,
		event.InvolvedObject.Name,
		event.InvolvedObject.Namespace,
		event.Reason,
		event.Source.Host,
		message,
	)
	if err != nil {
		glog.Warning(err)
	} else {
		if ok := prometheus.Unregister(counter); !ok {
			glog.Warningf("unregister not OK %v", message)
		}
	}
}

func substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}
