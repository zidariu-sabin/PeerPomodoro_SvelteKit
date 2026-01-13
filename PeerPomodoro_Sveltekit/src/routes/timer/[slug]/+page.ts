import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params, fetch }) => {
    const sessionId = params.slug;

    try {
        const response = await fetch(`http://localhost:8080/session/${sessionId}`);
        
        if (response.status === 404) {
            return {
                isValid: false,
                error: "Session not found or expired"
            };
        }

        if (!response.ok) {
            return {
                isValid: false,
                error: "Failed to validate session"
            };
        }

        return {
            isValid: true,
            sessionId: sessionId
        };
    } catch (e) {
        console.error("Session validation error:", e);
        return {
            isValid: false,
            error: "Could not connect to server"
        };
    }
};
