# Azure SignalR Broadcaster

This small demo illustrates how you can use [*Azure SignalR*](https://docs.microsoft.com/en-us/azure/azure-signalr/) from Go.

At the time of writing, there is no official SignalR server implementation. If you have a webserver written in Go that would like to push messages to e.g. browser clients via SignalR, you can achieve that by using Azure's SignalR service. It offers a REST API with which you can send messages without having to implement the SignalR protocol. Azure is running the SignalR server for you. All you have to do is to implement the *negotiate* endpoint that is described in the [Azure SignalR documentation](https://github.com/Azure/azure-signalr/blob/dev/docs/rest-api.md).

If you want to try the sample,

* build the code using Go,
* run *azure-signalr-server -h* to see the command line args,
* provide meaningful arguments,
* and open *http://localhost:8081/* (or whatever port you specified in the command line) in a browser.
