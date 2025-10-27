package ui

import (
	"net/http"
)

// DashboardView handles the main dashboard view with sidebar
func (h *UIHandler) DashboardView(w http.ResponseWriter, r *http.Request) {
	h.DashboardHandler(w, r)
}

// MinimalDashboardView handles the minimal dashboard view without sidebar
func (h *UIHandler) MinimalDashboardView(w http.ResponseWriter, r *http.Request) {
	h.MinimalDashboardHandler(w, r)
}

// ShowcaseHandler renders the component showcase page
func (h *UIHandler) ShowcaseHandler(w http.ResponseWriter, r *http.Request) {
	showcaseContent := `
<div>
    <!-- Header -->
    <div class="mb-8">
        <p class="text-gray-600 dark:text-gray-400">Interactive UI component library with HTMX integration</p>
    </div>

    <!-- Navigation Tabs -->
    <div class="mb-8 flex flex-wrap gap-1 border-b border-gray-200 dark:border-gray-700 pb-0">
        <button hx-get="/dashboard/showcase/component/buttons"
                hx-target="#component-content"
                hx-swap="innerHTML"
                class="px-4 py-3 font-semibold text-flight-primary border-b-2 border-flight-primary hover:opacity-80 transition-opacity">
            Buttons
        </button>
        <button hx-get="/dashboard/showcase/component/forms"
                hx-target="#component-content"
                hx-swap="innerHTML"
                class="px-4 py-3 font-semibold text-gray-600 dark:text-gray-400 border-b-2 border-transparent hover:text-flight-primary transition-colors">
            Forms & Inputs
        </button>
        <button hx-get="/dashboard/showcase/component/typography"
                hx-target="#component-content"
                hx-swap="innerHTML"
                class="px-4 py-3 font-semibold text-gray-600 dark:text-gray-400 border-b-2 border-transparent hover:text-flight-primary transition-colors">
            Typography
        </button>
        <button hx-get="/dashboard/showcase/component/cards"
                hx-target="#component-content"
                hx-swap="innerHTML"
                class="px-4 py-3 font-semibold text-gray-600 dark:text-gray-400 border-b-2 border-transparent hover:text-flight-primary transition-colors">
            Cards & Tables
        </button>
        <button hx-get="/dashboard/showcase/component/badges"
                hx-target="#component-content"
                hx-swap="innerHTML"
                class="px-4 py-3 font-semibold text-gray-600 dark:text-gray-400 border-b-2 border-transparent hover:text-flight-primary transition-colors">
            Badges & Chips
        </button>
        <button hx-get="/dashboard/showcase/component/alerts"
                hx-target="#component-content"
                hx-swap="innerHTML"
                class="px-4 py-3 font-semibold text-gray-600 dark:text-gray-400 border-b-2 border-transparent hover:text-flight-primary transition-colors">
            Alerts
        </button>
    </div>

    <!-- Component Content Area -->
    <div id="component-content" class="bg-gray-50 dark:bg-gray-800 rounded-lg p-8 min-h-96">
        <!-- Default: Buttons Component -->
        <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">Buttons</h2>

        <div class="space-y-6">
            <!-- Primary Button Variants -->
            <div>
                <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-3">Primary Buttons</h3>
                <div class="flex flex-wrap gap-3">
                    <button class="px-6 py-2 bg-flight-primary text-white rounded-lg font-semibold hover:bg-opacity-90 transition-all">
                        Primary Button
                    </button>
                    <button class="px-6 py-2 bg-flight-primary text-white rounded-lg font-semibold hover:bg-opacity-90 transition-all opacity-50 cursor-not-allowed">
                        Disabled
                    </button>
                    <button class="px-6 py-2 bg-flight-primary text-white rounded-lg font-semibold hover:bg-opacity-90 transition-all">
                        <span class="inline-block mr-2">✈️</span>
                        With Icon
                    </button>
                </div>
            </div>

            <!-- Secondary Button Variants -->
            <div>
                <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-3">Secondary Buttons</h3>
                <div class="flex flex-wrap gap-3">
                    <button class="px-6 py-2 border-2 border-flight-primary text-flight-primary rounded-lg font-semibold hover:bg-flight-primary hover:text-white transition-all">
                        Secondary
                    </button>
                    <button class="px-6 py-2 border-2 border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-lg font-semibold hover:bg-gray-100 dark:hover:bg-gray-700 transition-all">
                        Neutral
                    </button>
                    <button class="px-6 py-2 border-2 border-red-500 text-red-500 rounded-lg font-semibold hover:bg-red-500 hover:text-white transition-all">
                        Danger
                    </button>
                </div>
            </div>

            <!-- Button Sizes -->
            <div>
                <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-3">Button Sizes</h3>
                <div class="flex flex-wrap gap-3 items-center">
                    <button class="px-3 py-1 text-sm bg-flight-primary text-white rounded font-semibold">Small</button>
                    <button class="px-4 py-2 bg-flight-primary text-white rounded-lg font-semibold">Medium</button>
                    <button class="px-6 py-3 text-lg bg-flight-primary text-white rounded-lg font-semibold">Large</button>
                </div>
            </div>

            <!-- Button States -->
            <div>
                <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-3">Button States</h3>
                <div class="flex flex-wrap gap-3">
                    <button class="px-6 py-2 bg-flight-primary text-white rounded-lg font-semibold hover:shadow-lg transition-all">Hover State</button>
                    <button class="px-6 py-2 bg-flight-primary text-white rounded-lg font-semibold active:scale-95 transition-all">Active State</button>
                    <button class="px-6 py-2 bg-gray-300 text-gray-600 rounded-lg font-semibold cursor-not-allowed">Disabled State</button>
                </div>
            </div>
        </div>
    </div>

    <!-- Loading Indicator -->
    <div class="mt-4 text-sm text-gray-500 dark:text-gray-400">
        <span hx-indicator="true" style="display:none">Loading component...</span>
    </div>
</div>
`
	data := map[string]interface{}{
		"Title":   "Component Showcase",
		"Content": showcaseContent,
		"Theme":   getThemeFromRequest(r),
	}
	RenderTemplate(w, "layouts/sidebar.html", data)
}

