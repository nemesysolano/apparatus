export type ChatMessageDirection = 'incoming' | 'outgoing'
export type ChatMessage = {
    id: string
    content: string
    direction: ChatMessageDirection
    timestamp: Date
}