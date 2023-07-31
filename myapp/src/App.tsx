import React, { useState, useEffect, ChangeEvent, FC } from 'react';
import { BrowserRouter as Router, Route, Switch, useHistory, useRouteMatch, useLocation, Link} from 'react-router-dom';
import './App.css';


import LogoutIcon from '@mui/icons-material/Logout';
import PersonAddAltIcon from '@mui/icons-material/PersonAddAlt';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import SendIcon from '@mui/icons-material/Send';

import { Button } from '@mui/material';
import { createTheme, ThemeProvider } from '@mui/material/styles';



const theme = createTheme({
  palette: {
    primary: {
      main: '#fff', 
    },
    secondary: {
      main: '#000000',
    },
    // 追加の色を指定
    error: {
      main: '#f44336',
    },
    warning: {
      main: '#ff9800',
    },
    info: {
      main: '#2196f3',
    },
    success: {
      main: '#4caf50',
    },
  },
});

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
    <div className='home'>
      <div className='title'>RINE</div>
      <div className='login'>
      <Link to="/login">
        <button>ログイン</button>
      </Link>
      </div>
      <div className='signin'>
      <Link to="/signup">
        <button>新規登録</button>
      </Link>
      </div>
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
      const res = await fetch('http://localhost:8080/login', {
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
        history.push(`/rooms`);
        console.log(data)
      } else {
        alert("ログインエラーが発生しました。やり直してください。");
      }
    } catch (err) {
      setError((err as Error).message);
    }
  };

  return (
    <div className='loginscreen'>
      <div className='title'>ログイン</div>
      <div className='loginname'>
      <input value={username} onChange={(e: ChangeEvent<HTMLInputElement>) => setUsername(e.target.value)} placeholder="ユーザーネーム" />
      </div>
      <div className='loginpassword'>
      <input value={password} onChange={(e: ChangeEvent<HTMLInputElement>) => setPassword(e.target.value)} placeholder="パスワード" type="password" />
      </div>
      <div>
      <button onClick={login}>ログイン</button>
      </div>
      <div>
      <button onClick={() => history.push('/signup')}>新規登録</button>
      </div>
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
      const res = await fetch('http://localhost:8080/signup', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ UserName: username, UserPassword: password }),
      });

      if (res.ok) {
        history.push('/');
      } else {
        alert("そのユーザー名はもう使われているので、別の名前を設定してください。");
      }
    } catch (err) {
      setError((err as Error).message);
    }
  };

  return (
    <div className='signinscreen'>
      <div className='title'>新規登録</div>
      <div className='newuser'>
      <input value={username} onChange={(e: ChangeEvent<HTMLInputElement>) => setUsername(e.target.value)} placeholder="ユーザーネーム" />
      </div>
      <div className='newpassword'>
      <input value={password} onChange={(e: ChangeEvent<HTMLInputElement>) => setPassword(e.target.value)} placeholder="パスワード" type="password" />
      </div>
      <div className='signinbuttons'>
      <button onClick={() => history.goBack()}>戻る</button>
      <button onClick={signup}>登録</button>
      </div>
    </div>
  );
};

const Profile : FC =() => {
  const userName = localStorage.getItem('userName');
  const history = useHistory();
  return(
    <div className='profile_main'>
    
    <div className="profile">
      <div className='username_title'>ユーザーネーム</div>
      <div className='profile_username'>{userName}</div>
      <div className='exitprofile'>
        <button onClick={() => history.goBack()}>戻る</button>
      </div>
    </div>
    </div>
  )
};

const MakeRoom : FC =() => {
  const userId = localStorage.getItem('userId');
  const userName = localStorage.getItem('userName');
  const [opponent, setOpponent] = useState<string>('');
  const [error, setError] = useState<string | null>(null);
  const history = useHistory();

  const follow = async () => {
    try {
      const res = await fetch('http://localhost:8080/follow', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ UserId:userId, UserName: userName, OpponentUserName: opponent }),
      });

      if (res.ok) {
        history.push(`/rooms?userid=${userId}`);
      } else {
        alert("すでにそのユーザーをフォローしている、もしくはそのユーザーは存在しません。");
      }
    } catch (err) {
      setError((err as Error).message);
    }
  };

  return (
    <div className='followscreen'>
      <ThemeProvider theme={theme}>
      <div className='backbutton'>
      <Button onClick={()=>history.goBack()}><ArrowBackIcon fontSize="large"/></Button>  
      </div>
      <div className='follow'>
        <div className='title'>ユーザーをフォロー</div>
        <input value={opponent} onChange={(e: ChangeEvent<HTMLInputElement>) => setOpponent(e.target.value)} placeholder="ユーザーネームを入力してください" />
        <div className='followbuttons'>
        <button onClick={follow}>フォロー</button>
        </div>
      </div>
      </ThemeProvider>
    </div>
  );
};

