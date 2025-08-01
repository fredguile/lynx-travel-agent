import debug from 'debug';

import { WHITELISTED_AI_FIELDS } from './constants';
import type { WhitelistedAIFields } from './types';

/**
 * Creates a logger instance with the specified namespace using the debug library.
 * If the DEBUG environment variable is set and not present in localStorage, it sets it in localStorage.
 * @param namespace - The namespace for the logger.
 * @returns A debug logger instance.
 */
export function createLogger(namespace: string) {
    if (process.env.DEBUG && !localStorage.DEBUG) {
        localStorage.DEBUG = process.env.DEBUG;
    }
    return debug(namespace);
}

/**
 * Checks if the clicked element is whitelisted for AI actions based on the current URL, tag name, and element index.
 * @param currentUrl The current page URL
 * @param event The MouseEvent from the click
 * @returns true if the element is whitelisted, false otherwise
 */
export function isWhitelistedAIField(currentUrl: string, element: HTMLElement): boolean {
    const whitelistForUrl = (WHITELISTED_AI_FIELDS as WhitelistedAIFields)[currentUrl];
    if (!whitelistForUrl) return false;
    const tagName = element.tagName;
    const whitelistedLabels = whitelistForUrl[tagName as keyof typeof whitelistForUrl];
    if (!whitelistedLabels) return false;
    const label = findElementLabel(element);
    if (!label) return false;
    return whitelistedLabels.includes(label);
}

/**
 * Converts a base64 image string to a Blob object.
 * @param base64Image The base64 image string (with or without data URL prefix)
 * @returns Blob representing the image
 */
export function base64ImageToBlob(base64Image: string): Blob {
    let byteString: string;
    if (base64Image.startsWith('data:image/')) {
        byteString = atob(base64Image.split(',')[1]);
    } else {
        byteString = atob(base64Image || '');
    }
    const ab = new ArrayBuffer(byteString.length);
    const ia = new Uint8Array(ab);
    for (let i = 0; i < byteString.length; i++) {
        ia[i] = byteString.charCodeAt(i);
    }
    return new Blob([ab], { type: 'image/png' });
}

/**
 * Infers the label for a given form or read-only element.
 * - Only works for form elements (input, select, textarea, etc.) and read-only fields.
 * - If the element has an ID, searches for a <label> with a matching 'for' attribute.
 * - If no ID and inside a table, looks for a previous column <div> with class 'fieldLabel'.
 * @param el The HTMLElement to find the label for.
 * @returns The inferred label string, or null if not found.
 */
export function findElementLabel(el: HTMLElement): string | null {
    // Only consider form elements
    const formTags = ['INPUT', 'SELECT', 'TEXTAREA', 'BUTTON', 'FIELDSET', 'DIV'];
    if (!formTags.includes(el.tagName)) {
        return null;
    }

    // If element has an ID, try to find a <label for="...">
    const id = el.getAttribute('id');
    if (el.tagName !== 'DIV' && id) {
        const label = document.querySelector(`label[for="${id}"]`);
        if (label) {
            return label.textContent?.trim() || null;
        }
    }

    // If no ID and inside a table, look for previous column with .fieldLabel
    let parent: HTMLElement | null = el.parentElement;
    while (parent && parent.tagName !== 'TABLE') {
        parent = parent.parentElement;
    }
    if (parent && parent.tagName === 'TABLE') {
        // Find the closest ancestor <tr> and <td>
        let td: HTMLElement | null = el;
        while (td && td.tagName !== 'TD') {
            td = td.parentElement;
        }
        if (td && td.parentElement) {
            const tr = td.parentElement;
            // Find all <td> siblings before this one
            const tds = Array.from(tr.children);
            const idx = tds.indexOf(td);
            // First, look for a div.fieldLabel in previous columns
            for (let i = idx - 1; i >= 0; i--) {
                const prevTd = tds[i] as HTMLElement;
                const labelDiv = prevTd.querySelector('div.fieldLabel');
                if (labelDiv && labelDiv.textContent) {
                    return labelDiv.textContent.trim();
                }
            }
            // If not found, check the cell above in the previous row (same column index)
            const prevTr = tr.previousElementSibling as HTMLElement | null;
            if (prevTr && prevTr.tagName === 'TR' && idx >= 0 && idx < prevTr.children.length) {
                const aboveTd = prevTr.children[idx] as HTMLElement;
                if (aboveTd.classList.contains('fieldLabel')) {
                    return aboveTd.textContent?.trim() || null;
                }
            }
        }
    }

    return null;
}



/**
 * Returns true if the element and all its ancestors are visible (not display: none).
 * @param el The HTMLElement to check.
 * @returns boolean
 */
export function isElementVisible(el: HTMLElement): boolean {
    let current: HTMLElement | null = el;
    while (current) {
        if (current instanceof HTMLElement) {
            const style = window.getComputedStyle(current);
            if (style.display === 'none') {
                return false;
            }
        }
        current = current.parentElement;
    }
    return true;
}

/**
 * Waits for all specified elements to appear in the DOM before executing a callback.
 * Uses MutationObserver to watch for DOM changes and executes the callback when all
 * selectors successfully find their corresponding elements, or after a timeout period.
 * 
 * @param selectors - Array of CSS selectors to wait for. All selectors must find elements before callback is executed.
 * @param callback - Function to execute when all elements are found or timeout is reached.
 * @param timeout - Maximum time to wait in milliseconds before giving up (default: 3000ms).
 * 
 * @example
 * ```typescript
 * // Wait for a form and button to appear
 * waitForElements(['#myForm', '.submit-button'], () => {
 *   console.log('Both elements are now present');
 * }, 5000);
 * ```
 */
export function waitForElements(selectors: string[], callback: () => void, timeout = 3000) {
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

/**
 * Creates a MutationObserver to detect significant content changes (dynamic refreshes).
 * The observer watches for changes that indicate a content refresh, such as:
 * - Multiple nodes being added/removed
 * - Changes to readOnlyField elements, inputs, or textareas
 * 
 * @param onContentChange - Callback function to execute when significant content changes are detected
 * @returns The MutationObserver instance that has been started
 * 
 * @example
 * ```typescript
 * const observer = createContentChangeObserver(() => {
 *   console.log('Content has changed significantly');
 *   // Handle the content change
 * });
 * 
 * // Later, if needed:
 * observer.disconnect();
 * ```
 */
export function createContentChangeObserver(onContentChange: () => void): MutationObserver {
    const contentObserver = new MutationObserver((mutations) => {
        let significantChange = false;

        for (const mutation of mutations) {
            // Check for significant changes that indicate content refresh
            if (mutation.type === 'childList') {
                // If multiple nodes were added/removed, likely a content refresh
                if (mutation.addedNodes.length > 3 || mutation.removedNodes.length > 3) {
                    significantChange = true;
                    break;
                }

                // Check if readOnlyField elements were added/removed (key indicators for this app)
                for (const node of [...mutation.addedNodes, ...mutation.removedNodes]) {
                    if (node instanceof Element) {
                        if (node.classList?.contains('readOnlyField') ||
                            node.querySelector?.('.readOnlyField') ||
                            node.tagName === 'INPUT' ||
                            node.tagName === 'TEXTAREA') {
                            significantChange = true;
                            break;
                        }
                    }
                }
            }

            if (significantChange) break;
        }

        if (significantChange) {
            onContentChange();
        }
    });

    // Start observing body for content changes
    contentObserver.observe(document.body, {
        childList: true,
        subtree: true,
        attributes: false, // Don't watch attribute changes to avoid noise
        characterData: false // Don't watch text content changes to avoid noise
    });

    return contentObserver;
}

