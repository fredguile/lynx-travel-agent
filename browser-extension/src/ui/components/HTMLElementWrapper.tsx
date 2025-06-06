/** @jsxImportSource @compiled/react */
import { useEffect, useRef } from "react";
import { css } from '@compiled/react';

import { useAIAutoSuggestStore } from "../state/AIAutoSuggestStore";
import { createLogger } from "../../utils";

const log = createLogger('HTMLElementWrapper');

export function HTMLElementWrapper({ wrapperId, element }: { wrapperId: number, element: HTMLElement }) {
    const containerRef = useRef<HTMLDivElement>(null);

    const [
        { highlight },
    ] = useAIAutoSuggestStore();

    useEffect(() => {
        if (containerRef.current && element) {
            containerRef.current.appendChild(element);
        }
        // Cleanup: remove the element when unmounting
        return () => {
            if (containerRef.current && element && containerRef.current.contains(element)) {
                containerRef.current.removeChild(element);
            }
        };
    }, [element]);

    return <div
        id={`ai-auto-suggest-element-wrapper-${wrapperId}`}
        css={wrapperStyle}
        style={{ border: highlight ? '3px solid red' : 'none' }}
        ref={containerRef}
    />;
}

const wrapperStyle = css({
    position: 'relative',
});
