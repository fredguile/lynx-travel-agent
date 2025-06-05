import type { WhitelistedAIFields } from './types';

export const ENDPOINTS = {
    ANALYSE_USER_CLICK: "http://localhost:5678/webhook/1831ad0f-9c9b-4fb3-99e6-1ce8c0857931",
    AI_AUTO_SUGGEST: "http://localhost:5678/webhook/835185f1-a443-49c1-8821-892fdee51424",
}

export const WHITELISTED_AI_FIELDS: WhitelistedAIFields = {
    'https://www.lynx-reservations.com/lynx/#_FILE_UPDATE': {
        'TEXTAREA': ['Default Remark:']
    },
    'https://www.lynx-reservations.com/lynx/#_FLIGHT_INFO_PANEL': {
        'INPUT': ['Notes:', 'Flight Number', 'Flight Origin', '	Flight Destination', '	Flight Departure Time', 'Flight Arrival Time']
    },
    'https://www.lynx-reservations.com/lynx/#_FILE_PAX_PANEL': {
        'INPUT': ['Address:', 'Suburb:', 'Post/Zip Code:', 'Phone No:']
    }
};

export const SCREENSHOT_RENDER_DELAY_MS = 32;