// ShowcasePageHandler renders just the showcase content (for HTMX loading)
func (h *UIHandler) ShowcasePageHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":   "Component Showcase",
		"Content": "Component Showcase",
		"Theme":   getThemeFromRequest(r),
	}
	RenderTemplate(w, "showcase.html", data)
}

// ComponentButtonsHandler returns the buttons component
func (h *UIHandler) ComponentButtonsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buttonHTML := `
<h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">Buttons</h2>

<div class="space-y-8">
    <!-- Primary Buttons -->
    <div>
        <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-4">Primary Buttons</h3>
        <div class="flex flex-wrap gap-3 p-4 bg-gray-50 dark:bg-gray-900 rounded-lg">
            <button class="btn-primary">
                Primary Button
            </button>
            <button class="btn-primary">
                <span>✈️</span>
                With Icon
            </button>
            <button class="btn-primary" disabled>
                Disabled
            </button>
        </div>
    </div>

    <!-- Secondary Buttons -->
    <div>
        <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-4">Secondary Buttons</h3>
        <div class="flex flex-wrap gap-3 p-4 bg-gray-50 dark:bg-gray-900 rounded-lg">
            <button class="btn-secondary">
                Secondary
            </button>
            <button class="btn-secondary">
                <span>🔗</span>
                With Icon
            </button>
            <button class="btn-secondary" disabled>
                Disabled
            </button>
        </div>
    </div>

    <!-- Neutral Buttons -->
    <div>
        <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-4">Neutral Buttons</h3>
        <div class="flex flex-wrap gap-3 p-4 bg-gray-50 dark:bg-gray-900 rounded-lg">
            <button class="btn-neutral">
                Neutral
            </button>
            <button class="btn-neutral">
                <span>⚙️</span>
                Settings
            </button>
            <button class="btn-neutral" disabled>
                Disabled
            </button>
        </div>
    </div>

    <!-- Danger Buttons -->
    <div>
        <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-4">Danger Buttons</h3>
        <div class="flex flex-wrap gap-3 p-4 bg-gray-50 dark:bg-gray-900 rounded-lg">
            <button class="btn-danger">
                Delete
            </button>
            <button class="btn-danger">
                <span>🗑️</span>
                Remove
            </button>
            <button class="btn-danger" disabled>
                Disabled
            </button>
        </div>
    </div>

    <!-- Button Sizes -->
    <div>
        <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-4">Button Sizes</h3>
        <div class="flex flex-wrap gap-3 items-center p-4 bg-gray-50 dark:bg-gray-900 rounded-lg">
            <button class="btn-primary btn-small">Small</button>
            <button class="btn-primary">Medium</button>
            <button class="btn-primary btn-large">Large</button>
        </div>
    </div>

    <!-- Button States -->
    <div>
        <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-4">Button States</h3>
        <div class="flex flex-wrap gap-3 p-4 bg-gray-50 dark:bg-gray-900 rounded-lg">
            <button class="btn-primary">Default</button>
            <button class="btn-primary active">Active</button>
            <button class="btn-primary" disabled>Disabled</button>
        </div>
    </div>
</div>
`
	w.Write([]byte(buttonHTML))
}

