import React from "react";
import { ChatMessage } from "../model";

export type ChatMessagesContextType = {
    messages: ChatMessage[];
    setMessages: (messages: ChatMessage[]) => void
}

export const ChatMessagesContext = React.createContext<ChatMessagesContextType | null>(null);

export const pushChatMessage = (messages: ChatMessage[], message: ChatMessage) => {
    const updatedMessages = [...messages];
    if (updatedMessages.length > 100) {
        updatedMessages.shift(); // Remove the oldest message if we exceed 100
    }
    updatedMessages.push(message);
    return updatedMessages;
}

export const pushPromptMessage = (messages: ChatMessage[], prompt: string) => {
    const message: ChatMessage = {
        id: crypto.randomUUID(),
        content: prompt,
        direction: 'outgoing',
        timestamp: new Date()
    };
    return pushChatMessage(messages, message);
}

export const pushAnswerMessage = (messages: ChatMessage[], answer: string) => {
    const message: ChatMessage = {
        id: crypto.randomUUID(),
        content: answer,
        direction: 'incoming',
        timestamp: new Date()
    };
    return pushChatMessage(messages, message);
}

export const pushPromptAndAnswerMessages = (messages: ChatMessage[], prompt: string, answer: string) => {
    let updatedMessages = pushPromptMessage(messages, prompt);
    updatedMessages = pushAnswerMessage(updatedMessages, answer);
    return updatedMessages;
}