import { PUBLIC_POCKETBASE_URL } from '$env/static/public';
import PocketBase, { type ClientResponseError } from 'pocketbase';
import { browser } from '$app/environment';

// Configuration
const DEFAULT_URL = 'http://localhost:8090';
const CONNECTION_TIMEOUT = 10000; // 10 seconds

export function createInstance(url?: string): PocketBase {
    const pocketbaseUrl = url || PUBLIC_POCKETBASE_URL || DEFAULT_URL;
    
    const pb = new PocketBase(pocketbaseUrl);
    
    // Set request timeout
    pb.beforeSend = function (url, options) {
        // Add timeout to requests
        const controller = new AbortController();
        const timeoutId = setTimeout(() => controller.abort(), CONNECTION_TIMEOUT);
        
        options.signal = controller.signal;
        
        // Clear timeout if request completes
        const originalSignal = options.signal;
        if (originalSignal) {
            originalSignal.addEventListener('abort', () => clearTimeout(timeoutId));
        }
        
        return { url, options };
    };
    
    return pb;
}

// Global instance
export const pb = createInstance();

// Error handling utilities
export interface AppError {
    message: string;
    code?: string;
    status?: number;
    details?: any;
}

export function handlePocketBaseError(error: unknown): AppError {
    if (error instanceof Error) {
        // Handle ClientResponseError from PocketBase
        if ('status' in error) {
            const pbError = error as ClientResponseError;
            return {
                message: pbError.message || 'An error occurred',
                code: pbError.data?.code,
                status: pbError.status,
                details: pbError.data
            };
        }
        
        // Handle generic errors
        return {
            message: error.message || 'An unexpected error occurred',
            code: 'UNKNOWN_ERROR'
        };
    }
    
    return {
        message: 'An unexpected error occurred',
        code: 'UNKNOWN_ERROR'
    };
}

// Authentication helpers
export const auth = {
    async signIn(email: string, password: string) {
        try {
            return await pb.collection('users').authWithPassword(email, password);
        } catch (error) {
            throw handlePocketBaseError(error);
        }
    },
    
    async signUp(email: string, password: string, passwordConfirm: string, name?: string) {
        try {
            const data = {
                email,
                password,
                passwordConfirm,
                name: name || email.split('@')[0] // Use email prefix as default name
            };
            return await pb.collection('users').create(data);
        } catch (error) {
            throw handlePocketBaseError(error);
        }
    },
    
    signOut() {
        pb.authStore.clear();
    },
    
    isValid() {
        return pb.authStore.isValid;
    },
    
    get user() {
        return pb.authStore.model;
    }
};

// Auto-refresh authentication in browser
if (browser && pb.authStore.isValid) {
    pb.collection('users').authRefresh().catch(() => {
        // Refresh failed, clear invalid auth
        pb.authStore.clear();
    });
}