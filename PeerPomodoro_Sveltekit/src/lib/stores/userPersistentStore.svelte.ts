import { browser } from '$app/environment';

/**
 * STRUCTURE FOR PERSISTENT DATA:
 * To add new persistent features (like a Task List):
 * 1. Add a property with $state in the class.
 * 2. Add the property to the load() method with a fallback.
 * 3. Add a setter method that calls this.save().
 */

interface PersistentData {
    userName: string;
    theme: 'light' | 'dark' | 'system';
    soundEnabled: boolean;
    // tasks: Task[]; // Example for future features
}

const STORAGE_KEY = 'peerpomodoro_user_data';

class UserPersistentStore {
    // --- Identity ---
    userName = $state('User-' + Math.floor(Math.random() * 1000));

    // --- Preferences ---
    theme = $state<'light' | 'dark' | 'system'>('system');
    soundEnabled = $state(true);

    constructor() {
        if (browser) {
            this.load();
        }
    }

    private load() {
        try {
            const stored = localStorage.getItem(STORAGE_KEY);
            if (stored) {
                const data = JSON.parse(stored);
                // Load each property if it exists in the storage
                if (data.userName !== undefined) this.userName = data.userName;
                if (data.theme !== undefined) this.theme = data.theme;
                if (data.soundEnabled !== undefined) this.soundEnabled = data.soundEnabled;
            }
        } catch (e) {
            console.error('Failed to load user persistent data:', e);
        }
    }

    private save() {
        if (!browser) return;
        try {
            const data: PersistentData = {
                userName: this.userName,
                theme: this.theme,
                soundEnabled: this.soundEnabled
            };
            localStorage.setItem(STORAGE_KEY, JSON.stringify(data));
        } catch (e) {
            console.error('Failed to save user persistent data:', e);
        }
    }

    // --- Setters (One for each piece of data to ensure persistence) ---

    setUserName(name: string) {
        this.userName = name;
        this.save();
    }

    setTheme(theme: PersistentData['theme']) {
        this.theme = theme;
        this.save();
    }

    setSoundEnabled(enabled: boolean) {
        this.soundEnabled = enabled;
        this.save();
    }
}

export const userPersistentStore = new UserPersistentStore();
