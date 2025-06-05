/// <reference types="@types/firefox-webext-browser"/>

interface StorageData {
  currentUrl?: string;
  currentBookingRef?: string;
}

interface StorageProperty<T> {
  get(): Promise<T | null>;
  set(value: T | null): Promise<void>;
  onPropertyChange(callback: (value: T | null) => void): void;
}

export const sessionStorage = {
  currentUrl: {
    get: () => createStorageGetter('currentUrl'),
    set: (value: string | null) => createStorageSetter('currentUrl', value),
    onPropertyChange: (callback: (url: string | null) => void) => createStorageListener('currentUrl', callback)
  } as StorageProperty<string>,

  currentBookingRef: {
    get: () => createStorageGetter('currentBookingRef'),
    set: (value: string | null) => createStorageSetter('currentBookingRef', value),
    onPropertyChange: (callback: (bookingRef: string | null) => void) => createStorageListener('currentBookingRef', callback)
  } as StorageProperty<string>
};

// Utility functions
async function createStorageGetter(key: string): Promise<string | null> {
  try {
    const result = await browser.storage.local.get(key);
    return result[key] || null;
  } catch (error) {
    console.error(`Error getting ${key} from storage:`, error);
    return null;
  }
}

async function createStorageSetter(key: string, value: string | null): Promise<void> {
  try {
    if (value) {
      await browser.storage.local.set({ [key]: value });
    } else {
      await browser.storage.local.remove(key);
    }
  } catch (error) {
    console.error(`Error setting ${key} in storage:`, error);
  }
}

function createStorageListener(key: string, callback: (value: string | null) => void): void {
  browser.storage.local.onChanged.addListener((changes) => {
    if (changes[key]) {
      const newValue = changes[key].newValue || null;
      callback(newValue);
    }
  });
} 