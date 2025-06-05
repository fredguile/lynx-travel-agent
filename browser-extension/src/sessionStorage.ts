/// <reference types="@types/firefox-webext-browser"/>

interface StorageData {
  currentBookingRef?: string;
}

export const sessionStorage = {
  // Get the current booking reference
  async getCurrentBookingRef(): Promise<string | null> {
    try {
      const result = await browser.storage.local.get('currentBookingRef');
      return result.currentBookingRef || null;
    } catch (error) {
      console.error('Error getting currentBookingRef from storage:', error);
      return null;
    }
  },

  // Set the current booking reference
  async setCurrentBookingRef(bookingRef: string | null): Promise<void> {
    try {
      if (bookingRef) {
        await browser.storage.local.set({ currentBookingRef: bookingRef });
      } else {
        await browser.storage.local.remove('currentBookingRef');
      }
    } catch (error) {
      console.error('Error setting currentBookingRef in storage:', error);
    }
  },

  // Listen for changes to the booking reference
  onBookingRefChange(callback: (bookingRef: string | null) => void): void {
    browser.storage.local.onChanged.addListener((changes) => {
      if (changes.currentBookingRef) {
        const newValue = changes.currentBookingRef.newValue || null;
        callback(newValue);
      }
    });
  }
}; 