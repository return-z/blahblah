
import React, { useState } from 'react';


export function NameModal({ onSaveName }) {
  const [name, setName] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();
    onSaveName(name);
  };

  return (
    <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full flex items-center justify-center">
      <div className="bg-white p-5 rounded-md flex flex-col items-center space-y-4">
        <h2 className="text-lg font-semibold">Enter Your Name</h2>
        <form onSubmit={handleSubmit} className="flex flex-col items-center space-y-2">
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="border p-2 rounded"
            placeholder="Your name"
            required
          />
          <button type="submit" className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-700">Enter Chat</button>
        </form>
      </div>
    </div>
  );
}