// ComponentFormsHandler returns the forms and inputs component
func (h *UIHandler) ComponentFormsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	formsHTML := `
<h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">Forms & Inputs</h2>

<div class="max-w-2xl">
    <div class="bg-gray-50 dark:bg-gray-900 p-6 rounded-lg space-y-6">
        <!-- Text Input Group -->
        <div class="form-group">
            <label for="text-input">Full Name</label>
            <input type="text" id="text-input" placeholder="John Doe" class="w-full" />
        </div>

        <!-- Email Input Group -->
        <div class="form-group">
            <label for="email-input">Email Address</label>
            <input type="email" id="email-input" placeholder="john@example.com" class="w-full" />
        </div>

        <!-- Select Dropdown Group -->
        <div class="form-group">
            <label for="select-input">Select Option</label>
            <select id="select-input" class="w-full">
                <option>Choose an option</option>
                <option>Option 1</option>
                <option>Option 2</option>
                <option>Option 3</option>
            </select>
        </div>

        <!-- Textarea Group -->
        <div class="form-group">
            <label for="textarea-input">Message</label>
            <textarea id="textarea-input" placeholder="Enter your message..." rows="4" class="w-full"></textarea>
        </div>

        <!-- Checkbox Group -->
        <div class="form-group">
            <label class="flex items-center gap-3 cursor-pointer">
                <input type="checkbox" class="w-5 h-5" />
                <span>I agree to the terms and conditions</span>
            </label>
        </div>

        <!-- Radio Buttons Group -->
        <div class="form-group">
            <label class="block mb-3">Preferences</label>
            <div class="space-y-2">
                <label class="flex items-center gap-3 cursor-pointer">
                    <input type="radio" name="preference" value="opt1" class="w-5 h-5" />
                    <span>Option A</span>
                </label>
                <label class="flex items-center gap-3 cursor-pointer">
                    <input type="radio" name="preference" value="opt2" class="w-5 h-5" />
                    <span>Option B</span>
                </label>
                <label class="flex items-center gap-3 cursor-pointer">
                    <input type="radio" name="preference" value="opt3" class="w-5 h-5" />
                    <span>Option C</span>
                </label>
            </div>
        </div>

        <!-- Form Actions -->
        <div class="flex gap-3 pt-4 border-t border-gray-200 dark:border-gray-700">
            <button class="btn-primary">Submit</button>
            <button class="btn-secondary">Cancel</button>
        </div>
    </div>
</div>
`
	w.Write([]byte(formsHTML))
}

