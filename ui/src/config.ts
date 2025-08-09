import { browser } from '$app/environment';

// Application configuration with environment variable support
export const config = {
    website: {
        name: "Photo Cifu",
        baseUrl: browser ? window.location.origin : "https://home-server.tailadc076.ts.net",
        description: "Photo gallery application with workflow-based image processing. Built with SvelteKit, PocketBase, and Go.",
    },
    features: {
        createProfileStep: true,
        maxGallerySize: 100, // Maximum number of images per gallery
        supportedImageFormats: ['jpg', 'jpeg', 'png', 'gif', 'webp'],
        maxFileSize: 100 * 1024 * 1024, // 100MB in bytes
    },
    ui: {
        defaultPageSize: 20,
        thumbnailSize: '100x100',
        loadingTimeout: 30000, // 30 seconds
    }
} as const;

// Legacy exports for backward compatibility
export const WebsiteName = config.website.name;
export const WebsiteBaseUrl = config.website.baseUrl;
export const WebsiteDescription = config.website.description;
export const CreateProfileStep = config.features.createProfileStep;
