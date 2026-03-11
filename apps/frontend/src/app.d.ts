declare module 'virtual:pwa-register' {
    export function registerSW(options?: { immediate?: boolean }): void;
}

interface Window {
    Square: {
        payments(applicationId: string, locationId: string): any;
    };
}
