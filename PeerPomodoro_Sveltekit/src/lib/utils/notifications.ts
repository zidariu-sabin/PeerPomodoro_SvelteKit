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
        new Notification(title, {
            body,
            icon: '/favicon.svg'
        });
    }
}