import React, { useContext } from 'react'
import { BaloonMessage, Scroll } from '../../components'
import { ChatMessagesContext } from '../../context'

export const ChatView: React.FC = () => {
    const context = useContext(ChatMessagesContext)
 
    return (
        <Scroll>
            {context!!.messages.map((message) => (
                <BaloonMessage key={message.id} message={message}/>
            ))}
        </Scroll>
    )
}