import { browser } from '$app/environment';
import { syncTimer, timer, type TimerConfigurableData } from './pomodoroStore.svelte';
import { requestNotificationPermission, sendNotification } from '$lib/utils/notifications';
import { PUBLIC_WS_BASE_URL } from '$env/static/public';

export type ConnectionState = 'connecting' | 'connected' | 'disconnected' | 'error';

export interface User {
    id: string;
    name: string;
}

interface SessionState {
    id: string | null;
    error: string | null;
}

class ConnectionStore {
    socket: WebSocket | null = null;
    state: ConnectionState = $state('disconnected');
    error: string | null = $state(null);
    session: SessionState = $state({ id: null, error: null });
    clientId: string | null = $state(null);
    users: User[] = $state([]);

    constructor() {}

    connect() {
        if (!browser) return;
        if (this.socket?.readyState === WebSocket.OPEN || 
            this.socket?.readyState === WebSocket.CONNECTING) return;

        this.state = 'connecting';
        this.error = null;
        this.session.error = null;
        this.clientId = null;

        try {
            const wsUrl = `ws://${import.meta.env.PUBLIC_API_BASE_URL}/connect` || 'ws://localhost:8080/connect';
            this.socket = new WebSocket(wsUrl);

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
        switch (msg.type) {
            case 'session_created':
                this.session.id = msg.payload.session_id;
                break;
            case 'session_joined':
                this.session.id = msg.payload.session_id;
                this.clientId = msg.payload.client_id;
                if (msg.payload.timer) {
                    syncTimer(msg.payload.timer);
                }
                if (msg.payload.clients) {
                    this.users = msg.payload.clients;
                }
                break;
            case 'user_joined':
                const newUser = msg.payload.client;
                if (!this.users.find(u => u.id === newUser.id)) {
                    this.users.push(newUser);
                }
                break;
            case 'user_updated':
                const index = this.users.findIndex(u => u.id === msg.payload.client_id);
                if (index !== -1) {
                    // Force full array update for reactivity
                    const newUsers = [...this.users];
                    newUsers[index] = { ...newUsers[index], name: msg.payload.name };
                    this.users = newUsers;
                }
                break;
            case 'user_left':
                this.users = this.users.filter(u => u.id !== msg.payload.client_id);
                break;
            case 'timer_update':
                if (msg.payload.timer) {
                    const newTimer = msg.payload.timer;
                    
                    if (newTimer.is_completed && !timer.isCompleted) {
                        sendNotification("Session Completed!", "All rounds finished. Well done!");
                    } else if (!newTimer.is_completed) {
                        if (timer.isWorkPeriod && !newTimer.is_work_period) {
                             sendNotification("Break Time!", "Great job! Take a short break.");
                        } else if (!timer.isWorkPeriod && newTimer.is_work_period) {
                             sendNotification("Work Time!", "Break is over. Focus time!");
                        }
                    }

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
        requestNotificationPermission();
        // Map camelCase to snake_case for backend
        const payload = {
            work_time: timerData.workTime,
            break_time: timerData.breakTime,
            total_rounds: timerData.totalRounds
        };
        this.sendMessage({ type: 'create_session', payload });
    }

    joinSession(sessionId: string, userName: string) {
        requestNotificationPermission();
        this.sendMessage({
            type: 'join_session',
            payload: { session_id: sessionId, user_name: userName }
        });
    }

    updateName(name: string) {
        this.sendMessage({
            type: 'update_user',
            payload: { name }
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
