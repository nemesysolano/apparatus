import React from 'react'
import './Scroll.css'

export type ScrollProperties = {
    className?: string,
    children?: React.ReactNode | React.ReactNode[],
}

export const Scroll: React.FC<ScrollProperties> = (properties: ScrollProperties) => {
    const className = `scroll-container ${properties.className || ' '}`.trim()

    return (
        <div className={className}>
            <div className="scroll-content">
                {properties.children}
            </div>
        </div>
    )
}
