//@ts-nocheck
console.log("sw: service worker loaded");

self.addEventListener('install', () => self.skipWaiting()); // Remove service worker for BookBrowser