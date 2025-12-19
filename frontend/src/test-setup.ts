import '@testing-library/jest-dom';

// Mock localStorage for persisted stores
const localStorageMock = (() => {
    let store: Record<string, string> = {};
    return {
        getItem: (key: string) => store[key] ?? null,
        setItem: (key: string, value: string) => { store[key] = value; },
        removeItem: (key: string) => { delete store[key]; },
        clear: () => { store = {}; },
        length: 0,
        key: () => null,
    };
})();

Object.defineProperty(globalThis, 'localStorage', {
    value: localStorageMock,
});
