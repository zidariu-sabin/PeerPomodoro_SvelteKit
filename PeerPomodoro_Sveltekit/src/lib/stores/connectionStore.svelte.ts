import { browser } from '$app/environment';
import { syncTimer, type TimerConfigurableData } from './pomodoroStore.svelte';

export type ConnectionState = 'connecting' | 'connected' | 'disconnected' | 'error';

interface SessionState {
    id: string | null;
    error: string | null;
}

class ConnectionStore {
    socket: WebSocket | null = null;
    state: ConnectionState = $state('disconnected');
    error: string | null = $state(null);
    session: SessionState = $state({ id: null, error: null });

    constructor() {}

    connect() {
        if (!browser) return;
        if (this.socket?.readyState === WebSocket.OPEN) return;

        this.state = 'connecting';
        this.error = null;

        try {
            this.socket = new WebSocket('ws://localhost:8080/connect');

            this.socket.onopen = () => {
                this.state = 'connected';
                console.log('WebSocket connected');
            };

            this.socket.onclose = () => {
                this.state = 'disconnected';
                this.socket = null;
                console.log('WebSocket disconnected');
            };

            this.socket.onerror = (event) => {
                this.state = 'error';
                this.error = 'WebSocket error occurred';
                console.error('WebSocket error:', event);
            };

            this.socket.onmessage = (event) => {
                try {
                    const msg = JSON.parse(event.data);
                    this.handleMessage(msg);
                } catch (e) {
                    console.error('Failed to parse message:', event.data);
                }
            };

        } catch (e) {
            this.state = 'error';
            this.error = e instanceof Error ? e.message : 'Unknown error';
        }
    }

    handleMessage(msg: any) {
        console.log('Received:', msg);
        switch (msg.type) {
            case 'session_created':
                this.session.id = msg.payload.session_id;
                break;
            case 'session_joined':
                this.session.id = msg.payload.session_id;
                if (msg.payload.timer) {
                    syncTimer(msg.payload.timer);
                }
                break;
            case 'timer_update':
                if (msg.payload.timer) {
                    syncTimer(msg.payload.timer);
                }
                break;
            case 'error':
                this.session.error = msg.payload.message;
                break;
        }
    }

    sendMessage(data: any) {
        if (this.socket?.readyState === WebSocket.OPEN) {
            this.socket.send(JSON.stringify(data));
        } else {
            console.warn('Cannot send message, socket not open');
        }
    }

    createSession(timerData: TimerConfigurableData) {
        // Map camelCase to snake_case for backend
        const payload = {
            work_time: timerData.workTime,
            break_time: timerData.breakTime,
            total_rounds: timerData.totalRounds
        };
        this.sendMessage({ type: 'create_session', payload });
    }

    joinSession(sessionId: string, userName: string) {
        this.sendMessage({
            type: 'join_session',
            payload: { session_id: sessionId, user_name: userName }
        });
    }

    startTimer() {
        this.sendMessage({ type: 'start_timer', payload: {} });
    }

    pauseTimer() {
        this.sendMessage({ type: 'pause_timer', payload: {} });
    }

    resetTimer() {
        this.sendMessage({ type: 'stop_timer', payload: {} });
    }
}

export const connection = new ConnectionStore();