// ComponentTypographyHandler returns the typography component
func (h *UIHandler) ComponentTypographyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	typographyHTML := `
<h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">Typography</h2>

<div class="space-y-6">
    <div>
        <h1 class="text-4xl font-bold text-gray-900 dark:text-white mb-2">Heading 1</h1>
        <p class="text-sm text-gray-500 dark:text-gray-400">Main page title</p>
    </div>

    <div>
        <h2 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">Heading 2</h2>
        <p class="text-sm text-gray-500 dark:text-gray-400">Section title</p>
    </div>

    <div>
        <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">Heading 3</h3>
        <p class="text-sm text-gray-500 dark:text-gray-400">Subsection title</p>
    </div>

    <div>
        <h4 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">Heading 4</h4>
        <p class="text-sm text-gray-500 dark:text-gray-400">Component title</p>
    </div>

    <div class="pt-4 border-t border-gray-200 dark:border-gray-700">
        <p class="text-base text-gray-700 dark:text-gray-300 mb-2">Regular paragraph text</p>
        <p class="text-sm text-gray-600 dark:text-gray-400 mb-2">Small text</p>
        <p class="text-xs text-gray-500 dark:text-gray-500 mb-2">Extra small text</p>
    </div>

    <div class="pt-4 border-t border-gray-200 dark:border-gray-700">
        <p class="font-bold text-gray-900 dark:text-white">Bold text</p>
        <p class="italic text-gray-700 dark:text-gray-300">Italic text</p>
        <p class="underline text-gray-700 dark:text-gray-300">Underlined text</p>
    </div>

    <div class="pt-4 border-t border-gray-200 dark:border-gray-700">
        <p class="bg-yellow-100 dark:bg-yellow-900 text-yellow-800 dark:text-yellow-200 p-2 rounded">Highlighted text</p>
    </div>
</div>
`
	w.Write([]byte(typographyHTML))
}

// ComponentCardsHandler returns the cards and tables component
func (h *UIHandler) ComponentCardsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	cardsHTML := `
<h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">Cards & Tables</h2>

<div class="space-y-8">
    <!-- Standard Cards -->
    <div>
        <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-4">Standard Cards</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Basic Card -->
            <div class="card">
                <div class="card-header">
                    <h4 class="card-title">Flight Details</h4>
                </div>
                <div class="card-body">
                    <p class="text-gray-600 dark:text-gray-400">Route: LFPG-EGLL</p>
                    <p class="text-gray-600 dark:text-gray-400">Aircraft: Boeing 777-300ER</p>
                    <p class="text-gray-600 dark:text-gray-400">Distance: 344 nm</p>
                </div>
                <div class="card-footer">
                    <button class="btn-primary">View Details</button>
                </div>
            </div>

            <!-- Featured Card -->
            <div class="card card-featured">
                <div class="card-header">
                    <h4 class="card-title">Featured Flight</h4>
                    <span class="badge badge-primary">Premium</span>
                </div>
                <div class="card-body">
                    <p class="text-gray-600 dark:text-gray-400">Route: KJFK-LFPG</p>
                    <p class="text-gray-600 dark:text-gray-400">Aircraft: Airbus A380</p>
                    <p class="text-gray-600 dark:text-gray-400">Distance: 3,626 nm</p>
                </div>
                <div class="card-footer">
                    <button class="btn-primary">Book Flight</button>
                </div>
            </div>
        </div>
    </div>

    <!-- Data Table -->
    <div>
        <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-4">Flight Log</h3>
        <div class="overflow-x-auto">
            <table class="w-full text-left border-collapse">
                <thead>
                    <tr class="bg-gray-100 dark:bg-gray-900 border-b border-gray-200 dark:border-gray-700">
                        <th class="px-4 py-3 font-semibold text-gray-900 dark:text-white">Flight Route</th>
                        <th class="px-4 py-3 font-semibold text-gray-900 dark:text-white">Aircraft</th>
                        <th class="px-4 py-3 font-semibold text-gray-900 dark:text-white">Distance</th>
                        <th class="px-4 py-3 font-semibold text-gray-900 dark:text-white">Status</th>
                    </tr>
                </thead>
                <tbody>
                    <tr class="border-b border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-900 transition-colors">
                        <td class="px-4 py-3 text-gray-700 dark:text-gray-300">LFPG-EGLL</td>
                        <td class="px-4 py-3 text-gray-700 dark:text-gray-300">B777-300ER</td>
                        <td class="px-4 py-3 text-gray-700 dark:text-gray-300">344 nm</td>
                        <td class="px-4 py-3"><span class="badge badge-success">Completed</span></td>
                    </tr>
                    <tr class="border-b border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-900 transition-colors">
                        <td class="px-4 py-3 text-gray-700 dark:text-gray-300">KJFK-UUWW</td>
                        <td class="px-4 py-3 text-gray-700 dark:text-gray-300">A380-800</td>
                        <td class="px-4 py-3 text-gray-700 dark:text-gray-300">5,395 nm</td>
                        <td class="px-4 py-3"><span class="badge badge-primary">In Progress</span></td>
                    </tr>
                    <tr class="border-b border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-900 transition-colors">
                        <td class="px-4 py-3 text-gray-700 dark:text-gray-300">EGLL-KMIA</td>
                        <td class="px-4 py-3 text-gray-700 dark:text-gray-300">B787-9</td>
                        <td class="px-4 py-3 text-gray-700 dark:text-gray-300">3,459 nm</td>
                        <td class="px-4 py-3"><span class="badge badge-warning">Pending</span></td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</div>
`
	w.Write([]byte(cardsHTML))
}

