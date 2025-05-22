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
 * Gets the index of the given element among all elements of the same tag name in the document.
 * @param element - The HTMLElement to find the index for.
 * @returns The index of the element among elements with the same tag name.
 */
export function getElementIndexInDocument(element: HTMLElement) {
    const { tagName } = element;
    // query all elements with the same tag name
    const elements = document.querySelectorAll(tagName);
    // return the index of the element
    return Array.from(elements).indexOf(element);
}


/**
 * Checks if the clicked element is whitelisted for AI actions based on the current URL, tag name, and element index.
 * @param currentUrl The current page URL
 * @param event The MouseEvent from the click
 * @returns true if the element is whitelisted, false otherwise
 */
export function isWhitelistedAIField(currentUrl: string, event: MouseEvent): boolean {
    const el = event.target as HTMLElement | null;
    if (!el) return false;
    const whitelistForUrl = (WHITELISTED_AI_FIELDS as WhitelistedAIFields)[currentUrl];
    if (!whitelistForUrl) return false;
    const tagName = el.tagName;
    const whitelistedIndexes = whitelistForUrl[tagName as keyof typeof whitelistForUrl];
    if (!whitelistedIndexes) return false;
    const idx = getElementIndexInDocument(el);
    return whitelistedIndexes.includes(idx);
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
