import React, { createRef, FC, useEffect } from 'react'
import './BaloonMessage.css'
import { ChatMessage } from '../../model'
import incomingImage from './incoming.jpeg' 
import outgoingImage from './outgoing.jpeg'
import { marked } from 'marked' 
export type BaloonMessageProperties = {
    message: ChatMessage
    className?: string,
}

export const BaloonMessage: FC<BaloonMessageProperties> = (properties: BaloonMessageProperties) => {
    const className = `baloon-message ${properties.className || ''}`.trim()
    const contentAlignClass = properties.message.direction === 'outgoing' ? 'baloon-message-content-outgoing' : 'baloon-message-content-incoming'
    const image = properties.message.direction === 'outgoing' ? outgoingImage : incomingImage
    const imageAlt = properties.message.direction
    const spanRef = createRef<HTMLSpanElement>()
    const madkDownText = marked.parse(properties.message.content)

    useEffect(() => {
        if (spanRef.current) {
            spanRef!!.current!!.innerHTML = madkDownText as string
        }
    }, [madkDownText, spanRef])

    return (
        <div className={className}>
            <div className={`baloon-message-content ${contentAlignClass}`}>
                <p><img src={image} alt={imageAlt}/><span ref={spanRef}>{madkDownText}</span></p>
            </div>
        </div>
    )
}