// ComponentBadgesHandler returns the badges and chips component
func (h *UIHandler) ComponentBadgesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	badgesHTML := `
<h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">Badges & Chips</h2>

<div class="space-y-8">
    <!-- Badge Variants -->
    <div>
        <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-4">Badge Variants</h3>
        <div class="p-4 bg-gray-50 dark:bg-gray-900 rounded-lg flex flex-wrap gap-3">
            <span class="badge badge-primary">
                <span>✈️</span>
                Primary
            </span>
            <span class="badge badge-success">
                <span>✓</span>
                Success
            </span>
            <span class="badge badge-error">
                <span>✕</span>
                Error
            </span>
            <span class="badge badge-warning">
                <span>⚠️</span>
                Warning
            </span>
            <span class="badge badge-default">
                <span>•</span>
                Default
            </span>
        </div>
    </div>

    <!-- Dismissible Chips -->
    <div>
        <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-4">Dismissible Chips</h3>
        <div class="p-4 bg-gray-50 dark:bg-gray-900 rounded-lg flex flex-wrap gap-3">
            <div class="chip badge-primary">
                <span>🎯 Priority Filter</span>
                <button class="chip-close" onclick="this.parentElement.remove()">×</button>
            </div>
            <div class="chip badge-success">
                <span>✓ Completed</span>
                <button class="chip-close" onclick="this.parentElement.remove()">×</button>
            </div>
            <div class="chip badge-warning">
                <span>⏳ Pending Review</span>
                <button class="chip-close" onclick="this.parentElement.remove()">×</button>
            </div>
        </div>
    </div>

    <!-- Tag Cloud -->
    <div>
        <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-4">Topic Tags</h3>
        <div class="p-4 bg-gray-50 dark:bg-gray-900 rounded-lg flex flex-wrap gap-2">
            <span class="px-3 py-1 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded text-sm hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer transition-colors">
                ✈️ Aviation
            </span>
            <span class="px-3 py-1 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded text-sm hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer transition-colors">
                🛫 Flight Operations
            </span>
            <span class="px-3 py-1 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded text-sm hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer transition-colors">
                🗺️ Routes
            </span>
            <span class="px-3 py-1 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded text-sm hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer transition-colors">
                🛩️ Aircraft
            </span>
            <span class="px-3 py-1 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded text-sm hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer transition-colors">
                👨‍✈️ Pilots
            </span>
        </div>
    </div>
</div>
`
	w.Write([]byte(badgesHTML))
}

