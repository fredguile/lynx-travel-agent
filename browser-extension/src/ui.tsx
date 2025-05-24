import { AIAutosuggestOverlay } from './ui/components/AIAutosuggestOverlay';

import { createLogger, renderReactComponent } from './utils';

const log = createLogger('ui');

export function wrapElementWithAutosuggest(element: HTMLElement) {
    const rect = element.getBoundingClientRect();
    const left = rect.left + rect.width / 2 + window.scrollX;
    const top = rect.top + window.scrollY;

    log('wrapElementWithAutosuggest', { left, top })

    let handler: ReturnType<typeof renderReactComponent> | null = null;

    const onClose = () => {
        if (handler?.reactRoot) {
            handler.reactRoot.unmount();
            handler.el.remove();
        }
    };

    handler = renderReactComponent(
        <AIAutosuggestOverlay left={left} top={top} onClose={onClose} />
    );
    element.parentElement?.appendChild(handler.el);

    return onClose;
}