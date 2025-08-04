import { Answer } from "./Answer";
import { Prompt } from "./Prompt";

export const usePromptService = () => (
    async (prompt: Prompt): Promise<Answer> => {
        const response = await fetch(`/prompt?user-id=${prompt.userId}&question=${prompt.question}`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            }
        });

        if (!response.ok) {
            throw new Error("Failed to save prompt");
        }

        return response.json();
    }
)