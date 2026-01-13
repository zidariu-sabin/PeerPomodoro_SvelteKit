import { browser } from '$app/environment';
import { PUBLIC_WS_BASE_URL, PUBLIC_API_BASE_URL } from '$env/static/public';

export type ConnectionState = 'connecting' | 'connected' | 'disconnected' | 'error';
type MessageHandler = (msg: any) => void;

class SocketService {
    private socket: WebSocket | null = null;
    private messageHandler: MessageHandler | null = null;
    
    // State is managed here but subscribed to by the store
    public state: ConnectionState = 'disconnected';
    public error: string | null = null;
    
    // Callbacks for state changes
    private onStateChange: ((state: ConnectionState, error: string | null) => void) | null = null;

    constructor() {}

    setStateCallback(callback: (state: ConnectionState, error: string | null) => void) {
        this.onStateChange = callback;
    }

    setMessageHandler(handler: MessageHandler) {
        this.messageHandler = handler;
    }

    connect() {
        if (!browser) return;
        if (this.socket?.readyState === WebSocket.OPEN || 
            this.socket?.readyState === WebSocket.CONNECTING) return;

        this.updateState('connecting', null);

        try {
            // Determine WebSocket URL: use env var if set, otherwise derive from current location or default
            let wsUrl = PUBLIC_WS_BASE_URL;
            if (!wsUrl && PUBLIC_API_BASE_URL) {
                 wsUrl = `ws://${PUBLIC_API_BASE_URL.replace('http://', '').replace('https://', '')}`;
            }
            if (!wsUrl) {
                // Fallback if env vars are missing (though they should be there)
                wsUrl = 'ws://localhost:8080';
            }
            
            this.socket = new WebSocket(`${wsUrl}/connect`);

            this.socket.onopen = () => {
                this.updateState('connected', null);
                console.log('WebSocket connected');
            };

            this.socket.onclose = () => {
                this.updateState('disconnected', null);
                this.socket = null;
                console.log('WebSocket disconnected');
            };

            this.socket.onerror = (event) => {
                this.updateState('error', 'WebSocket error occurred');
                console.error('WebSocket error:', event);
            };

            this.socket.onmessage = (event) => {
                try {
                    const msg = JSON.parse(event.data);
                    if (this.messageHandler) {
                        this.messageHandler(msg);
                    }
                } catch (e) {
                    console.error('Failed to parse message:', event.data);
                }
            };

        } catch (e) {
            this.updateState('error', e instanceof Error ? e.message : 'Unknown error');
        }
    }

    send(data: any) {
        if (this.socket?.readyState === WebSocket.OPEN) {
            this.socket.send(JSON.stringify(data));
        } else {
            console.warn('Cannot send message, socket not open');
        }
    }

    private updateState(newState: ConnectionState, newError: string | null) {
        this.state = newState;
        this.error = newError;
        if (this.onStateChange) {
            this.onStateChange(newState, newError);
        }
    }
}

export const socketService = new SocketService();
