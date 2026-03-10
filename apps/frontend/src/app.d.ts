declare global {
	interface Window {
		Square: {
			payments(applicationId: string, locationId: string): any;
		};
	}
}

export {};
