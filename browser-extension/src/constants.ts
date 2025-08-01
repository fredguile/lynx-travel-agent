import type { WhitelistedAIFields } from './types';

export const ENDPOINTS = {
    // ANALYSE_USER_CLICK: "http://localhost:5678/webhook/1831ad0f-9c9b-4fb3-99e6-1ce8c0857931",
    // AI_AUTO_SUGGEST: "http://localhost:5678/webhook/12ad90c0-beb8-44e8-b46e-373d7f4809ea",
    AI_CHATBOT: "https://pn8n.dodmcdund.cc/webhook/022d6461-e68d-4531-b713-833953c388c2/chat",
}

export const WHITELISTED_AI_FIELDS: WhitelistedAIFields = {
    'https://www.lynx-reservations.com/lynx/#_FILE_UPDATE': {
        'TEXTAREA': ['Default Remark:']
    },
    'https://www.lynx-reservations.com/lynx/#_FLIGHT_INFO_PANEL': {
        'INPUT': ['Notes:', 'Flight Number', 'Flight Origin', 'Flight Destination', 'Flight Departure Time', 'Flight Arrival Time'],
        'TEXTAREA': ['Notes:']
    },
    'https://www.lynx-reservations.com/lynx/#_FILE_PAX_PANEL': {
        'INPUT': ['Address:', 'Suburb:', 'Post/Zip Code:', 'Phone No:']
    }
};

export const SCREENSHOT_RENDER_DELAY_MS = 32;
