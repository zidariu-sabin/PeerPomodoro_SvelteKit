import { goto } from '$app/navigation';
import type { TimerConfigurableData } from "$lib/stores/pomodoroStore.svelte";
import { PUBLIC_API_BASE_URL } from '$env/static/public';

export async function createSession(timerData: TimerConfigurableData) {
    try {
        const response = await fetch(`${PUBLIC_API_BASE_URL}/create-session`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                work_time: timerData.workTime,
                break_time: timerData.breakTime,
                total_rounds: timerData.totalRounds
            })
        });

        if (!response.ok) {
            throw new Error('Failed to create session');
        }

        const data = await response.json();
        if (!data.session_id) {
            throw new Error('Invalid response: missing session_id');
        }
        await goto(`/timer/${data.session_id}`);
    } catch (error) {
        console.error('Error creating session:', error);
        throw error; // Re-throw so caller can handle UI if needed
    }
}
