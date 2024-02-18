import React, { useState, useEffect } from 'react';
import { Events } from './Events'
import { Form } from './Form' 

export function Chat({ username, chatroom }){
  const [websocket, setWebsocket] = useState(null); 
  const [messages, setMessages] = useState([]); 

  useEffect(() => {
    var ws = new WebSocket(`ws://localhost:5990/ws/${chatroom}`)
    console.log('trying to connect to ws')
    ws.onopen = () => {
      console.log("Websocket is connected")
    };
    
    ws.onmessage = (e) => {
      const message = e.data;
      console.log("Message from server: ", message);
      console.log(messages);
      setMessages((prevMessages) => [...prevMessages, message]);
    };

    ws.onerror = (e) => {
      console.log("WS error: ", e);
    };

    ws.onclose = (e) => {
      console.log("websocket is closed");
    };

    setWebsocket(ws);

    return () => {
      ws.close();
    };
  }, []);

  return (
  <div class="flex flex-col h-screen">
    <h3 className="text-center bg-gray-100 text-2xl font-bold">Chat App</h3>
    <div class="flex-grow overflow-auto">
      <Events messages={messages} username={username} />
    </div>
    <div class="p-4 bg-gray-100 w-full fixed bottom-0 left-0">
      <Form websocket={websocket} />
    </div>
  </div>
  );
}
