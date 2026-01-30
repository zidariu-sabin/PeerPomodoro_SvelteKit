export interface ScheduleNotificationPayload {
    delay: number;
    title: string;
    body: string;
}

export interface WorkerMessage {
    type: 'SCHEDULE_NOTIFICATION' | 'CANCEL_NOTIFICATION' | 'TIMER_COMPLETED';
    payload?: ScheduleNotificationPayload;
}
