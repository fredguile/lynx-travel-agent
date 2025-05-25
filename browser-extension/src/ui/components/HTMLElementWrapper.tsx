/** @jsxImportSource @compiled/react */
import { useEffect, useRef } from "react";

export function HTMLElementWrapper({ element }: { element: HTMLElement }) {
    const containerRef = useRef<HTMLDivElement>(null);

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

    return <div css={{ position: 'relative' }} ref={containerRef} />;
}