import type { WhitelistedAIFields } from './types';

export const ENDPOINTS = {
    ANALYSE_USER_CLICK: "http://localhost:5678/webhook/1831ad0f-9c9b-4fb3-99e6-1ce8c0857931",
}

export const WHITELISTED_AI_FIELDS: WhitelistedAIFields = {
    'https://www.lynx-reservations.com/lynx/#_FILE_UPDATE': {
        'TEXTAREA': [0]
    },
    'https://www.lynx-reservations.com/lynx/#_FLIGHT_INFO_PANEL': {
        'INPUT': [231, 239, 240, 241, 242, 243]
    },
    'https://www.lynx-reservations.com/lynx/#_FILE_PAX_PANEL': {
        'INPUT': [179, 180, 181, 186, 187, 188]
    }
};

export const SCREENSHOT_RENDER_DELAY_MS = 32;
