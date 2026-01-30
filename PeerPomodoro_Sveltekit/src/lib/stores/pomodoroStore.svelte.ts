
import { requestNotificationPermission, sendNotification, scheduleNotification, cancelNotification } from '$lib/utils/notifications';

export interface TimerConfigurableData{
  workTime: number
  breakTime: number
  totalRounds: number
}

export interface TimerDisplayData{
  totalRounds: number
  currentRound: number
  secondsRemaining: number
  isWorkPeriod: boolean
  isRunning: boolean
  isCompleted: boolean
  isTimerCreated: boolean
  endTime: number | null
}

type PomodoroState = TimerConfigurableData & TimerDisplayData
export const timer: PomodoroState = $state({
  workTime: 25,
  breakTime: 5,
  totalRounds: 4,
  currentRound: 1,
  secondsRemaining: 0,
  isWorkPeriod: true,
  isRunning: false,
  isCompleted: false,
  isTimerCreated: false,
  endTime: null
})
let intervalId: number | null = null

export function setTimerData(timerData: TimerConfigurableData) {
  timer.workTime = timerData.workTime
  timer.breakTime = timerData.breakTime
  timer.totalRounds = timerData.totalRounds

  initializeTimer()
  startTimer()
}

export function syncTimer(backendTimer: any) {
  // When backend syncs, we stop our local scheduling because the server drives the state.
  // We cancel any local SW notifications to avoid double-firing or ghost notifications.
  cancelNotification();

  timer.workTime = backendTimer.work_time
  timer.breakTime = backendTimer.break_time
  timer.totalRounds = backendTimer.total_rounds
  timer.currentRound = backendTimer.current_round
  timer.secondsRemaining = backendTimer.seconds_remaining
  timer.isWorkPeriod = backendTimer.is_work_period
  timer.isRunning = backendTimer.is_running
  timer.isCompleted = backendTimer.is_completed
  timer.isTimerCreated = backendTimer.is_timer_created
  timer.endTime = null; // Backend drives the timer, no local target time.

  // Disable local interval in favor of backend updates
  if (intervalId !== null) {
    clearInterval(intervalId)
    intervalId = null
  }
}

export function refreshTimer() {
    tick();
}

function tick() {
  if (timer.endTime !== null) {
      const now = Date.now();
      const remaining = Math.ceil((timer.endTime - now) / 1000);
      
      timer.secondsRemaining = remaining;
      
      // Debug log (optional, reduced frequency or kept as requested)
      console.log("Seconds Remaining:", timer.secondsRemaining, "Real time: ", new Date().toLocaleTimeString());

      if (timer.secondsRemaining <= 0) {
          timer.secondsRemaining = 0;
          handlePeriodComplete();
      }
  } else {
      if (timer.isRunning && !timer.endTime) {
           timer.endTime = Date.now() + (timer.secondsRemaining * 1000);
      }
  }
}

  export function handlePeriodComplete() {
    console.log("Period complete. Handling transition...");
    // Determine next state properties
    let nextSecondsRemaining = 0;
    let title = "";
    let body = "";

    if (timer.isWorkPeriod) {
      // Transition to Break
      title = "Break Time!";
      body = "Great job! Take a short break.";
      nextSecondsRemaining = timer.breakTime * 60;
      timer.isWorkPeriod = false;
    } else {
      // Transition to Work or Complete
      if (timer.currentRound >= timer.totalRounds) {
        pause();
        timer.isCompleted = true;
        sendNotification("Session Completed!", "All rounds finished. Well done!");
        return;
      }
      
      title = "Work Time!";
      body = "Break is over. Focus time!";
      nextSecondsRemaining = timer.workTime * 60;
      timer.isWorkPeriod = true;
      timer.currentRound++;
    }
    // Update state
    timer.secondsRemaining = nextSecondsRemaining;
    
    // Send immediate notification for the transition
    // sendNotification(title, body);

    // Schedule NEXT notification
    if (timer.isRunning) {
        timer.endTime = Date.now() + (timer.secondsRemaining * 1000);
        scheduleNextNotification()
    }
  }

export function pause() {
  timer.isRunning = false
  timer.endTime = null;
  cancelNotification(); // Cancel any pending SW notification
  if (intervalId !== null) {
    clearInterval(intervalId)
    intervalId = null
  }
}

// Local-Only: Resets the timer state
export function resetLocalTimer() {
  pause()
  initializeTimer()
}

function initializeTimer() {
    timer.isTimerCreated = true
    timer.currentRound = 1
    timer.isCompleted = false
    timer.secondsRemaining = timer.workTime * 60
    timer.isWorkPeriod = true
    timer.isRunning = false
    timer.endTime = null;
  }

export function startTimer(){
  requestNotificationPermission()
  if (timer.isRunning) return
  if (intervalId !== null) clearInterval(intervalId)
  
  timer.isRunning = true
  
  // Initialize end time based on current secondsRemaining
  timer.endTime = Date.now() + (timer.secondsRemaining * 1000);

  intervalId = window.setInterval(() => {
      tick()
    }, 1000)
}

function scheduleNextNotification() {
    // Schedule the notification for when this period ends
  const nextTitle = timer.isWorkPeriod ? "Break Time!" : "Work Time!";
  const nextBody = timer.isWorkPeriod ? "Great job! Take a short break." : "Break is over. Focus time!";

  let targetTitle = nextTitle;
  let targetBody = nextBody;
  
  if (!timer.isWorkPeriod && timer.currentRound >= timer.totalRounds) {
      targetTitle = "Session Completed!";
      targetBody = "All rounds finished. Well done!";
  }

  scheduleNotification(targetTitle, targetBody, timer.secondsRemaining * 1000);
}

