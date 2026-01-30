import type { WorkerMessage } from '$lib/types/workerMessages';
import { dev } from '$app/environment';

export async function requestNotificationPermission(): Promise<NotificationPermission | null> {
    if (typeof window === 'undefined' || !('Notification' in window)) {
        console.warn('This browser does not support desktop notification');
        return null;
    }

    if (Notification.permission !== 'granted' && Notification.permission !== 'denied') {
        await Notification.requestPermission();
    }

    return Notification.permission;
}

export function sendNotification(title: string, body: string) {
    if (typeof window === 'undefined' || !('Notification' in window)) return;

    if (Notification.permission === 'granted') {
        // Fallback to main thread notification if SW not active (though SW is preferred)
        if (navigator.serviceWorker && navigator.serviceWorker.controller) {
            new Notification(title, {
                body,
            });
        } else {
             new Notification(title, {
                body,
            });
        }
    }
}

export function scheduleNotification(title: string, body: string, delayMs: number) {
    if (typeof window === 'undefined' || !('serviceWorker' in navigator)) return;

    if (navigator.serviceWorker.controller) {
        const message: WorkerMessage = {
            type: 'SCHEDULE_NOTIFICATION',
            payload: {
                delay: delayMs,
                title,
                body
            }
        };
        navigator.serviceWorker.controller.postMessage(message);
    }
}

export function cancelNotification() {
    if (typeof window === 'undefined' || !('serviceWorker' in navigator)) return;

    if (navigator.serviceWorker.controller) {
        const message: WorkerMessage = {
            type: 'CANCEL_NOTIFICATION'
        };
        navigator.serviceWorker.controller.postMessage(message);
    }
}
