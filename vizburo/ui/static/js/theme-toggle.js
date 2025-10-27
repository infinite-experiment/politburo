// Theme Toggle Script - Minimal vanilla JS (~500 bytes minified)
(function() {
    'use strict';

    const THEME_COOKIE = 'theme_preference';
    const VALID_THEMES = ['light', 'dark', 'high-contrast', 'obsidian'];
    const DEFAULT_THEME = 'light';

    // Load theme from cookie or localStorage on page load
    function initTheme() {
        let theme = getThemeFromStorage();
        if (!VALID_THEMES.includes(theme)) {
            theme = DEFAULT_THEME;
        }
        applyTheme(theme);
    }

    // Get theme from cookie (server-set) or localStorage (client-set)
    function getThemeFromStorage() {
        // Check server-set cookie first
        const cookieTheme = getCookie(THEME_COOKIE);
        if (cookieTheme && VALID_THEMES.includes(cookieTheme)) {
            return cookieTheme;
        }

        // Fall back to localStorage
        return localStorage.getItem('theme_preference') || DEFAULT_THEME;
    }

    // Get cookie value by name
    function getCookie(name) {
        const value = `; ${document.cookie}`;
        const parts = value.split(`; ${name}=`);
        if (parts.length === 2) return parts.pop().split(';').shift();
        return null;
    }

    // Apply theme by setting data-theme attribute
    function applyTheme(theme) {
        if (!VALID_THEMES.includes(theme)) {
            theme = DEFAULT_THEME;
        }

        document.documentElement.setAttribute('data-theme', theme);
        localStorage.setItem('theme_preference', theme);

        // Send to server to set cookie
        setThemeCookie(theme);

        // Update active button styling
        updateThemeButtons(theme);
    }

    // Set theme cookie on server
    function setThemeCookie(theme) {
        fetch('/ui/api/theme', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: `theme=${encodeURIComponent(theme)}`,
        }).catch(err => console.warn('Failed to set theme cookie:', err));
    }

    // Update button styling to show active theme
    function updateThemeButtons(theme) {
        document.querySelectorAll('.theme-btn').forEach(btn => {
            const btnTheme = btn.getAttribute('data-theme');
            if (btnTheme === theme) {
                btn.classList.add('bg-white', 'dark:bg-gray-600', 'text-gray-900', 'dark:text-white', 'font-bold');
                btn.classList.remove('text-gray-600', 'dark:text-gray-400');
            } else {
                btn.classList.remove('bg-white', 'dark:bg-gray-600', 'text-gray-900', 'dark:text-white', 'font-bold');
                btn.classList.add('text-gray-600', 'dark:text-gray-400');
            }
        });
    }

    // Handle theme button clicks
    function setupThemeButtons() {
        document.querySelectorAll('.theme-btn').forEach(btn => {
            btn.addEventListener('click', function(e) {
                e.preventDefault();
                const theme = this.getAttribute('data-theme');
                applyTheme(theme);
            });
        });
    }

    // Initialize on DOM ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', function() {
            initTheme();
            setupThemeButtons();
        });
    } else {
        initTheme();
        setupThemeButtons();
    }
})();
