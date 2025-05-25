/** @jsxImportSource @compiled/react */
import { ReactElement, useCallback, useEffect, useRef, useState } from 'react';
import { css } from '@compiled/react';

import { RedCross } from './RedCross';
import icon from 'url:../../icons/icon-light-32.png';
import { base64ImageToBlob, createLogger } from '../../utils';
import { ENDPOINTS } from '../../constants';

interface AIAutoSuggestProps {
  wrapperId: number;
  currentUrl: string;
  children: ReactElement;
}

const HIDE_TIMER_MS = 800;
const log = createLogger('aiAutoSuggest');

export const AIAutoSuggest = ({ wrapperId, currentUrl, children }: AIAutoSuggestProps) => {
  const [visible, setVisible] = useState(false);
  const childrenRef = useRef<HTMLDivElement>(null);
  const [redCrossVisible, setRedCrossVisible] = useState(false);
  const [redCrossLeft, setRedCrossLeft] = useState(0);
  const hideTimer = useRef<NodeJS.Timeout | null>(null);

  const onClick = useCallback(async () => {
    setRedCrossVisible(true);

    const controller = new AbortController();

    try {
      // Capture screenshot
      const response = await browser.runtime.sendMessage({ action: 'capture_screenshot' });

      // Convert base64 to binary Blob
      const blob = base64ImageToBlob(response.screenshot);

      // Prepare multipart/form-data
      const formData = new FormData();
      formData.append('screenshot', blob, 'screenshot.png');
      log('Sending screenshot to:', ENDPOINTS.ANALYSE_USER_CLICK);
      const res = await fetch(`${ENDPOINTS.ANALYSE_USER_CLICK}?currentUrl=${encodeURIComponent(currentUrl)}`, {
        method: 'POST',
        body: formData,
        signal: controller.signal,
      });

      const text = await res.text();
      log('Server response:', text);
    } catch (err: any) {
      if (err.name === 'AbortError') {
        log('Fetch aborted');
      } else {
        log('Error posting screenshot:', err);
      }
    } finally {
      setRedCrossVisible(false);
    }
  }, []);

  const handleMouseEnter = () => {
    if (hideTimer.current) {
      clearTimeout(hideTimer.current);
      hideTimer.current = null;
    }
    setVisible(true);
  };

  const handleMouseLeave = () => {
    hideTimer.current = setTimeout(() => {
      setVisible(false);
    }, HIDE_TIMER_MS);
  };

  useEffect(() => {
    const childrenEl = childrenRef.current;
    if (childrenEl) {
      const rect = childrenEl.getBoundingClientRect();
      setRedCrossLeft(rect.width + 20);
    }
  }, []);

  return (
    <>
      <div
        id={`ai-auto-suggest-${wrapperId}`}
        css={aiAutoSuggestStyle}
        style={{ display: visible ? 'flex' : 'none' }}
        onMouseEnter={handleMouseEnter}
        onMouseLeave={handleMouseLeave}
      >
        <img
          src={icon}
          css={iconStyle}
          alt="Pan PAC AI Helper"
        />
        <a css={linkStyle} onClick={onClick}>Suggest content</a>
      </div>
      <div
        onMouseEnter={handleMouseEnter}
        onMouseLeave={handleMouseLeave}
        css={childrenStyle}
        ref={childrenRef}
      >
        {children}
      </div>
      {redCrossVisible && <RedCross left={redCrossLeft} />}
    </>
  );
};

const aiAutoSuggestStyle = css({
  position: 'absolute',
  top: '5px',
  flexDirection: 'row',
  justifyContent: 'center',
  alignItems: 'center',
  width: 'auto',
  minWidth: 0,
  maxWidth: 200,
  padding: '4px 4px',
  border: '1px solid #7F9DB9',
  borderRadius: '6px',
  background: '#FFFFFF',
  boxShadow: '0 4px 16px rgba(0,0,0,0.15)',
  pointerEvents: 'auto',
  zIndex: 1000000,
});

const iconStyle = css({
  width: 16,
  height: 16,
  marginRight: 4,
  flexShrink: 0,
});

const linkStyle = css({
  display: 'flex',
  flexDirection: 'column',
  fontSize: '1.0em',
  fontWeight: 'bold',
  lineHeight: '12px',
  textDecoration: 'none !important',
  color: '#333333',
  ':hover': {
    color: '#333333',
    textDecoration: 'underline !important',
  },
  ':visited': {
    color: '#333333',
    textDecoration: 'none !important',
  },
});

const childrenStyle = css({
  display: 'inline-block',
  position: 'relative',
});