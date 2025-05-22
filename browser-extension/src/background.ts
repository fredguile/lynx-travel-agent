/// <reference types="@types/firefox-webext-browser"/>

// Listen for messages from content scripts
browser.runtime.onMessage.addListener(async (message, sender) => {
    if (message.action === 'capture_screenshot') {
        try {
            // Get the windowId from the sender, or fallback to the current window's id
            const windowId = sender.tab?.windowId || (await browser.windows.getCurrent()).id;
            if (windowId !== undefined) {
                const screenshot = await browser.tabs.captureVisibleTab(windowId, { format: 'png' });
                return Promise.resolve({ screenshot, mimeType: 'image/png' });
            }
        } catch (err) {
            return Promise.resolve({ error: (err instanceof Error ? err.message : String(err)) });
        }
    }
}); 