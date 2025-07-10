/// <reference types="@types/firefox-webext-browser"/>
import { debounce } from 'lodash';
import { analyseScreenContext } from './backend/analyseScreenContext';
import { LYNX_ELEMENTS_TO_WATCH, WHITELISTED_URLS } from './constants';
import { sessionStorage } from './sessionStorage';
// import { wrapElementWithAutoSuggest } from './ui';
import { createLogger, createContentChangeObserver, findElementLabel, isWhitelistedAIField, waitForElements, takeScreenshot } from './utils';

const log = createLogger('contentScript');

log('Content script loaded');

// Flag to prevent infinite loops during our own DOM modifications
let isWrappingDOMFields = false;

async function onPageRefreshed() {
    const currentUrl = location.href;
    log('onPageRefreshed:', currentUrl);

    // Store current URL in browser storage for global access
    await sessionStorage.currentUrl.set(currentUrl);

    // Skip if we're currently modifying the DOM to prevent infinite loops
    if (isWrappingDOMFields) {
        log('Skipping onPageRefreshed - currently modifying DOM');
        return;
    }

    // If currentUrl is whitelisted, take screenshot and analyse context
    if (WHITELISTED_URLS.includes(currentUrl)) {
        const screenshot = await takeScreenshot();
        const { fields } = await analyseScreenContext(currentUrl, screenshot);
        await sessionStorage.currentFields.set(fields);
    }

    // Initialize currentBookingRef from storage
    let currentBookingRef = await sessionStorage.currentBookingRef.get();

    if (currentUrl == 'https://www.lynx-reservations.com/lynx/#FILE_DETAILS') {
        // Retrieve booking reference from read-only fields
        const readOnlyFields = document.evaluate('//*[@class="readOnlyField"]', document, null, XPathResult.ORDERED_NODE_SNAPSHOT_TYPE, null);
        for (let i = 0; i < readOnlyFields.snapshotLength; i++) {
            const element = readOnlyFields.snapshotItem(i) as HTMLElement;
            const label = findElementLabel(element);

            if (label?.includes('File') || label?.includes('Quote Number')) {
                const quoteNumber = element.innerText;
                if (quoteNumber.startsWith('FT')) {
                    currentBookingRef = quoteNumber;
                    // Store in browser storage for global access
                    await sessionStorage.currentBookingRef.set(currentBookingRef);
                    log('Retrieved booking reference:', currentBookingRef);
                    break;
                }
            }
        }
    }

    // Set flag to indicate we're modifying DOM
    // isWrappingDOMFields = true;

    // for (const element of document.querySelectorAll('input, textarea') as NodeListOf<HTMLElement>) {
    //     if (isWhitelistedAIField(currentUrl, element) && currentBookingRef) {
    //         // Check if element is already wrapped by looking for our placeholder
    //         const isAlreadyWrapped = element.closest('[id^="ai-auto-suggest-placeholder-"]') !== null;
    //         if (!isAlreadyWrapped) {
    //             wrapElementWithAutoSuggest(element);
    //         }
    //     }
    // }

    // // Reset flag after DOM modifications are complete
    // isWrappingDOMFields = false;
}

// Detect page refreshes (F5, refresh button, etc.)
document.addEventListener('DOMContentLoaded', () => {
    log('DOMContentLoaded event detected - page refreshed');
    onPageRefreshed();
});

// Fallback for cases where DOMContentLoaded has already fired
if (document.readyState === 'loading') {
    // DOMContentLoaded will fire
    log('Document still loading, DOMContentLoaded will handle refresh detection');
} else {
    // DOMContentLoaded has already fired, run immediately
    log('Document already loaded, running onPageRefreshed immediately');
    onPageRefreshed();
}

// Listen for popstate (back/forward navigation)
window.addEventListener('popstate', onPageRefreshed);

// Monkey-patch pushState and replaceState to detect SPA navigation
(['pushState', 'replaceState'] as (keyof History)[]).forEach((method) => {
    const original = history[method] as (...args: any[]) => any;
    (history as any)[method] = function (this: History, ...args: any[]): any {
        const result = original.apply(this, args);
        // Use MutationObserver to wait for DOM updates after navigation
        waitForElements(LYNX_ELEMENTS_TO_WATCH, onPageRefreshed);

        return result;
    } as History[typeof method];
});

// Create content change observer to detect significant content changes (dynamic refreshes)
createContentChangeObserver(debounce(() => {
    // Only trigger if we're not currently modifying DOM ourselves
    if (!isWrappingDOMFields) {
        log('Significant content change detected - likely dynamic refresh');
        onPageRefreshed();
    }
}, 1e3));