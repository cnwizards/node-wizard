![Docker Build](https://github.com/cnwizards/node-wizard/actions/workflows/main.yml/badge.svg?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/cnwizards/node-wizard)](https://goreportcard.com/report/github.com/cnwizards/node-wizard)
# Node Wizard
Node Wizard is a controller that monitors node readiness. It cordons and drains nodes that are not ready, and uncordons them when they become ready again.

## Why do we need this?
The purpose of Node Wizard is to automate the response to nodes entering the "NotReady" state. It evacuates workloads from the affected node and cordons it off until it becomes ready again. Automating this process offers a faster response time compared to waiting for human intervention. The controller instantly reacts, cordons off the node, evacuates the workloads, and reschedules them on other nodes with minimal downtime.

Additionally, Node Wizard accounts for cases where the node may recover on its own over time. In such situations, there may not be an immediate urgency, allowing for investigation at a later time. When the node becomes ready again, the controller automatically uncordons it.

## Features
There are several features that Node Wizard offers:

* `Draining`: Non-graceful draining parameters can be set via an environment variable.
* `Uncordon`: The node will be uncordoned when it is ready.
* `Ignore Some Nodes`: Some nodes can be ignored by the controller by labeling with `node-wizard/ignore=true` (it can be useful for the ready nodes but some maintenance is going on).
* `Leader Election`: Application uses leader election mechanism. This is useful for high availability.
* `Metrics`: As now, two metrics are exposed:
&nbsp;

    | Metric Name | Metric Type | Description |
    | ----------- | ----------- | ----------- |
    | `node_wizard_uncordon_count` | Counter | Counter metric that shows the number of uncordon operations performed for each node. |
    | `node_wizard_drained_count` | Counter | Counter metric that shows the number of drain operations performed for each node. |

## Features to be added
* `Time to uncordon`: Time to uncordon feature is planned to be added in the future.
* `Time to cordon`: The default node monitor grace period is 40 seconds. As this is quite a long time, the Node Wizard does not wait by default. However, this feature can be added in the future.

## How to install?
```
# to add the Helm repository
helm repo add cnwizards https://charts.cloudnativewizards.dev 
# to install the Helm charts
helm install node-wizard cnwizards/node-wizard --namespace node-wizard --create-namespace
```

#### ⚠️ This project is still under development. ⚠️