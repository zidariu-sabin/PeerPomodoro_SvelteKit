import { browser } from '$app/environment';

export type ConnectionState = 'connecting' | 'connected' | 'disconnected' | 'error';

class ConnectionStore {
    socket: WebSocket | null = null;
    state: ConnectionState = $state('disconnected');
    error: string | null = $state(null);

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
                this.sendMessage({ type: 'join', payload: { connection_status: "connection_successfull" } });
            };

            this.socket.onclose = () => {
                this.state = 'disconnected';
                this.socket = null;
                console.log('WebSocket disconnected');
                this.sendMessage({ type: 'leave', payload: { connection_status: "connection_interrupted" } });
            };

            this.socket.onerror = (event) => {
                this.state = 'error';
                this.error = 'WebSocket error occurred';
                console.error('WebSocket error:', event);
            };
            //on received backed message
            this.socket.onmessage = (event) => {
                // console.log('Message received:', event.data);
            };

        } catch (e) {
            this.state = 'error';
            this.error = e instanceof Error ? e.message : 'Unknown error';
        }
    }

    sendMessage(data: any) {
        if (this.socket?.readyState === WebSocket.OPEN) {
            this.socket.send(JSON.stringify(data));
        } else {
            console.warn('Cannot send message, socket not open');
        }
    }
}

export const connection = new ConnectionStore();
