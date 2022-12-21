# Eventrouter

Receives Kubernetes events and turns them into Prometheus metrics.

## Important Events

* CrashLoopBackOff
* ImagePullBackOff
* Evicted
* FailedMount / FailedAttachVolume
* FailedSchedulingEvents
* NodeNotReady
* Rebooted
* HostPort Conflict
