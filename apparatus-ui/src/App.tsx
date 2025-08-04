import React, { useEffect, useState } from 'react';
import './App.css';
import { ChatView } from './views/chat/ChatView';
import airplaneLogo from './airplane.png';
import { ChatMessage, usePromptService } from './model';
import { ChatMessagesContext } from './context';
import { Strings } from './utils';
import { pushPromptAndAnswerMessages } from './context/ChatMessagesContext';

type AppState = {
  prompt: string
  sendDisabled: boolean
}

function App() {
  const promptService = usePromptService();
  const [contextMessages, setContexMessages] = useState<ChatMessage[]>([]);
  const [state, setState] = useState<AppState>({
    prompt: '',
    sendDisabled: true
  })

  const onPromptChange = (event: React.ChangeEvent<HTMLInputElement>) =>  setState({ 
    ...state, 
    prompt: event.target.value,
    sendDisabled: event.target.value.trim() === ''
  })

  const onSendClick = () => {
    (async () => {
      let result: ChatMessage[] = contextMessages;

      setState({ ...state, sendDisabled: true })
      try {
        const response = await promptService({ userId: 'rafael.solano', question: state.prompt })
        const answer = Strings.removeThinkingFromAnswer(response.answer)
        result = pushPromptAndAnswerMessages(contextMessages, state.prompt, answer)
      } catch (error: Error | any) {
        result = pushPromptAndAnswerMessages(contextMessages, state.prompt, error.message || 'Unknown error')
      } finally {
        setContexMessages(result);
      }
    })()
  }

  useEffect(() => {
    setState({ prompt: '', sendDisabled: false })
  }, [contextMessages]);

  return (
    <div className="app">
      <header>
          <img src="./torot-96.png" className="logo" alt="Apparatus Logo" />
      </header>
      <main>
        <ChatMessagesContext.Provider value={{ messages: contextMessages, setMessages: setContexMessages }}>
          <ChatView />
        </ChatMessagesContext.Provider>
      </main>
      <footer>
        <input type='text' placeholder='Type your message here...' onChange={onPromptChange} />
        <img src={airplaneLogo} alt='send' className={state.sendDisabled ? 'disabled' : ''} onClick={onSendClick} />
      </footer>
    </div>
  );
}

export default App;
