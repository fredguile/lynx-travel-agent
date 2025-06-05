/** @jsxImportSource @compiled/react */
import { ReactElement, useCallback, useEffect, useRef } from 'react';
import { css } from '@compiled/react';

import { RedCross } from './RedCross';
import { useAIAutoSuggestStore, AIAutoSuggestContainer } from '../state/AIAutoSuggestStore';
import icon from 'url:../../icons/icon-light-32.png';
import { createLogger } from '../../utils';
import { sessionStorage } from '../../sessionStorage';

interface AIAutoSuggestProps {
  wrapperId: number;
  children: ReactElement;
}

const HIDE_TIMER_MS = 800;
const CROSS_LEFT_OFFSET = 20;

const log = createLogger('aiAutoSuggest');

export const AIAutoSuggest = ({ wrapperId, children }: AIAutoSuggestProps) => {
  const childrenRef = useRef<HTMLDivElement>(null);
  const hideTimer = useRef<NodeJS.Timeout | null>(null);

  const [
    { visible, crossVisible, crossLeft, loading },
    { setVisible, setCrossLeft, aiSuggestContent }
  ] = useAIAutoSuggestStore();

  useEffect(() => {
    // position cross relative to the children
    const childrenEl = childrenRef.current;
    if (childrenEl) {
      const rect = childrenEl.getBoundingClientRect();
      setCrossLeft(rect.width + CROSS_LEFT_OFFSET);
    }
  }, []);

  const handleMouseEnter = useCallback(() => {
    if (hideTimer.current) {
      clearTimeout(hideTimer.current);
      hideTimer.current = null;
    }
    setVisible(true);
  }, [setVisible]);

  const handleMouseLeave = useCallback(() => {
    hideTimer.current = setTimeout(() => {
      if (!loading) {
        setVisible(false);
      }
    }, HIDE_TIMER_MS);
  }, [setVisible, loading]);

  const onClick = useCallback(() => {
    (async () => {
      const currentUrl = await sessionStorage.currentUrl.get();
      const currentBookingRef = await sessionStorage.currentBookingRef.get();
      aiSuggestContent({
        currentUrl: currentUrl || '',
        currentBookingRef: currentBookingRef || '',
        onSuccess: (text: string) => {
          log('Server response:', text);
        },
        onError: (err: Error) => {
          log('Error posting screenshot:', err);
        },
      });
    })();
  }, [aiSuggestContent]);

  return (
    <AIAutoSuggestContainer scope={`ai-auto-suggest-${wrapperId}`}>
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
        <a css={linkStyle} onClick={onClick}>
          {loading ? 'Loading...' : 'Suggest content'}
        </a>
      </div>
      <div
        onMouseEnter={handleMouseEnter}
        onMouseLeave={handleMouseLeave}
        css={childrenStyle}
        ref={childrenRef}
      >
        {children}
      </div>
      {crossVisible && <RedCross left={crossLeft} />}
    </AIAutoSuggestContainer>
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