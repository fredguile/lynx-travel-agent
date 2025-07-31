import { createChat } from "@n8n/chat";

import { ENDPOINTS } from "./constants";
import { sessionStorage } from "./sessionStorage";

// const log = createLogger("ui");

// let wrapperId = 0;

export function wrapElementWithAutoSuggest(element: HTMLElement) {
//   log("wrapElementWithAutoSuggest", element);

//   // Get the parent node before detaching
//   const parentNode = element.parentNode;
//   if (!parentNode) {
//     throw new Error(`Element has no parent node: ${element.outerHTML}`);
//   }

//   // Create a placeholder to mark the original position
//   const placeholderEl = document.createElement("div");
//   placeholderEl.id = `ai-auto-suggest-placeholder-${wrapperId}`;
//   parentNode.insertBefore(placeholderEl, element);

//   // Detach element from DOM
//   element.remove();

//   // Render AIAutoSuggest with the element as children
//   renderReactPortal(
//     <AIAutoSuggestContainer scope={`ai-auto-suggest-${wrapperId}`}>
//       <AIAutoSuggest wrapperId={wrapperId}>
//         <HTMLElementWrapper wrapperId={wrapperId} element={element} />
//       </AIAutoSuggest>
//     </AIAutoSuggestContainer>,
//     placeholderEl
//   );

//   wrapperId++;
}

const N8N_CHAT_CSS = `#n8n-chat {
    --chat--color-primary: #329ad0;
    --chat--color-primary-shade-50:rgb(103, 199, 251);
    --chat--color-primary-shade-100:rgb(175, 227, 255);
    --chat--color-secondary: #EF7728;
    --chat--color-secondary-shade-50:rgb(255, 161, 98);
    --chat--color-white: #ffffff;
    --chat--color-light: #f2f4f8;
    --chat--color-light-shade-50: #e6e9f1;
    --chat--color-light-shade-100: #c2c5cc;
    --chat--color-medium: #d2d4d9;
    --chat--color-dark: #333333;
    --chat--color-disabled: #777980;
    --chat--color-typing: #404040;

    --chat--spacing: 1rem;
    --chat--border-radius: 0.25rem;
    --chat--transition-duration: 0.15s;

    --chat--window--width: 400px;
    --chat--window--height: 600px;

    --chat--header-height: auto;
    --chat--header--padding: var(--chat--spacing);
    --chat--header--background: var(--chat--color-dark);
    --chat--header--color: var(--chat--color-light);
    --chat--header--border-top: none;
    --chat--header--border-bottom: none;
    --chat--header--border-bottom: none;
    --chat--header--border-bottom: none;
    --chat--heading--font-size: 2em;
    --chat--header--color: var(--chat--color-light);
    --chat--subtitle--font-size: inherit;
    --chat--subtitle--line-height: 1.8;

    --chat--textarea--height: 50px;

    --chat--message--font-size: 1rem;
    --chat--message--padding: var(--chat--spacing);
    --chat--message--border-radius: var(--chat--border-radius);
    --chat--message-line-height: 1.8;
    --chat--message--bot--background: var(--chat--color-white);
    --chat--message--bot--color: var(--chat--color-dark);
    --chat--message--bot--border: none;
    --chat--message--user--background: var(--chat--color-secondary);
    --chat--message--user--color: var(--chat--color-white);
    --chat--message--user--border: none;
    --chat--message--pre--background: rgba(0, 0, 0, 0.05);

    --chat--toggle--background: var(--chat--color-primary);
    --chat--toggle--hover--background: var(--chat--color-primary-shade-50);
    --chat--toggle--active--background: var(--chat--color-primary-shade-100);
    --chat--toggle--color: var(--chat--color-white);
    --chat--toggle--size: 64px;
}`;

export async function injectN8nChat() {
  const fileNumber = await sessionStorage.currentBookingRef.get();

  // Dynamically import CSS styles
  try {
    await import("@n8n/chat/dist/style.css");
  } catch (error) {
    console.warn("Could not dynamically import @n8n/chat CSS:", error);
  }

  createChat({
    webhookUrl: ENDPOINTS.AI_CHATBOT,
    metadata: {
      fileNumber,
    },
    mode: "window",
    showWelcomeScreen: true,
    initialMessages: [
      "Hi there! 👋",
      fileNumber
        ? `Please ask me anything regarding the file identified by ${fileNumber}.`
        : "Please ask me anything regarding your booking. You must specify the file number for me to indentify the booking.",
    ],
    i18n: {
      en: {
        title: "Lynx Travel Agent",
        subtitle: "",
        footer: "AI Travel Assistant built by @fredguile &",
        getStarted: "New Conversation",
        inputPlaceholder: "Type your question..",
        closeButtonTooltip: "Close chat",
      },
    },
    enableStreaming: false,
  });

  // Inject custom chat CSS variables into the document
  const styleEl = document.createElement("style");
  styleEl.textContent = N8N_CHAT_CSS;
  document.body.appendChild(styleEl);
}
