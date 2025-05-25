import { createRoot, Root } from 'react-dom/client';
import { createPortal } from 'react-dom';
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
 * Infers the label for a given form element.
 * - Only works for form elements (input, select, textarea, etc.).
 * - If the element has an ID, searches for a <label> with a matching 'for' attribute.
 * - If no ID and inside a table, looks for a previous column <div> with class 'fieldLabel'.
 * @param el The HTMLElement to find the label for.
 * @returns The inferred label string, or null if not found.
 */
export function findElementLabel(el: HTMLElement): string | null {
    // Only consider form elements
    const formTags = ['INPUT', 'SELECT', 'TEXTAREA', 'BUTTON', 'FIELDSET'];
    if (!formTags.includes(el.tagName)) {
        return null;
    }

    // If element has an ID, try to find a <label for="...">
    const id = el.getAttribute('id');
    if (id) {
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
 * Renders a React component as a portal directly into the specified container (default: document.body).
 * @param component The React element to render.
 * @param container The DOM node to portal into (default: document.body).
 * @returns An object containing the React root.
 */
export function renderReactPortal(component: React.ReactElement, container: HTMLElement = document.body): { reactRoot: Root } {
    // Create a root on the container (body)
    const reactRoot = createRoot(container);
    reactRoot.render(createPortal(component, container));
    return { reactRoot };
}
