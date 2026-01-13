import { syncTimer, timer, type TimerConfigurableData } from './pomodoroStore.svelte';
import { requestNotificationPermission, sendNotification } from '$lib/utils/notifications';
import { socketService, type ConnectionState } from '$lib/services/socketService';

export interface User {
    id: string;
    name: string;
}

interface SessionState {
    id: string | null;
    error: string | null;
}

class SessionStore {
    state: ConnectionState = $state('disconnected');
    error: string | null = $state(null);
    session: SessionState = $state({ id: null, error: null });
    clientId: string | null = $state(null);
    users: User[] = $state([]);

    constructor() {
        // Subscribe to socket state changes
        socketService.setStateCallback((newState, newError) => {
            this.state = newState;
            this.error = newError;
            // Clear session error on connect/disconnect
            if (newState !== 'connected') {
                this.session.error = null;
            }
        });

        // Register message handler
        socketService.setMessageHandler((msg) => this.handleMessage(msg));
    }

    connect() {
        socketService.connect();
    }

    private handleMessage(msg: any) {
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

    createSession(timerData: TimerConfigurableData) {
        requestNotificationPermission();
        const payload = {
            work_time: timerData.workTime,
            break_time: timerData.breakTime,
            total_rounds: timerData.totalRounds
        };
        socketService.send({ type: 'create_session', payload });
    }

    joinSession(sessionId: string, userName: string) {
        requestNotificationPermission();
        socketService.send({
            type: 'join_session',
            payload: { session_id: sessionId, user_name: userName }
        });
    }

    updateName(name: string) {
        socketService.send({
            type: 'update_user',
            payload: { name }
        });
    }

    startTimer() {
        socketService.send({ type: 'start_timer', payload: {} });
    }

    pauseTimer() {
        socketService.send({ type: 'pause_timer', payload: {} });
    }

    resetTimer() {
        socketService.send({ type: 'stop_timer', payload: {} });
    }
}

export const sessionStore = new SessionStore();
