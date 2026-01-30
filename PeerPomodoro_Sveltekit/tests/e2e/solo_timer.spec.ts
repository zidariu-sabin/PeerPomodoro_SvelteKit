import { test, expect } from '@playwright/test';

test.describe('Solo Timer Notification Test', () => {
  test('should initialize timer and trigger notifications correctly', async ({ page }) => {
    // 0. Install clock immediately to ensure consistent time reference (Epoch 1970)
    await page.clock.install();

    // 1. Grant notification permissions
    await page.context().grantPermissions(['notifications']);

    // 2. Go to home page
    await page.goto('http://localhost:5173/');

    // 3. Spy on the Notification API
    // We'll store notifications in a global array in the browser context to assert them later
    await page.evaluate(() => {
        (window as any).notifications = [];
        const originalNotification = window.Notification;
        
        // Mock Notification constructor
        const MockNotification = function(title: string, options?: NotificationOptions) {
            (window as any).notifications.push({ title, options });
            return new (originalNotification as any)(title, options); // Call original if needed, or just mock
        } as any;
        
        // Copy static properties
        Object.assign(MockNotification, originalNotification);
        MockNotification.permission = 'granted';
        MockNotification.requestPermission = async () => 'granted';
        
        window.Notification = MockNotification;
    });

    // 4. Initialize Solo Timer: 1 min work, 1 min break, 4 repetitions
    await page.fill('#workMinutes', '1');
    await page.fill('#breakMinutes', '1');
    await page.fill('#totalRounds', '4');

    // 5. Start the timer
    // Note: The "Start Timer!" button is wrapped in an anchor tag.
    await page.click('text=Start Timer!');

    // 6. Verify we are on the timer page
    await expect(page).toHaveURL(/\/timer/);
    await expect(page.locator('h2')).toHaveText('Work');
    await expect(page.locator('.text-6xl')).toHaveText('01:00');

    // 7. Verify notifications over time
    // We use clock manipulation to speed up the test
    // Install the clock *after* page load if possible, or use page.clock which persists?
    // Playwright's page.clock.install() replaces Date, setTimeout, etc.
    
    // NOTE: Svelte 5 runes / stores might rely on microtasks or specific timing. 
    // fastForward should handle setInterval correctly.
    
    // --- End of Work Round 1 (1 min) ---
    // Fast forward 1 minute (plus a little buffer for execution)
    await page.clock.fastForward(60 * 1000 + 100);

    // Expect transition to Break
    // Check UI
    await expect(page.locator('h2')).toHaveText('Break');
    // Check Notification
    let notifications = await page.evaluate(() => (window as any).notifications);
    expect(notifications).toEqual(
        expect.arrayContaining([
            expect.objectContaining({ title: 'Break Time!' })
        ])
    );
    // Clear notifications for next assertion to be clean
    await page.evaluate(() => (window as any).notifications = []);

    // --- End of Break Round 1 (1 min) ---
    await page.clock.fastForward(60 * 1000 + 100);

    // Expect transition to Work (Round 2)
    await expect(page.locator('h2')).toHaveText('Work');
    await expect(page.locator('.text-xl.mb-4')).toContainText('Round 2 / 4');
    
    notifications = await page.evaluate(() => (window as any).notifications);
    console.log(notifications)
    expect(notifications).toEqual(
        expect.arrayContaining([
            expect.objectContaining({ title: 'Work Time!' })
        ])
    );
    await page.evaluate(() => (window as any).notifications = []);

    // --- Fast Forward through remaining rounds ---
    // Round 2 Work (1m) -> Break (1m) -> Round 3 Work (1m) -> Break (1m) -> Round 4 Work (1m) -> Break (1m)
    // Actually: 
    // R1 Work -> Break
    // R2 Work -> Break
    // R3 Work -> Break
    // R4 Work -> Completion (No break after last work usually? Or maybe break then complete? 
    // Let's check logic: if (timer.currentRound >= timer.totalRounds) -> Complete.
    // This check is in the "Break" transition logic?
    // logic: if (timer.isWorkPeriod) { ... becomes Break ... } else { if (current >= total) Complete; else Work ... }
    // So: 
    // R1 Work -> Break
    // Break -> R2 Work
    // R2 Work -> Break
    // Break -> R3 Work
    // R3 Work -> Break
    // Break -> R4 Work
    // R4 Work -> Break (Wait, logic says `if (timer.currentRound >= timer.totalRounds)` inside the `else` block of `isWorkPeriod`.
    // The `else` block handles "Break is ending, starting Work".
    // 
    // Let's trace carefully:
    // Start: Work, R=1.
    // Work Ends -> `handlePeriodComplete`: `isWorkPeriod` is true. Switch to Break. `nextSecondsRemaining` = breakTime. `isWorkPeriod` = false.
    // Break Ends -> `handlePeriodComplete`: `isWorkPeriod` is false.
    //    Check: `if (timer.currentRound >= timer.totalRounds)`? 1 >= 4? No.
    //    Switch to Work. `isWorkPeriod` = true. `currentRound++` (becomes 2).
    // ...
    // R4 Work Ends -> Break.
    // Break (Round 4) Ends -> `handlePeriodComplete`: `isWorkPeriod` is false.
    //    Check: `if (timer.currentRound >= timer.totalRounds)`? 4 >= 4? Yes!
    //    `timer.isCompleted = true`. Send "Session Completed!".
    
    // So the sequence is: W(1), B(1), W(2), B(2), W(3), B(3), W(4), B(4) -> Complete.
    
    // We already did W(1), B(1). We are at Start of W(2).
    // Remaining: W(2), B(2), W(3), B(3), W(4), B(4).
    // Total remaining time: 6 mins.
    
    await page.clock.fastForward(6 * 60 * 1000 + 1000);

    // Assert Completion
    await expect(page.locator('.text-2xl.text-green-600')).toContainText('All rounds completed!');
    
    notifications = await page.evaluate(() => (window as any).notifications);
    expect(notifications).toEqual(
        expect.arrayContaining([
            expect.objectContaining({ title: 'Session Completed!' })
        ])
    );

  });
});
