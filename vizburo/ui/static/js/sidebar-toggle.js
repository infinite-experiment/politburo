// Sidebar Toggle Script - Minimal vanilla JS (~500 bytes minified)
(function() {
    'use strict';

    const SIDEBAR_STATE_KEY = 'sidebar_visible';
    const SIDEBAR_ID = 'sidebar';
    const TOGGLE_BTN_ID = 'sidebar-toggle-btn';
    const OVERLAY_ID = 'sidebar-overlay';

    // Initialize sidebar state
    function initSidebar() {
        const isVisible = getStoredSidebarState();
        setSidebarVisibility(isVisible);
    }

    // Get sidebar visibility from sessionStorage
    function getStoredSidebarState() {
        const stored = sessionStorage.getItem(SIDEBAR_STATE_KEY);
        return stored === 'true'; // Default to true on first load
    }

    // Set sidebar visibility
    function setSidebarVisibility(visible) {
        const sidebar = document.getElementById(SIDEBAR_ID);
        const overlay = document.getElementById(OVERLAY_ID);

        if (!sidebar) return;

        if (visible) {
            sidebar.classList.remove('hidden');
            if (overlay) overlay.classList.add('hidden');
        } else {
            sidebar.classList.add('hidden');
            if (overlay) overlay.classList.remove('hidden');
        }

        sessionStorage.setItem(SIDEBAR_STATE_KEY, visible ? 'true' : 'false');
    }

    // Toggle sidebar visibility
    function toggleSidebar() {
        const isCurrentlyVisible = !document.getElementById(SIDEBAR_ID).classList.contains('hidden');
        setSidebarVisibility(!isCurrentlyVisible);
    }

    // Setup toggle button
    function setupToggleButton() {
        const toggleBtn = document.getElementById(TOGGLE_BTN_ID);
        const overlay = document.getElementById(OVERLAY_ID);

        if (toggleBtn) {
            toggleBtn.addEventListener('click', function(e) {
                e.preventDefault();
                toggleSidebar();
            });
        }

        // Close sidebar when clicking overlay
        if (overlay) {
            overlay.addEventListener('click', function(e) {
                e.preventDefault();
                setSidebarVisibility(false);
            });
        }
    }

    // Handle window resize (show sidebar on desktop, hide on mobile)
    function handleResize() {
        const mediaQuery = window.matchMedia('(min-width: 768px)'); // md breakpoint

        function handleMediaChange(e) {
            if (e.matches) {
                // Desktop: show sidebar
                document.getElementById(SIDEBAR_ID).classList.remove('hidden');
                const overlay = document.getElementById(OVERLAY_ID);
                if (overlay) overlay.classList.add('hidden');
            } else {
                // Mobile: hide sidebar
                const isVisible = getStoredSidebarState();
                setSidebarVisibility(isVisible);
            }
        }

        mediaQuery.addEventListener('change', handleMediaChange);
    }

    // Initialize on DOM ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', function() {
            initSidebar();
            setupToggleButton();
            handleResize();
        });
    } else {
        initSidebar();
        setupToggleButton();
        handleResize();
    }
})();
