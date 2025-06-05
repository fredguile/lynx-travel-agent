import { AIAutoSuggest } from './ui/components/AIAutoSuggest';
import { HTMLElementWrapper } from './ui/components/HTMLElementWrapper';
import { createLogger, renderReactPortal } from './utils';

const log = createLogger('ui');

let wrapperId = 0;

export function wrapElementWithAutoSuggest(element: HTMLElement) {
    log('wrapElementWithAutoSuggest', element);

    // Get the parent node before detaching
    const parentNode = element.parentNode;
    if (!parentNode) {
        throw new Error(`Element has no parent node: ${element.outerHTML}`);
    }

    // Create a placeholder to mark the original position
    const placeholderEl = document.createElement('div');
    placeholderEl.id = `ai-auto-suggest-placeholder-${wrapperId}`;
    parentNode.insertBefore(placeholderEl, element);

    // Detach element from DOM
    element.remove();

    // Render AIAutoSuggest with the element as children
    renderReactPortal(
        <AIAutoSuggest wrapperId={wrapperId}>
            <HTMLElementWrapper element={element} />
        </AIAutoSuggest>,
        placeholderEl
    );

    wrapperId++;
}