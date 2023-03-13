# Eventrouter

Receives Kubernetes events and turns them into Prometheus metrics.

## Important Events

Good background information on watching Kubernetes [events](https://isitobservable.io/observability/kubernetes/how-to-collect-kubernetes-events)

* **CrashLoopBackOff** - which happens when a pod starts, crashes, starts again, and then crashes again
* **ImagePullBackOff** - which happens when the node is unable to retrieve the image
* **Evicted** - which can happen when a node determines that a pod needs to be evicted or terminated to free up some resources (CPU, memory...etc). When this happens, K8s is supposed to reschedule the pod on another node
* **FailedMount** / FailedAttachVolume - when pods require a persistent volume or storage, this event prevents them from starting if the storage is not accessible
* **FailedSchedulingEvents** - when the scheduler is not able to find a node to run your pods
* **NodeNotReady** - when a node cannot be used to run a pod because of an underlying issue
* **Rebooted**
* **HostPort Conflict**

# TODO

* Investigate possible pruning options for `message` to reduce cardinality.
* Delete event remotes Prometheus metric label set - tricky