/// <reference types="@types/firefox-webext-browser"/>
import { ENDPOINTS } from './constants';
import { wrapElementWithAutoSuggest } from './ui';
import { createLogger, isWhitelistedAIField, base64ImageToBlob } from './utils';

const log = createLogger('contentScript');

log('Content script loaded');

let currentBookingRef: string | null = null;

async function onPageRefreshed() {
    const currentUrl = location.href;
    log('Page changed:', currentUrl);

    if (currentUrl == 'https://www.lynx-reservations.com/lynx/#FILE_DETAILS') {
        const response = await browser.runtime.sendMessage({ action: 'capture_screenshot' });
        const blob = base64ImageToBlob(response.screenshot);

        log('analysing screen context', currentUrl);

        const formData = new FormData();
        formData.append('screenshot', blob, 'screenshot.png');
        let res = await fetch(`${ENDPOINTS.ANALYSE_BOOKING_REF}?currentUrl=${encodeURIComponent(currentUrl)}`, {
            method: 'POST',
            body: formData,
        });
        const screenContext = await res.text();

        log('screen context:', screenContext);

        if (!screenContext.includes('No booking reference found')) {
            currentBookingRef = screenContext;
        }
    }

    for (const element of document.querySelectorAll('input, textarea') as NodeListOf<HTMLElement>) {
        if (isWhitelistedAIField(currentUrl, element) && currentBookingRef) {
            wrapElementWithAutoSuggest(currentUrl, currentBookingRef, element);
        }
    }
}

// Helper to wait for elements to appear in the DOM
function waitForElements(selectors: string[], callback: () => void, timeout = 3000) {
    const start = Date.now();
    const observer = new MutationObserver(() => {
        const allPresent = selectors.every(sel => document.querySelector(sel));
        if (allPresent) {
            observer.disconnect();
            callback();
        } else if (Date.now() - start > timeout) {
            observer.disconnect();
            callback(); // Optionally call anyway after timeout
        }
    });
    observer.observe(document.body, { childList: true, subtree: true });

    // Initial check in case elements are already present
    if (selectors.every(sel => document.querySelector(sel))) {
        observer.disconnect();
        callback();
    }
}

// Listen for popstate (back/forward navigation)
window.addEventListener('popstate', onPageRefreshed);

// Monkey-patch pushState and replaceState to detect SPA navigation
(['pushState', 'replaceState'] as (keyof History)[]).forEach((method) => {
    const original = history[method] as (...args: any[]) => any;
    (history as any)[method] = function (this: History, ...args: any[]): any {
        const result = original.apply(this, args);
        // Use MutationObserver to wait for DOM updates after navigation
        waitForElements([
            "div.readOnlyField",
            "input",
            "textarea"
        ], onPageRefreshed);
        return result;
    } as History[typeof method];
});

// Initial extraction
onPageRefreshed();
