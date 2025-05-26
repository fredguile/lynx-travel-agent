import { createStore, defaults, StoreActionApi, createHook, createContainer } from 'react-sweet-state';

import { ENDPOINTS } from '../../constants';
import { base64ImageToBlob, createLogger } from '../../utils';

if (process.env.NODE_ENV === 'development') {
  defaults.devtools = true;
}

const log = createLogger('aiAutoSuggestStore');

interface State {
  visible: boolean;
  crossVisible: boolean;
  crossLeft: number;
  loading: boolean;
  error: string | null;
}

const initialState: State = {
  visible: false,
  crossVisible: false,
  crossLeft: 0,
  loading: false,
  error: null,
};

interface SuggestContentParams {
  currentUrl: string;
  currentBookingRef: string;
  onSuccess?: (response: string) => void;
  onError?: (error: Error) => void;
}

const actions = {
  setVisible: (visible: boolean) => ({ setState }: StoreActionApi<State>) => setState({ visible }),
  setCrossLeft: (crossLeft: number) => ({ setState }: StoreActionApi<State>) => setState({ crossLeft }),
  aiSuggestContent:
    ({ currentUrl, currentBookingRef, onSuccess, onError }: SuggestContentParams) =>
      async ({ setState }: StoreActionApi<State>) => {
        setState({ loading: true, crossVisible: true, error: null });

        try {
          log('taking screenshot', currentUrl);

          const response = await browser.runtime.sendMessage({ action: 'capture_screenshot' });
          const blob = base64ImageToBlob(response.screenshot);

          setState({ crossVisible: false });

          log('analysing screen context', currentUrl);

          const formData = new FormData();
          formData.append('screenshot', blob, 'screenshot.png');
          let res = await fetch(`${ENDPOINTS.ANALYSE_USER_CLICK}?currentUrl=${encodeURIComponent(currentUrl)}`, {
            method: 'POST',
            body: formData,
          });
          const screenContext = await res.text();

          log('requesting ai auto suggest', currentUrl);

          res = await fetch(`${ENDPOINTS.AI_AUTO_SUGGEST}`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              currentUrl,
              currentBookingRef,
              screenContext,
            }),
          });
          const aiSuggestion = await res.text();

          setState({ loading: false, error: null });
          onSuccess?.(aiSuggestion);
        } catch (err: any) {
          setState({ loading: false, error: err.message || 'Unknown error' });
          onError?.(err);
        } finally {
          setState({ crossVisible: false });
        }
      },
};

export const AIAutoSuggestStore = createStore({
  initialState,
  actions,
  name: 'AIAutoSuggestStore',
});

export const useAIAutoSuggestStore = createHook(AIAutoSuggestStore);

export const AIAutoSuggestContainer = createContainer();
