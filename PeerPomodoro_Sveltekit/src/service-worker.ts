/// <reference lib="webworker" />
/// <reference types="@sveltejs/kit" />

import type { WorkerMessage } from '$lib/types/workerMessages';
import { handlePeriodComplete } from '$lib/stores/pomodoroStore.svelte';
import { sendNotification } from '$lib/utils/notifications';

self.addEventListener('install', () => {
	console.log('Service worker installing123');
});

self.addEventListener('activate', () => {
	console.log('Service worker activated123');
});	


const sw = self as unknown as ServiceWorkerGlobalScope;

sw.addEventListener('install', (event) => {
    sw.skipWaiting();
});

sw.addEventListener('activate', (event) => {
    event.waitUntil(sw.clients.claim());
});

let timerId: any = null;

sw.addEventListener('message', (event) => {
    if (!event.data || typeof event.data !== 'object') return;

    const message = event.data as WorkerMessage;
    const { type, payload } = message;

    if (type === 'SCHEDULE_NOTIFICATION' && payload) {
        if (timerId) {
            clearTimeout(timerId);
            timerId = null;
        }

        const { delay, title, body } = payload;
        
        // If delay is effectively immediate, show it now, otherwise schedule
        if (delay <= 0) {
            sendNotification(title, body);
        } else {
            timerId = setTimeout(() => {
                sendNotification(title, body);
                handlePeriodComplete();
                timerId = null;
            }, delay);
        }

    } else if (type === 'CANCEL_NOTIFICATION') {
        if (timerId) {
            clearTimeout(timerId);
            timerId = null;
        }
    }
});

// function showNotification(title: string, body: string) {
//     sw.registration.showNotification(title, {
//         body,
//         requireInteraction: true, // Keeps it visible until user clicks
//         tag: 'pomodoro-timer' // Replaces older notifications with this tag
//     });
// }
