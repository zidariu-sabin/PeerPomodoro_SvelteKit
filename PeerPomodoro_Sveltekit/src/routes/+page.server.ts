
import {timer, setTimerData, type TimerData} from '$lib/stores/pomodoroStore.svelte'

// can be further developed for the collaborative session in order to create websockets 

// export const actions = {
//     default: async (request: Request) => {
//         const data = await request.formData();
//         let formData: TimerData
//         data.entries().forEach((field: {key, value})=> { formData[field.key] = field.value })
//     }   
// }