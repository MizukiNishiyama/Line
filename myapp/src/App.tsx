import React, { useState, useEffect, ChangeEvent, FC } from 'react';
import { BrowserRouter as Router, Route, Switch, useHistory, useRouteMatch, useLocation, Link} from 'react-router-dom';
import './App.css';

type User = {
  UserId: string;
  UserName: string;
  UserPassword: string;
};

type Message = {
  MessageId: string;
  MessageContent: string;
  MessageTime: string;
  UserId: string;
  RoomId: string;
  UserName:string;
};

type Room = {
  RoomId: string;
  UserId1: string;
  UserId2: string;
  UserName1: string;
  UserName2: string;
};

const HomePage: FC = () => {
  return (
    <div>
      <Link to="/login">
        <button>Login</button>
      </Link>
      <Link to="/signup">
        <button>Signup</button>
      </Link>
    </div>
  );
};

const Login: FC = () => {
  const [username, setUsername] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [error, setError] = useState<string | null>(null);
  const history = useHistory();

  const login = async () => {
    try {
      const res = await fetch('http://localhost:8000/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ UserName: username, UserPassword: password }),
      });

      if (res.ok) {
        const data = await res.json();
        const userId = data.UserId;
        localStorage.setItem('userId', userId);
        localStorage.setItem('userName', username);
        history.push(`/rooms?userid=${userId}`);
        console.log(data)
      } else {
        throw new Error('Failed to login');
      }
    } catch (err) {
      setError((err as Error).message);
    }
  };

  return (
    <div>
      <h2>Login</h2>
      <input value={username} onChange={(e: ChangeEvent<HTMLInputElement>) => setUsername(e.target.value)} placeholder="Username" />
      <input value={password} onChange={(e: ChangeEvent<HTMLInputElement>) => setPassword(e.target.value)} placeholder="Password" type="password" />
      <button onClick={login}>Log in</button>
      <button onClick={() => history.push('/signup')}>Sign up</button>
      {error && <p>Error: {error}</p>}
    </div>
  );
};

const Signup: FC = () => {
  const [username, setUsername] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [error, setError] = useState<string | null>(null);
  const history = useHistory();

  const signup = async () => {
    try {
      const res = await fetch('http://localhost:8000/signup', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ UserName: username, UserPassword: password }),
      });

      if (res.ok) {
        history.push('/');
      } else {
        throw new Error('Failed to sign up');
      }
    } catch (err) {
      setError((err as Error).message);
    }
  };

  return (
    <div>
      <h2>Sign Up</h2>
      <input value={username} onChange={(e: ChangeEvent<HTMLInputElement>) => setUsername(e.target.value)} placeholder="Username" />
      <input value={password} onChange={(e: ChangeEvent<HTMLInputElement>) => setPassword(e.target.value)} placeholder="Password" type="password" />
      <button onClick={signup}>Sign up</button>
      {error && <p>Error: {error}</p>}
    </div>
  );
};

const RoomList: FC = () => {
  const [rooms, setRooms] = useState<Room[]>([]);
  const [error, setError] = useState<string | null>(null);
  const history = useHistory();
  const location = useLocation();
  const queryParams = new URLSearchParams(location.search);
  const userId = queryParams.get('userid');

  useEffect(() => {
    const fetchRooms = async () => {
      try {
        const res = await fetch(`http://localhost:8000/rooms?userid=${userId}`);
        const data = await res.json();
        console.log(data)

        if (res.ok) {
          setRooms(data);
        } else {
          throw new Error('Failed to fetch rooms');
        }
      } catch (err) {
        setRooms([]);
        setError((err as Error).message);
      }
    };

    if (userId) {
      fetchRooms();
    }
  }, [userId]);


  return (
    <div>
      <h2>Rooms</h2>
      {rooms.map(room => (
        <div key={room.RoomId} onClick={() => history.push(`/chat/${room.RoomId}`)} className="room-list-item">
          {room.UserName1} - {room.UserName2}
        </div>
      ))}
      {error && <p>Error: {error}</p>}
    </div>
  );
};

let socket: WebSocket | null = null;

