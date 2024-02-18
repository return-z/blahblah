import React, { useState, useEffect } from 'react';
import { Chat } from './components/Chat';
import { NameModal } from './components/NameModal';
import { Chatrooms } from './components/Chatrooms';
import axios from 'axios';

const API = axios.create({ baseURL: 'http://localhost:5990/' });

export default function App(){
  const [username, setUserName] = useState('')
  const [chatrooms, setChatrooms] = useState([])
  const auth = (data) => API.post('/auth', data);

  const handleSaveName = async (name) => {
    const authData = { 'username': name }
    console.log(name);
    try {
      const { data } = await auth(authData);
      setChatrooms(data.map(chatroom => chatroom.Name))
      setUserName(name);
    }
    catch(error){
      alert("User doesn't exist");
      console.log(error);
    }
  }

  return (
  <div class="flex flex-col h-screen">
    { !username ? <NameModal onSaveName={handleSaveName} /> 
        : <Chatrooms username={username} chatrooms={chatrooms}/> }
  </div>
  );
}
