import React, { useState } from 'react';

export function Form({ websocket, chatroom, username }){
  const [value, setValue] = useState('');

  function onSubmit(e){
    e.preventDefault();
    if (websocket && websocket.readyState === WebSocket.OPEN){
      websocket.send(value);
    }
    else{
      console.error("Websocket is not open")
    }
    setValue('');
  }

  return (
    <form onSubmit={ onSubmit }>
    <label for="search" class="mb-2 text-sm font-medium text-gray-900 sr-only dark:text-white">Search</label>
      <div class="relative">
        <input value={value} onChange={(e) => setValue(e.target.value)}id="message" class="block w-full p-4 ps-10 text-sm text-gray-900 border border-gray-300 rounded-lg bg-gray-50 focus:ring-blue-500 focus:border-blue-500" placeholder="Message" required />
        <button type="submit" class="text-white absolute end-2.5 bottom-2.5 bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-4 py-2 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">
          Send</button>
      </div>
    </form>
  );
}