const ChatRoom: FC = () => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [newMessage, setNewMessage] = useState<string>('');
  const [error, setError] = useState<string | null>(null);
  const match = useRouteMatch<{ roomId: string }>('/chat/:roomId');


  useEffect(() => {
    const fetchMessages = async () => {
      try {
        const res = await fetch(`http://localhost:8000/fetchmessage?roomid=${match?.params.roomId}`);
        const data = await res.json();
        console.log(data);

        if (res.ok) {
          setMessages(data);
          console.log("fetchmessage success" )
        } else {
          throw new Error('Failed to fetch messages');
        }
      } catch (err) {
        setError((err as Error).message);
      }
    };

    if (match?.params.roomId) {
      fetchMessages();
    }

    try {
      socket = new WebSocket('ws://localhost:8000/ws');
    } catch (error) {
      console.error("Failed to create WebSocket:", error);
      return;
    }
  
    if(socket){
      socket.onopen = () => {
        // WebSocket has connected
        if (socket) {
          socket.onmessage = (event) => {
            const message = JSON.parse(event.data);
            setMessages((prevMessages) => [...prevMessages, message]);
          };
        }
      };
    }
  
    return () => {
      if (socket && socket.readyState === WebSocket.OPEN) {
        // Only close the connection if it is open
        socket.close();
      }
      socket = null;
    };

  }, [match?.params.roomId]);

  // const sendMessage = async () => {
  //   try {
  //     const userId = localStorage.getItem('userId');
  //     const userName = localStorage.getItem('userName');
  //     console.log(userId)
  //     const res = await fetch('http://localhost:8000/sendmessage', {
  //       method: 'POST',
  //       headers: {
  //         'Content-Type': 'application/json',
  //       },
  //       body: JSON.stringify({ MessageContent: newMessage, MessageTime: ` (${new Date().getMonth()+1}/${new Date().getDate()} ${new Date().getHours()}:${new Date().getMinutes()})`, UserId: userId, RoomId: match?.params.roomId, UserName: userName }),
  //     });
  //     const data = await res.json();

  //     if (res.ok) {
  //       setMessages([...messages, data.message]);
  //       setNewMessage('');
  //       console.log(data.message)
  //     } else {
  //       throw new Error('Failed to send message');
  //     }
  //   } catch (err) {
  //     setError((err as Error).message);
  //   }
  const sendMessage = async () => {
    try {
      const userId = localStorage.getItem('userId');
      const userName = localStorage.getItem('userName');

      if (socket && socket.readyState === WebSocket.OPEN && userId && userName) {
        socket.send(JSON.stringify({
          MessageContent: newMessage, 
          MessageTime: ` (${new Date().getMonth()+1}/${new Date().getDate()} ${new Date().getHours()}:${new Date().getMinutes()})`, 
          UserId: userId, 
          RoomId: match?.params.roomId, 
          UserName: userName 
        }));
        setNewMessage('');
      } else {
        console.log("WebSocket is not open. Ready state:", socket?.readyState);
      }
    } catch (err) {
      setError((err as Error).message);
    }
  };
  

  return (
    <div>
      <h2>Chat Room</h2>
      {messages && messages.map(message => (
    <div key={message.MessageTime} className="message-list-item">
      {message.MessageContent} {message.UserName} {message.MessageTime}
    </div>
    ))}

      <div className="message-input-container">
        <input value={newMessage} onChange={(e: ChangeEvent<HTMLInputElement>) => setNewMessage(e.target.value)} placeholder="New message" />
        <button onClick={sendMessage}>Send</button>
      </div>
      {error && <p>Error: {error}</p>}
    </div>
  );
};

const App: FC = () => {
  return (
    <Router>
      <Switch>
        <Route path="/signup">
          <Signup />
        </Route>
        <Route path="/login">
          <Login />
        </Route>
        <Route path="/rooms">
          <RoomList />
        </Route>
        <Route path="/chat/:roomId">
          <ChatRoom />
        </Route>
        <Route path="/">
          <HomePage />
        </Route>
      </Switch>
    </Router>
  );
};

export default App;
