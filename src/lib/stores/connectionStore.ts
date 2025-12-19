
interface Subscriber {
    id: number,
    name: string,
    sessionId: Session,
}

interface Session {
    id: number,
}

function connection(session: Session){

    const subscribers = new Set()

    function subscribe(subscriber: Subscriber){
        subscribers.add(subscriber)

        return () => {
            subscribers.delete(subscriber)
        }
    }

    function notify(){
        // subscribers.forEach()
    }

    return { subscribe, notify}

}