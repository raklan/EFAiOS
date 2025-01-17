function sendWsMessage(ws, msgType, data){
    console.info("sending message", {jsonType: msgType, data: data})
    if(!ws || !ws?.OPEN){
        console.assert(ws, 'WebSocket has not been initialized')
        console.assert(ws?.OPEN, 'WebSocket connection is not open')
        return;
    }
    ws.send(JSON.stringify({
        jsonType: msgType,
        data: data
    }))
}