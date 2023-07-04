
## To Do's
Here is a list of things that need to be done before this project can be considered complete.

* [ ] Code review
  * [ ] Optimizing controller code
  * [ ] Optimizing handler logic
----

* [ ] Feature
  * [ x ] Graceful draining options should be tunable via a configmap
  * [ x ] Non-graceful draining options can be enabled and should be tunable via a configmap
  * [ ] Time to uncordon the recovered node should be tunable via a configmap
  * [ ] Time to cordon the node should be tunable via a configmap
  * [ ] Some counter metrics should be added like how many pods evicted, how many nodes recovered, etc.

----
* [ ] Documentation
  * [ ] Fix Readme
  * [ ] Helm Chart