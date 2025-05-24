/// <reference types="@types/firefox-webext-browser"/>
import { fromEvent, filter, switchMap, tap } from 'rxjs';
import { timer, Observable } from 'rxjs';
import { createLogger, findElementLabel, isWhitelistedAIField, base64ImageToBlob } from './utils';
import { createRedCrossElement } from './ui';
import { ENDPOINTS, SCREENSHOT_RENDER_DELAY_MS } from './constants';

const log = createLogger('contentScript');

log('Content script loaded');

// --- Page Change Detection and Form Extraction ---
let currentUrl = location.href;

function onPageChange() {
    currentUrl = location.href;
    log('Page changed:', currentUrl);
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

// Listen for clicks on the document
fromEvent<MouseEvent>(window, 'click').pipe(
    // tap((event: MouseEvent) => log("Element label:", findElementLabel(event.target as HTMLElement))),
    filter((event: MouseEvent) => isWhitelistedAIField(currentUrl, event)),
    switchMap((event: MouseEvent) => {
        // Mark the click location with a red cross
        const cross = createRedCrossElement(event);
        document.body.appendChild(cross);

        // Wait a short moment to ensure the cross is rendered
        return timer(SCREENSHOT_RENDER_DELAY_MS).pipe(
            switchMap(() => new Observable<void>((subscriber) => {
                const controller = new AbortController();
                browser.runtime.sendMessage({ action: 'capture_screenshot' })
                    .then((response) => {
                        // Remove the cross after screenshot
                        if (cross.parentNode) {
                            cross.parentNode.removeChild(cross);
                        }

                        if (!response || !response.screenshot) {
                            subscriber.complete();
                            return;
                        }

                        // Convert base64 to binary Blob
                        const blob = base64ImageToBlob(response.screenshot);

                        // Prepare multipart/form-data
                        const formData = new FormData();
                        formData.append('screenshot', blob, 'screenshot.png');
                        log('Sending screenshot to:', ENDPOINTS.ANALYSE_USER_CLICK);
                        fetch(`${ENDPOINTS.ANALYSE_USER_CLICK}?currentUrl=${encodeURIComponent(currentUrl)}`, {
                            method: 'POST',
                            body: formData,
                            signal: controller.signal,
                        })
                            .then(res => res.text())
                            .then(json => {
                                log('Server response:', json);
                                subscriber.next();
                                subscriber.complete();
                            })
                            .catch(err => {
                                if (err.name === 'AbortError') {
                                    log('Fetch aborted');
                                } else {
                                    log('Error posting screenshot:', err);
                                }
                                subscriber.error(err);
                            });

                    })
                    .catch((err) => {
                        if (cross.parentNode) {
                            cross.parentNode.removeChild(cross);
                        }
                        log('Failed to capture screenshot:', err);
                        subscriber.error(err);
                    });
                // Teardown logic: abort fetch if unsubscribed
                return () => {
                    if (cross.parentNode) {
                        cross.parentNode.removeChild(cross);
                    }
                    controller.abort();
                };
            }))
        );
    })
).subscribe();

// Initial extraction
onPageChange();
