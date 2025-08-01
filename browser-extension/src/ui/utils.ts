import type { ReactElement } from 'react';
import { createRoot, Root } from 'react-dom/client';
import { createPortal } from 'react-dom';

/**
 * Renders a React component directly into the specified container.
 * @param component The React element to render.
 * @param container The DOM node to render into.
 * @returns An object containing the React root.
 */
export function renderReactComponent(component: ReactElement, container: HTMLElement): { reactRoot: Root } {
    // Create a root on the container
    const reactRoot = createRoot(container);
    reactRoot.render(component);
    return { reactRoot };
}

/**
 * Renders a React component as a portal directly into the specified container (default: document.body).
 * @param component The React element to render.
 * @param container The DOM node to portal into (default: document.body).
 * @returns An object containing the React root.
 */
export function renderReactPortal(component: ReactElement, container: HTMLElement = document.body): { reactRoot: Root } {
    // Create a root on the container (body)
    const reactRoot = createRoot(container);
    reactRoot.render(createPortal(component, container));
    return { reactRoot };
}