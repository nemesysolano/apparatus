export const Strings = {
    removeThinkingFromAnswer: (answer: string) => {
        const start = answer.indexOf('<think>')
        const end = answer.indexOf('</think>') + 8;
        if (start === -1 || end === -1) {
            return answer;
        }
        return answer.substring(0, start) + answer.substring(end);
    }
}