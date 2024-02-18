import React from 'react';

export function Events({ messages, username }){
  return (
    <div className="flex flex-col space-y-2 p-3 overflow-y-auto max-h-96 rounded-lg">
    {messages ? messages.map((message, index) => 
      <div key={index}>
        <span className="text-green-500">{username}</span>: { message }
      </div>
      ): <div>No messages</div>
    }
    </div>
  );
}
