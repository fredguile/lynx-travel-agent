/** @jsxImportSource @compiled/react */
import React, { useEffect, useRef } from 'react';

import iconDark from '../../icons/icon-dark-128.png';

interface AIAutosuggestOverlayProps {
  left: number;
  top: number;
  onClose: () => void;
}

const HIDE_DELAY_MS = 800;

export const AIAutosuggestOverlay: React.FC<AIAutosuggestOverlayProps> = ({ left, top, onClose }) => {
  const overlayRef = useRef<HTMLDivElement>(null);
  const hideTimer = useRef<NodeJS.Timeout | null>(null);

  // Start hide timer on mouse leave
  const startHideTimer = () => {
    hideTimer.current = setTimeout(() => {
      onClose();
    }, HIDE_DELAY_MS);
  };

  // Cancel hide timer on mouse enter
  const cancelHideTimer = () => {
    if (hideTimer.current) {
      clearTimeout(hideTimer.current);
      hideTimer.current = null;
    }
  };

  useEffect(() => {
    const overlay = overlayRef.current;
    if (!overlay) return;
    overlay.addEventListener('mouseleave', startHideTimer);
    overlay.addEventListener('mouseenter', cancelHideTimer);
    return () => {
      overlay.removeEventListener('mouseleave', startHideTimer);
      overlay.removeEventListener('mouseenter', cancelHideTimer);
      cancelHideTimer();
    };
  }, []);

  return (
    <div
      ref={overlayRef}
      css={{
        position: 'absolute',
        background: 'white',
        borderRadius: '12px',
        padding: '6px 14px',
        boxShadow: '0 2px 8px rgba(0,0,0,0.15)',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'flex-start',
        zIndex: 1000000,
        pointerEvents: 'auto',
        width: 'auto',
        minWidth: 0,
        height: 'auto',
        left: left,
        top: top,
        transform: 'translate(-50%, -100%)',
      }}
    >
      <a
        href="#"
        css={{
          display: 'flex',
          alignItems: 'center',
          fontWeight: 'bold',
          textDecoration: 'none',
          color: '#007bff',
        }}
      >
        <img
          src={`data:image/png;base64,${iconDark}`}
          alt="AI"
          css={{
            width: 20,
            height: 20,
            marginRight: 8,
            verticalAlign: 'middle',
          }}
        />
        Suggest content
      </a>
    </div>
  );
}; 