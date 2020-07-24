# Server Sent Events using go

From blog [Server-Sent Events with Go and React](https://medium.com/wesionary-team/server-sent-events-with-go-and-react-76df101a3efe).

Forked from [dipeshdulal/sse-go](https://github.com/dipeshdulal/sse-go).

To run server: (will run in `5000`  port)

1. `go run main.go`

To run frontend: (will start frontend server in `3000` port)

1. `yarn install`
1. `yarn start`

After running the frontend. You can send request to server or start or stop streaming by clicking in the button show below in screenshot.

![Screenshot](sc.png)

## What it does

Client sends the request to `/log` endpoint and every subscriber that is listening to the stream channel will receive message from server and frontend just displays the messages received from stream.

You can directly check streams going to `localhost:5000/sse` which will never finish responding to the events until client closes it. After opening the `/see` you can open another tab to `/log` and see that tab with `/sse` has output. Or, if you use curl you can just do.

```bash
curl -X GET localhost:5000/log
curl -X POST -d "MY DATA" localhost:5000/log
```

And check how every request that is made to `log` endpoint is streamed to `sse` using ServerSideEvents.

## To Initialize SSE from JavaScript

```js
var sse = new EventSource("http://localhost:5000/sse")
sse.onmessage = console.log
```

## And to close the EventSource.

```js
sse.close()
```

To learn more about SSE check the [mozilla documentation](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events)
