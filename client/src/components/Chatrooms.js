import React, { useState } from 'react';
import { Chat } from './Chat'

export function Chatrooms({ username, chatrooms }){
  const [chatroom, setChatroom] = useState('')
  const handleClick = (name) => {
    setChatroom(name);
  }
  return (
    <div className="chatrooms">
    { !chatroom ? (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">Subscribed Chatrooms</h1>
      <div className="grid grid-cols-3 gap-4">
        {chatrooms.map((name, index) => (
          <button
            key={index}
            className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
            onClick={() => handleClick(name)}
          >
            {name}
          </button>
        ))}
      </div>
    </div>
    ) : <Chat username={username} chatroom={chatroom} />
    }
    </div>
  );
}
