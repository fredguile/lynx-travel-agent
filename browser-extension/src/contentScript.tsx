/// <reference types="@types/firefox-webext-browser"/>
import { wrapElementWithAutoSuggest } from './ui';
import { createLogger, isWhitelistedAIField } from './utils';

const log = createLogger('contentScript');

log('Content script loaded');

// --- Page Change Detection ---
let currentUrl = location.href;

function onPageChange() {
    currentUrl = location.href;
    log('Page changed:', currentUrl);

    for (const element of document.querySelectorAll('input, textarea') as NodeListOf<HTMLElement>) {
        if (isWhitelistedAIField(currentUrl, element)) {
            wrapElementWithAutoSuggest(currentUrl, element);
        }
    }
}

// Listen for popstate (back/forward navigation)
window.addEventListener('popstate', onPageChange);

// Monkey-patch pushState and replaceState to detect SPA navigation
(['pushState', 'replaceState'] as (keyof History)[]).forEach((method) => {
    const original = history[method] as (...args: any[]) => any;
    (history as any)[method] = function (this: History, ...args: any[]): any {
        const result = original.apply(this, args);
        setTimeout(onPageChange, 0); // Call after state changes
        return result;
    } as History[typeof method];
});

// Initial extraction
onPageChange();
