/// <reference types="@types/firefox-webext-browser"/>
import { fromEvent, filter, switchMap } from 'rxjs';

console.log('Content script loaded');

// --- Global Variables ---
const ANALYSE_USER_CLICK_ENDPOINTS = "http://localhost:5678/webhook/1831ad0f-9c9b-4fb3-99e6-1ce8c0857931";

// --- Page Change Detection and Form Extraction ---
let currentUrl = location.href;
let lastScreenshotBase64: string | null = null;

function onPageChange() {
    currentUrl = location.href;
    console.log('Page changed:', currentUrl);
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

function createRedCrossElement(event: MouseEvent) {
    // Draw a red cross at the cursor position
    const crossSize = 20;
    const cross = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
    cross.setAttribute('width', `${crossSize}`);
    cross.setAttribute('height', `${crossSize}`);
    cross.style.position = 'absolute';
    cross.style.left = `${event.pageX - crossSize / 2}px`;
    cross.style.top = `${event.pageY - crossSize / 2}px`;
    cross.style.pointerEvents = 'none';
    cross.style.zIndex = '999999';
    cross.innerHTML = `
        <line x1="0" y1="0" x2="${crossSize}" y2="${crossSize}" stroke="red" stroke-width="3" />
        <line x1="${crossSize}" y1="0" x2="0" y2="${crossSize}" stroke="red" stroke-width="3" />
    `;

    return cross;
}

// Listen for clicks on the document
fromEvent<MouseEvent>(window, 'click').pipe(
    filter((event: MouseEvent) => {
        const el = event.target as HTMLElement | null;
        return !!el && (el.tagName === "INPUT" || el.tagName === "TEXTAREA" || el.tagName === "SELECT");
    }),
    switchMap((event: MouseEvent) => {
        // Mark the click location with a red cross
        const cross = createRedCrossElement(event);
        document.body.appendChild(cross);

        // Wait a short moment to ensure the cross is rendered
        return new Promise<void>((resolve) => setTimeout(resolve, 32)).then(() => {
            return browser.runtime.sendMessage({ action: 'capture_screenshot' })
                .then((response) => {
                    // Remove the cross after screenshot
                    if (cross.parentNode) {
                        cross.parentNode.removeChild(cross);
                    }

                    if (response && response.screenshot) {
                        lastScreenshotBase64 = response.screenshot;

                        // Convert base64 to binary Blob
                        let byteString: string;
                        if (lastScreenshotBase64?.startsWith('data:image/')) {
                            byteString = atob(lastScreenshotBase64.split(',')[1]);
                        } else {
                            byteString = atob(lastScreenshotBase64 || '');
                        }
                        const ab = new ArrayBuffer(byteString.length);
                        const ia = new Uint8Array(ab);
                        for (let i = 0; i < byteString.length; i++) {
                            ia[i] = byteString.charCodeAt(i);
                        }
                        const blob = new Blob([ab], { type: 'image/png' });

                        // Prepare multipart/form-data
                        const formData = new FormData();
                        formData.append('screenshot', blob, 'screenshot.png');

                        // Post the screenshot as binary (multipart/form-data) with currentUrl as query string
                        const urlWithQuery = `${ANALYSE_USER_CLICK_ENDPOINTS}?currentUrl=${encodeURIComponent(currentUrl)}`;
                        return fetch(urlWithQuery, {
                            method: 'POST',
                            body: formData,
                        })
                            .then(res => res.text())
                            .then(json => {
                                console.log('Server response:', json);
                            })
                            .catch(err => {
                                console.error('Error posting screenshot:', err);
                            });
                    }
                })
                .catch((err) => {
                    console.error('Failed to capture screenshot:', err);
                    // Remove the cross if error
                    if (cross.parentNode) {
                        cross.parentNode.removeChild(cross);
                    }
                });
        });
    })
).subscribe();

// Initial extraction
onPageChange();