// ComponentAlertsHandler returns the alerts component
func (h *UIHandler) ComponentAlertsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	alertsHTML := `
<h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">Alerts</h2>

<div class="space-y-4 max-w-2xl">
    <!-- Info Alert -->
    <div class="alert alert-info">
        <div class="alert-icon">ℹ️</div>
        <div class="alert-content">
            <h4 class="font-semibold mb-1">Information</h4>
            <p class="text-sm">This is an informational message. You may want to read this.</p>
        </div>
        <button class="alert-close">×</button>
    </div>

    <!-- Success Alert -->
    <div class="alert alert-success">
        <div class="alert-icon">✓</div>
        <div class="alert-content">
            <h4 class="font-semibold mb-1">Success</h4>
            <p class="text-sm">Your flight has been successfully filed and recorded.</p>
        </div>
        <button class="alert-close">×</button>
    </div>

    <!-- Warning Alert -->
    <div class="alert alert-warning">
        <div class="alert-icon">⚠️</div>
        <div class="alert-content">
            <h4 class="font-semibold mb-1">Warning</h4>
            <p class="text-sm">Please review this important information before proceeding.</p>
        </div>
        <button class="alert-close">×</button>
    </div>

    <!-- Error Alert -->
    <div class="alert alert-error">
        <div class="alert-icon">✕</div>
        <div class="alert-content">
            <h4 class="font-semibold mb-1">Error</h4>
            <p class="text-sm">An error occurred while processing your request. Please try again.</p>
        </div>
        <button class="alert-close">×</button>
    </div>

    <!-- Modal Demo Section -->
    <div class="mt-8 pt-8 border-t border-gray-200 dark:border-gray-700">
        <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-4">Modal Example</h3>
        <button class="btn-primary" onclick="document.getElementById('demoModal').classList.add('modal-open')">
            Open Modal
        </button>
    </div>
</div>

<!-- Modal Backdrop -->
<div id="demoModal" class="modal">
    <div class="modal-backdrop" onclick="document.getElementById('demoModal').classList.remove('modal-open')"></div>

    <!-- Modal Content -->
    <div class="modal-content modal-md">
        <!-- Modal Header -->
        <div class="modal-header">
            <h2 class="modal-title">Flight Details</h2>
            <button class="modal-close" onclick="document.getElementById('demoModal').classList.remove('modal-open')">×</button>
        </div>

        <!-- Modal Body -->
        <div class="modal-body">
            <div class="space-y-4">
                <p class="text-gray-600 dark:text-gray-400">
                    This is a modal dialog showcasing the flight details form.
                </p>

                <div class="form-group">
                    <label for="modal-route">Route</label>
                    <input type="text" id="modal-route" value="LFPG-EGLL" placeholder="Enter route" class="w-full" />
                </div>

                <div class="form-group">
                    <label for="modal-aircraft">Aircraft</label>
                    <input type="text" id="modal-aircraft" value="B777-300ER" placeholder="Enter aircraft" class="w-full" />
                </div>

                <div class="form-group">
                    <label for="modal-time">Flight Time</label>
                    <input type="text" id="modal-time" value="02:15" placeholder="Enter flight time" class="w-full" />
                </div>
            </div>
        </div>

        <!-- Modal Footer -->
        <div class="modal-footer">
            <button class="btn-secondary" onclick="document.getElementById('demoModal').classList.remove('modal-open')">
                Cancel
            </button>
            <button class="btn-primary">
                Save Flight
            </button>
        </div>
    </div>
</div>
`
	w.Write([]byte(alertsHTML))
}