const RoomList: FC = () => {
  const [rooms, setRooms] = useState<Room[]>([]);
  const [error, setError] = useState<string | null>(null);
  const userName = localStorage.getItem('userName');
  const userId = localStorage.getItem('userId');
  const history = useHistory();
  const location = useLocation();
  const queryParams = new URLSearchParams(location.search);

  useEffect(() => {
    const fetchRooms = async () => {
      try {
        const res = await fetch(`http://localhost:8080/rooms`,{
          method:'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ UserId:userId}),
        });
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

    if (location.pathname === "/rooms") {
      fetchRooms();
    }

    
  }, [location.pathname]);

  
  return (
    <div className='chatlist'>
      <div className='profile_useradd'>
        <div className='title'>チャットリスト</div>
        <div className='buttons'>
        {/* <div className='username'>{userName}</div> */}
        <ThemeProvider theme={theme}>
        <Link to="/follow">
          <Button  aria-label="follow">
            < PersonAddAltIcon color="primary" fontSize="large"/>
          </Button>
        </Link>
        <Link to="/">
          <Button onClick={()=>localStorage.clear()}><LogoutIcon fontSize="large"/></Button>  
        </Link> 
        </ThemeProvider>
        </div> 
      </div>
      <div className='rooms'>
        {rooms.map(room => (
          <div key={room.RoomId} onClick={() => history.push(`/chat/${room.RoomId}`)} className="room-list-item">
            {room.UserName1 == userName ? room.UserName2:room.UserName1}
          </div>
        ))}
      </div>
    </div>
  );
};



const ChatRoom: FC = () => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [newMessage, setNewMessage] = useState<string>('');
  const [error, setError] = useState<string | null>(null);
  const match = useRouteMatch<{ roomId: string }>('/chat/:roomId');
  const history = useHistory();

  useEffect(() => {
    const fetchMessages = async () => {
      try {
        const res = await fetch(`http://localhost:8080/fetchmessage?roomid=${match?.params.roomId}`);
        const data = await res.json();
        console.log(data)
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
  }, [match?.params.roomId]);

  const sendMessage = async () => {
    try {
      const userId = localStorage.getItem('userId');
      const userName = localStorage.getItem('userName');
      const res = await fetch('http://localhost:8080/sendmessage', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ MessageContent: newMessage, MessageTime: ` (${new Date().getMonth()+1}/${new Date().getDate()} ${new Date().getHours()}:${new Date().getMinutes()})`, UserId: userId, RoomId: match?.params.roomId, UserName: userName }),
      });
      const data = await res.json();

      if (res.ok) {
        setMessages([...messages, data]);
        setNewMessage('');
        console.log(data)
      } else {
        throw new Error('Failed to send message');
      }
    } catch (err) {
      setError((err as Error).message);
    }
  };
  

  return (
    <ThemeProvider theme={theme}>
    <div className='chatroom'>
  <div className='firstline'>
  <div className='exitchatroom'>
    <Button onClick={()=>history.goBack()}><ArrowBackIcon fontSize="large"/></Button>  
    </div>
    
  </div>
  <div className='message-list'>
    {messages && messages.map(message => (
      <div key={message.MessageTime} className="message-list-item">
        <div className='name'>{message.UserName}</div>
        <div className='content'>{message.MessageContent}</div> 
        <div className='time'>{message.MessageTime}</div>
      </div>
    ))}
  </div>
  <div className="message-input-container">
    <input value={newMessage} onChange={(e: ChangeEvent<HTMLInputElement>) => setNewMessage(e.target.value)} placeholder="メッセージを入力" />
    <Button color="secondary" onClick={sendMessage}><SendIcon fontSize="large"/></Button>
  </div>
</div>
</ThemeProvider>

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
        <Route path="/profile">
          <Profile />
        </Route>
        <Route path="/follow">
          <MakeRoom />
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
