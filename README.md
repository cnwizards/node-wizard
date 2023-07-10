![Docker Build](https://github.com/cnwizards/node-wizard/actions/workflows/main.yml/badge.svg?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/cnwizards/node-wizard)](https://goreportcard.com/report/github.com/cnwizards/node-wizard)
# Node Wizard

### ⚠️ This project is still under development. Please use it at your own risk. ⚠️
This is a controller that monitors the `Ready` state of the nodes. If a node is not ready, it will be cordoned off and drained. As soon as the node is ready again, it is uncorded.

## Why do we need this?
Sometimes nodes go into `NotReady` state for some reason and the immediate response should be to evacuate the workloads on that node and cordon it off until it is ready again. This controller automates that process. Why would we want to automate this process? The first answer is faster response time: instead of waiting for a human to react, this controller will react instantly and cordon off and evacuate the node, thus rescheduling the workloads on that node to other nodes with less downtime. The second answer is that sometimes the node can recover itself after some time and be ready again. In this case, there may be nothing urgent and it can be investigated later. This controller will uncordon the node when it is ready.

## Features?
There are two main features of this controller:

* `Draining`: Non-graceful draining parameters can be set via an environment variable.
* `Uncordon`: The node will be uncordoned when it is ready.
* `Ignore Some Nodes`: Some nodes can be ignored by the controller by labeling with `node-wizard/ignore=true`.
* `Leader Election`: Application uses leader election mechanism. This is useful for high availability.

## Features to be added
* `Time to uncordon`: The time to uncordon the recovered node can be set via an environment variable.
* `Time to cordon`: The default node monitor grace period is 40 seconds. As this is quite a long time, the Node Wizard does not wait by default. However, it can be set via an environment variable.
* `Metrics`: Some metrics are exposed to Prometheus.

## How to install?
```
https://github.com/cnwizards/helm-charts/tree/main/charts/node-wizard
```
