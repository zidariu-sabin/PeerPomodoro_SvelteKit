

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
  timer.workTime = backendTimer.work_time
  timer.breakTime = backendTimer.break_time
  timer.totalRounds = backendTimer.total_rounds
  timer.currentRound = backendTimer.current_round
  timer.secondsRemaining = backendTimer.seconds_remaining
  timer.isWorkPeriod = backendTimer.is_work_period
  timer.isRunning = backendTimer.is_running
  timer.isCompleted = backendTimer.is_completed
  timer.isTimerCreated = backendTimer.is_timer_created

  // Disable local interval in favor of backend updates
  if (intervalId !== null) {
    clearInterval(intervalId)
    intervalId = null
  }
}

function tick() {
  if(timer.secondsRemaining > 0)  timer.secondsRemaining --
   else {
      handlePeriodComplete()
    }
  }

  function handlePeriodComplete() {
    if (timer.isWorkPeriod) {
      timer.isWorkPeriod = false
      timer.secondsRemaining = timer.breakTime * 60
    } else {
      timer.isWorkPeriod = true
      
      if (timer.currentRound >= timer.totalRounds) {
        pause()
        timer.isCompleted = true
        return
      }
      timer.currentRound++
      
      timer.secondsRemaining = timer.workTime * 60
    }
  }

export function pause() {
  timer.isRunning = false
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
  }

export function startTimer(){
  if (timer.isRunning) return
  if (intervalId !== null) clearInterval(intervalId)
  
  timer.isRunning = true
  intervalId = window.setInterval(() => {
      tick()
    }, 1000)
}

