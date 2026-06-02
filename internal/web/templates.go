package web

// IndexTemplate is the embedded HTML/CSS/JS template for our local web wizard.
const IndexTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Config Maker Dashboard</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600&family=Outfit:wght@400;600;800&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg-color: #09090b;      /* Pure matte black (zinc-950) */
            --card-bg: #0f0f11;       /* Solid slate-gray card (zinc-900) */
            --border-color: #27272a;  /* Crisp dark-gray border (zinc-800) */
            --primary-accent: #3b82f6; /* High-contrast royal blue */
            --secondary-accent: #1d4ed8; /* Darker blue for state transitions */
            --text-main: #fafafa;     /* Clean bright off-white (zinc-50) */
            --text-muted: #8e8e93;    /* Professional medium gray */
            --success-color: #10b981; /* Pure emerald green */
            --error-color: #ef4444;   /* Solid red */
            --warning-color: #f59e0b; /* Solid amber */
        }

        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }

        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
            background-color: var(--bg-color);
            color: var(--text-main);
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            overflow-x: hidden;
            position: relative;
        }

        .container {
            width: 100%;
            max-width: 680px;
            padding: 40px 24px;
            z-index: 10;
        }

        /* Minimal Header */
        header {
            text-align: center;
            margin-bottom: 40px;
        }
        header h1 {
            font-size: 2.2rem;
            font-weight: 700;
            color: var(--text-main);
            letter-spacing: -0.8px;
            margin-bottom: 8px;
        }
        header p {
            color: var(--text-muted);
            font-size: 0.95rem;
            font-weight: 400;
        }

        /* Premium Minimalist Card */
        .glass-card {
            background: var(--card-bg);
            border: 1px solid var(--border-color);
            border-radius: 12px;
            padding: 36px;
            box-shadow: 0 12px 40px rgba(0, 0, 0, 0.4);
            transition: all 0.2s ease;
            position: relative;
            overflow: hidden;
        }

        /* Quiet Progress Bar */
        .progress-container {
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin-bottom: 28px;
            background: none;
            padding: 0;
            border: none;
        }
        .progress-bar-wrapper {
            flex-grow: 1;
            height: 4px;
            background: #1c1c1e;
            border-radius: 2px;
            margin: 0 16px;
            overflow: hidden;
            position: relative;
        }
        .progress-bar-fill {
            height: 100%;
            width: 0%;
            background-color: var(--primary-accent);
            border-radius: 2px;
            transition: width 0.3s ease;
        }
        .progress-step-text {
            font-size: 0.8rem;
            color: var(--text-muted);
            font-weight: 500;
            min-width: 65px;
        }

        /* Step Content transitions */
        .step-content {
            display: none;
            animation: fadeIn 0.2s ease forwards;
        }
        .step-content.active {
            display: block;
        }

        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(4px); }
            to { opacity: 1; transform: translateY(0); }
        }

        h2 {
            font-size: 1.3rem;
            font-weight: 600;
            margin-bottom: 20px;
            color: var(--text-main);
            letter-spacing: -0.3px;
        }

        /* Forms and Controls */
        .form-group {
            margin-bottom: 20px;
        }
        .form-group label {
            display: block;
            font-size: 0.85rem;
            color: var(--text-muted);
            margin-bottom: 8px;
            font-weight: 500;
        }
        .text-input {
            width: 100%;
            background: #161618;
            border: 1px solid var(--border-color);
            border-radius: 6px;
            padding: 12px 14px;
            color: var(--text-main);
            font-family: inherit;
            font-size: 0.95rem;
            transition: all 0.15s ease;
        }
        .text-input:focus {
            outline: none;
            border-color: var(--primary-accent);
            background: #18181b;
        }

        /* Minimal Switch Card */
        .toggle-card {
            background: #141416;
            border: 1px solid var(--border-color);
            border-radius: 8px;
            padding: 16px 20px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin-bottom: 16px;
            transition: all 0.2s ease;
        }
        .toggle-card:hover {
            border-color: #3f3f46;
            background: #161619;
        }
        .toggle-info {
            max-width: 80%;
        }
        .toggle-title {
            font-weight: 600;
            font-size: 0.95rem;
            margin-bottom: 2px;
        }
        .toggle-desc {
            font-size: 0.8rem;
            color: var(--text-muted);
            line-height: 1.4;
        }
        
        .switch {
            position: relative;
            display: inline-block;
            width: 44px;
            height: 24px;
        }
        .switch input {
            opacity: 0;
            width: 0;
            height: 0;
        }
        .slider {
            position: absolute;
            cursor: pointer;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background-color: #27272a;
            transition: .2s ease;
            border-radius: 24px;
        }
        .slider:before {
            position: absolute;
            content: "";
            height: 18px;
            width: 18px;
            left: 3px;
            bottom: 3px;
            background-color: #fff;
            transition: .2s ease;
            border-radius: 50%;
        }
        input:checked + .slider {
            background-color: var(--primary-accent);
        }
        input:checked + .slider:before {
            transform: translateX(20px);
        }

        /* Choices Grid (e.g. Dark/Light Mode selectors) */
        .options-grid {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 16px;
            margin-bottom: 24px;
        }
        .choice-box {
            background: #141416;
            border: 1px solid var(--border-color);
            border-radius: 8px;
            padding: 16px;
            text-align: center;
            cursor: pointer;
            transition: all 0.2s ease;
        }
        .choice-box:hover {
            border-color: #3f3f46;
            background: #161619;
        }
        .choice-box.selected {
            background: #172554;
            border-color: var(--primary-accent);
        }
        .choice-icon {
            font-size: 1.5rem;
            margin-bottom: 6px;
        }
        .choice-title {
            font-weight: 600;
            font-size: 0.9rem;
        }

        .select-input {
            width: 100%;
            background: #161618;
            border: 1px solid var(--border-color);
            border-radius: 6px;
            padding: 12px 14px;
            color: var(--text-main);
            font-family: inherit;
            font-size: 0.95rem;
            cursor: pointer;
            outline: none;
        }
        .select-input:focus {
            border-color: var(--primary-accent);
        }

        /* Minimal Summary Items list */
        .summary-list {
            background: #141416;
            border: 1px solid var(--border-color);
            border-radius: 8px;
            padding: 20px;
            margin-bottom: 20px;
        }
        .summary-row {
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 10px 0;
            border-bottom: 1px solid #1c1c1f;
        }
        .summary-row:last-child {
            border-bottom: none;
        }
        .summary-label {
            font-weight: 500;
            font-size: 0.9rem;
            color: var(--text-muted);
        }
        .summary-value {
            display: flex;
            align-items: center;
            font-size: 0.9rem;
            font-weight: 600;
        }
        .indicator {
            width: 10px;
            height: 10px;
            border-radius: 50%;
            display: inline-block;
            margin-right: 8px;
        }
        .indicator.enabled {
            background-color: var(--success-color);
        }
        .indicator.disabled {
            background-color: #3f3f46;
        }

        /* Sleek flat buttons */
        .btn-row {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-top: 28px;
        }
        .btn {
            font-family: inherit;
            font-size: 0.9rem;
            font-weight: 600;
            padding: 12px 24px;
            border-radius: 6px;
            cursor: pointer;
            transition: all 0.15s ease;
            border: none;
            display: inline-flex;
            align-items: center;
            justify-content: center;
            text-decoration: none;
        }
        .btn-prev {
            background: #18181b;
            color: var(--text-main);
            border: 1px solid var(--border-color);
        }
        .btn-prev:hover {
            background: #27272a;
            border-color: #3f3f46;
        }
        .btn-next {
            background: var(--primary-accent);
            color: white;
        }
        .btn-next:hover {
            background: var(--secondary-accent);
        }

        /* Solid Console Logger (Matte black/monochrome output) */
        .console-card {
            display: none;
            background: #000;
            border: 1px solid var(--border-color);
            border-radius: 8px;
            padding: 20px;
            font-family: 'SFMono-Regular', Consolas, "Liberation Mono", Menlo, Courier, monospace;
            height: 360px;
            overflow-y: auto;
            color: #d4d4d8;
            margin-bottom: 20px;
            box-shadow: inset 0 2px 10px rgba(0,0,0,0.9);
        }
        .console-card.active {
            display: block;
        }
        .console-line {
            margin-bottom: 6px;
            white-space: pre-wrap;
            line-height: 1.4;
            font-size: 0.85rem;
        }
        .console-pulse-container {
            display: flex;
            align-items: center;
            margin-bottom: 14px;
            color: var(--text-muted);
            font-size: 0.8rem;
        }
        .pulse-dot {
            width: 6px;
            height: 6px;
            border-radius: 50%;
            background-color: var(--primary-accent);
            margin-right: 8px;
            animation: pulse 1s infinite;
        }
        @keyframes pulse {
            0% { opacity: 0.4; }
            50% { opacity: 1; }
            100% { opacity: 0.4; }
        }

        /* Clean Finish View */
        .finish-view {
            display: none;
            text-align: center;
            padding: 20px 0;
        }
        .finish-view.active {
            display: block;
        }
        .finish-icon {
            font-size: 3rem;
            margin-bottom: 12px;
        }
    </style>
</head>
<body>

    <div class="container">
        <header>
            <h1>config-maker</h1>
            <p>Sleek Desktop Configuration Wizard</p>
            <div style="margin-top: 20px; display: flex; justify-content: center; gap: 16px;">
                <button class="btn btn-prev" onclick="manualImport()" style="padding: 10px 20px; font-size: 0.9rem; border-radius: 12px; margin-top: 0; background: rgba(14, 165, 233, 0.1); border-color: var(--primary-accent);">📥 Import Settings</button>
                <button class="btn btn-prev" onclick="manualExport()" style="padding: 10px 20px; font-size: 0.9rem; border-radius: 12px; margin-top: 0; background: rgba(139, 92, 246, 0.1); border-color: var(--secondary-accent);">💾 Export Settings</button>
            </div>
            <div id="toastNotification" style="display: none; margin-top: 16px; padding: 10px 20px; border-radius: 12px; font-size: 0.9rem; font-weight: 500; text-align: center; animation: fadeInSlide 0.3s ease;"></div>
        </header>

        <div class="glass-card" id="wizardCard">
            <!-- Progress Tracker -->
            <div class="progress-container" id="progressBarContainer">
                <span class="progress-step-text" id="progressStepNum">Step 1 of 6</span>
                <div class="progress-bar-wrapper">
                    <div class="progress-bar-fill" id="progressBarFill" style="width: 16.66%;"></div>
                </div>
                <span class="progress-step-text" id="progressPercent">16%</span>
            </div>

            <!-- Import Banner for skipping steps when settings are imported -->
            <div id="importBanner" style="display: none; background: rgba(14, 165, 233, 0.08); border: 1px solid var(--primary-accent); border-radius: 16px; padding: 16px; margin-bottom: 24px; justify-content: space-between; align-items: center;">
                <div style="text-align: left;">
                    <div style="font-weight: 600; font-size: 0.95rem; color: var(--primary-accent);">Imported Saved Settings</div>
                    <div style="font-size: 0.8rem; color: var(--text-muted);">You can skip straight to desktop wallpaper selection or review your choices.</div>
                </div>
                <button class="btn btn-next" onclick="showStep(4)" style="padding: 8px 16px; font-size: 0.85rem; border-radius: 10px; margin-top: 0; background: linear-gradient(90deg, var(--primary-accent), var(--secondary-accent)); box-shadow: none;">Go to Wallpaper ➔</button>
            </div>

            <!-- STEP 1: Oh My Zsh -->
            <div class="step-content active" id="step1">
                <h2>Zsh Shell & Extensions</h2>
                <div class="toggle-card">
                    <div class="toggle-info">
                        <div class="toggle-title">Install Oh-My-Zsh</div>
                        <div class="toggle-desc">Installs the popular community extension for Zsh to manage themes and plugins unattended.</div>
                    </div>
                    <label class="switch">
                        <input type="checkbox" id="installZsh" checked>
                        <span class="slider"></span>
                    </label>
                </div>
            </div>

            <!-- STEP 2: Git Config -->
            <div class="step-content" id="step2">
                <h2>Git Credentials Setup</h2>
                <div class="toggle-card" style="margin-bottom: 24px;">
                    <div class="toggle-info">
                        <div class="toggle-title">Enable Git Setup</div>
                        <div class="toggle-desc">Configure global credential storage and user credentials.</div>
                    </div>
                    <label class="switch">
                        <input type="checkbox" id="setupGit" checked onchange="toggleGitFields()">
                        <span class="slider"></span>
                    </label>
                </div>
                <div id="gitFields">
                    <div class="form-group">
                        <label for="gitName">Full Name / Login</label>
                        <input type="text" id="gitName" class="text-input" placeholder="e.g. John Doe">
                    </div>
                    <div class="form-group">
                        <label for="gitEmail">Email Address</label>
                        <input type="email" id="gitEmail" class="text-input" placeholder="e.g. john@example.com">
                    </div>
                </div>
            </div>

            <!-- STEP 3: Themes -->
            <div class="step-content" id="step3">
                <h2>Gnome Interface Theme</h2>
                <div class="toggle-card" style="margin-bottom: 24px;">
                    <div class="toggle-info">
                        <div class="toggle-title">Configure Theme</div>
                        <div class="toggle-desc">Apply customized GTK and Window themes from your local system resources.</div>
                    </div>
                    <label class="switch">
                        <input type="checkbox" id="setupTheme" checked onchange="toggleThemeFields()">
                        <span class="slider"></span>
                    </label>
                </div>
                <div id="themeFields">
                    <label style="display: block; font-size: 0.9rem; color: var(--text-muted); margin-bottom: 8px; font-weight: 500;">Select Mode</label>
                    <div class="options-grid">
                        <div class="choice-box selected" id="themeModeDark" onclick="selectThemeMode('1')">
                            <div class="choice-icon">🌙</div>
                            <div class="choice-title">Dark Mode</div>
                        </div>
                        <div class="choice-box" id="themeModeLight" onclick="selectThemeMode('2')">
                            <div class="choice-icon">☀️</div>
                            <div class="choice-title">Light Mode</div>
                        </div>
                    </div>
                    <div class="form-group">
                        <label for="themeSelect">Select Target Theme</label>
                        <select id="themeSelect" class="select-input">
                            <!-- Populated dynamically via JS -->
                        </select>
                    </div>
                </div>
            </div>

            <!-- STEP 4: Background -->
            <div class="step-content" id="step4">
                <h2>Desktop Wallpaper</h2>
                <div class="toggle-card" style="margin-bottom: 24px;">
                    <div class="toggle-info">
                        <div class="toggle-title">Apply Custom Background</div>
                        <div class="toggle-desc">Update user desktop backgrounds using selected wallpapers.</div>
                    </div>
                    <label class="switch">
                        <input type="checkbox" id="setupBg" checked onchange="toggleBgFields()">
                        <span class="slider"></span>
                    </label>
                </div>
                <div id="bgFields">
                    <div class="form-group">
                        <label for="bgSource">Select Background Option</label>
                        <select id="bgSource" class="select-input" onchange="toggleBgInputs()">
                            <option value="1">Predefined Background.jpeg</option>
                            <option value="2">Select from repository wallpapers</option>
                            <option value="3">Custom absolute image path</option>
                        </select>
                    </div>
                    <div class="form-group" id="repoWallpapersWrapper" style="display: none;">
                        <label for="repoWpSelect">Available Wallpaper</label>
                        <select id="repoWpSelect" class="select-input">
                            <!-- Populated from system values -->
                        </select>
                    </div>
                    <div class="form-group" id="customPathWrapper" style="display: none;">
                        <label for="customBgPath">Absolute Path</label>
                        <input type="text" id="customBgPath" class="text-input" placeholder="e.g. /home/user/Pictures/wallpaper.jpg">
                    </div>
                </div>
            </div>

            <!-- STEP 5: Docker and Default Shell -->
            <div class="step-content" id="step5">
                <h2>Environments & Shells</h2>
                <div class="toggle-card">
                    <div class="toggle-info">
                        <div class="toggle-title">Enable Docker Rootless</div>
                        <div class="toggle-desc">Installs Docker inside user environment. Fully safe, requires no sudo.</div>
                    </div>
                    <label class="switch">
                        <input type="checkbox" id="installDocker">
                        <span class="slider"></span>
                    </label>
                </div>
                <div class="toggle-card">
                    <div class="toggle-info">
                        <div class="toggle-title">Set Zsh as Default Shell</div>
                        <div class="toggle-desc">Sets up .bashrc script commands to redirect bash to login Zsh shell automatically.</div>
                    </div>
                    <label class="switch">
                        <input type="checkbox" id="defaultShellZsh" checked>
                        <span class="slider"></span>
                    </label>
                </div>
            </div>

            <!-- STEP 6: Summary & Confirmation -->
            <div class="step-content" id="step6">
                <h2>Setup Overview</h2>
                <div class="summary-list">
                    <div class="summary-row">
                        <span class="summary-label">Install Oh-My-Zsh</span>
                        <span class="summary-value" id="sumZsh"><span class="indicator"></span><span></span></span>
                    </div>
                    <div class="summary-row">
                        <span class="summary-label">Configure Git Credentials</span>
                        <span class="summary-value" id="sumGit"><span class="indicator"></span><span></span></span>
                    </div>
                    <div class="summary-row">
                        <span class="summary-label">Apply System Theme</span>
                        <span class="summary-value" id="sumTheme"><span class="indicator"></span><span></span></span>
                    </div>
                    <div class="summary-row">
                        <span class="summary-label">Set Desktop Background</span>
                        <span class="summary-value" id="sumBg"><span class="indicator"></span><span></span></span>
                    </div>
                    <div class="summary-row">
                        <span class="summary-label">Enable Docker Rootless</span>
                        <span class="summary-value" id="sumDocker"><span class="indicator"></span><span></span></span>
                    </div>
                    <div class="summary-row">
                        <span class="summary-label">Zsh as Default Shell</span>
                        <span class="summary-value" id="sumShell"><span class="indicator"></span><span></span></span>
                    </div>
                </div>
                <div class="toggle-card" style="margin-top: 24px;">
                    <div class="toggle-info" style="text-align: left;">
                        <div class="toggle-title">Save & Export Settings</div>
                        <div class="toggle-desc">Automatically save these choices to ~/.config/config-maker/config.json.</div>
                    </div>
                    <label class="switch">
                        <input type="checkbox" id="exportSettingsWeb" checked>
                        <span class="slider"></span>
                    </label>
                </div>
            </div>

            <!-- Streaming Console View -->
            <div class="console-pulse-container" id="pulseContainer" style="display: none;">
                <div class="pulse-dot"></div>
                <span id="consoleStatusText">Applying changes... Please wait.</span>
            </div>
            <div class="console-card" id="consoleCard">
                <!-- Exec logs shown live -->
            </div>

            <!-- Finish View -->
            <div class="finish-view" id="finishView">
                <div class="finish-icon">🎉</div>
                <h2>Setup Complete!</h2>
                <p style="color: var(--text-muted); margin-bottom: 24px; line-height: 1.6;">
                    All configurations have been successfully completed.<br>
                    You can close this tab and restart your terminal program to experience the upgrades.
                </p>
                <button class="btn btn-next" onclick="restartTerminal()" style="width: 100%;">Reopen & Reload Terminals</button>
            </div>

            <!-- Button navigation row -->
            <div class="btn-row" id="btnRow">
                <button class="btn btn-prev" id="btnPrev" onclick="prevStep()">Previous</button>
                <button class="btn btn-next" id="btnNext" onclick="nextStep()">Next</button>
            </div>
        </div>
    </div>

    <script>
        let currentStep = 1;
        const totalSteps = 6;
        let importedSettings = false;

        let selectedThemeMode = "1"; // "1" Dark, "2" Light
        let availableThemes = {
            dark: [],
            light: []
        };
        let wallpapers = [];

        // Load themes and wallpapers dynamically from the backend
        async function fetchSystemResources() {
            try {
                const response = await fetch('/api/resources');
                const data = await response.json();
                availableThemes.dark = data.dark_themes || [];
                availableThemes.light = data.light_themes || [];
                wallpapers = data.wallpapers || [];

                populateThemeDropdown();
                populateWallpaperDropdown();

                // Fetch loaded config
                const configRes = await fetch('/api/config');
                const configData = await configRes.json();
                applyConfigToUI(configData);
                if (configData && (configData.git_name || configData.git_email || configData.theme_name || configData.background_image)) {
                    showToast("Previously saved settings detected and imported!", true);
                    showStep(1); // Force banner refresh
                }
            } catch (err) {
                console.error("Failed to load resources:", err);
            }
        }

        function applyConfigToUI(cfg) {
            if (!cfg) return;

            if (cfg.git_name || cfg.git_email || cfg.theme_name || cfg.background_image) {
                importedSettings = true;
            }

            // Oh My Zsh
            document.getElementById("installZsh").checked = cfg.install_oh_my_zsh !== false;

            // Git Setup
            document.getElementById("setupGit").checked = cfg.configure_git === true;
            document.getElementById("gitName").value = cfg.git_name || "";
            document.getElementById("gitEmail").value = cfg.git_email || "";
            toggleGitFields();

            // Gnome Themes
            document.getElementById("setupTheme").checked = cfg.apply_theme !== false;
            selectedThemeMode = cfg.theme_mode || "1";
            document.getElementById("themeModeDark").classList.toggle("selected", selectedThemeMode === "1");
            document.getElementById("themeModeLight").classList.toggle("selected", selectedThemeMode === "2");
            populateThemeDropdown();
            if (cfg.theme_name) {
                document.getElementById("themeSelect").value = cfg.theme_name;
            }
            toggleThemeFields();

            // Background Desktop Wallpaper
            document.getElementById("setupBg").checked = cfg.apply_background !== false;
            const bgImage = cfg.background_image || "";
            if (bgImage === "" || bgImage.endsWith("Background.jpeg")) {
                document.getElementById("bgSource").value = "1";
            } else if (bgImage.includes("wallpapers/")) {
                document.getElementById("bgSource").value = "2";
                const parts = bgImage.split("/");
                const filename = parts[parts.length - 1];
                document.getElementById("repoWpSelect").value = filename;
            } else {
                document.getElementById("bgSource").value = "3";
                document.getElementById("customBgPath").value = bgImage;
            }
            toggleBgInputs();
            toggleBgFields();

            // Docker & Default Shell
            document.getElementById("installDocker").checked = cfg.enable_docker === true;
            document.getElementById("defaultShellZsh").checked = cfg.enable_zsh_default !== false;
        }

        function populateThemeDropdown() {
            const select = document.getElementById("themeSelect");
            select.innerHTML = "";
            const list = selectedThemeMode === "1" ? availableThemes.dark : availableThemes.light;

            if (list.length === 0) {
                const opt = document.createElement("option");
                opt.value = selectedThemeMode === "1" ? "Yaru-dark" : "Yaru";
                opt.textContent = selectedThemeMode === "1" ? "Yaru-dark" : "Yaru";
                select.appendChild(opt);
                return;
            }

            list.forEach(t => {
                const opt = document.createElement("option");
                opt.value = t;
                opt.textContent = t;
                select.appendChild(opt);
            });
        }

        function populateWallpaperDropdown() {
            const select = document.getElementById("repoWpSelect");
            select.innerHTML = "";
            wallpapers.forEach(wp => {
                const opt = document.createElement("option");
                opt.value = wp;
                opt.textContent = wp;
                select.appendChild(opt);
            });
        }

        function selectThemeMode(mode) {
            selectedThemeMode = mode;
            document.getElementById("themeModeDark").classList.toggle("selected", mode === "1");
            document.getElementById("themeModeLight").classList.toggle("selected", mode === "2");
            populateThemeDropdown();
        }

        function toggleGitFields() {
            document.getElementById("gitFields").style.opacity = document.getElementById("setupGit").checked ? "1" : "0.4";
            document.getElementById("gitFields").style.pointerEvents = document.getElementById("setupGit").checked ? "auto" : "none";
        }

        function toggleThemeFields() {
            document.getElementById("themeFields").style.opacity = document.getElementById("setupTheme").checked ? "1" : "0.4";
            document.getElementById("themeFields").style.pointerEvents = document.getElementById("setupTheme").checked ? "auto" : "none";
        }

        function toggleBgFields() {
            document.getElementById("bgFields").style.opacity = document.getElementById("setupBg").checked ? "1" : "0.4";
            document.getElementById("bgFields").style.pointerEvents = document.getElementById("setupBg").checked ? "auto" : "none";
        }

        function toggleBgInputs() {
            const src = document.getElementById("bgSource").value;
            document.getElementById("repoWallpapersWrapper").style.display = src === "2" ? "block" : "none";
            document.getElementById("customPathWrapper").style.display = src === "3" ? "block" : "none";
        }

        function prevStep() {
            if (currentStep > 1) {
                showStep(currentStep - 1);
            }
        }

        function nextStep() {
            if (currentStep < totalSteps) {
                showStep(currentStep + 1);
            } else if (currentStep === totalSteps) {
                submitConfig();
            }
        }

        function showStep(step) {
            document.getElementById(` + "`" + `step${currentStep}` + "`" + `).classList.remove("active");
            currentStep = step;
            document.getElementById(` + "`" + `step${currentStep}` + "`" + `).classList.add("active");

            // Progress bar
            const percent = Math.round((currentStep / totalSteps) * 100);
            document.getElementById("progressBarFill").style.width = percent + "%";
            document.getElementById("progressPercent").textContent = percent + "%";
            document.getElementById("progressStepNum").textContent = ` + "`" + `Step ${currentStep} of ${totalSteps}` + "`" + `;

            // Buttons
            document.getElementById("btnPrev").style.visibility = currentStep === 1 ? "hidden" : "visible";
            document.getElementById("btnNext").textContent = currentStep === totalSteps ? "Apply Config" : "Next";

            // Import banner visibility logic
            const banner = document.getElementById("importBanner");
            if (banner) {
                if (importedSettings && currentStep !== 4 && currentStep !== 6) {
                    banner.style.display = "flex";
                } else {
                    banner.style.display = "none";
                }
            }

            if (currentStep === 6) {
                renderSummary();
            }
        }

        function renderSummary() {
            updateSummaryRow("sumZsh", document.getElementById("installZsh").checked, "Install Oh-My-Zsh", "Skip Oh-My-Zsh");
            
            const gitEnabled = document.getElementById("setupGit").checked;
            const gitName = document.getElementById("gitName").value.trim();
            updateSummaryRow("sumGit", gitEnabled, gitEnabled ? ` + "`" + `Configure Git (${gitName || "Default"})` + "`" + ` : "Configure Git", "Skip Git");

            const themeEnabled = document.getElementById("setupTheme").checked;
            const themeName = document.getElementById("themeSelect").value;
            const modeName = selectedThemeMode === "1" ? "Dark" : "Light";
            updateSummaryRow("sumTheme", themeEnabled, themeEnabled ? ` + "`" + `Apply ${modeName} Theme (${themeName})` + "`" + ` : "Apply Theme", "Skip Theme");

            const bgEnabled = document.getElementById("setupBg").checked;
            let bgVal = "Default Background.jpeg";
            const bgSrc = document.getElementById("bgSource").value;
            if (bgSrc === "2") bgVal = "Wallpaper: " + document.getElementById("repoWpSelect").value;
            else if (bgSrc === "3") bgVal = "Custom Path: " + document.getElementById("customBgPath").value;
            updateSummaryRow("sumBg", bgEnabled, bgEnabled ? ` + "`" + `Set Wallpaper (${bgVal})` + "`" + ` : "Set Background", "Skip Background");

            updateSummaryRow("sumDocker", document.getElementById("installDocker").checked, "Install Docker Rootless", "Skip Docker");
            updateSummaryRow("sumShell", document.getElementById("defaultShellZsh").checked, "Set Zsh as default", "Keep Bash");
        }

        function updateSummaryRow(elementId, enabled, activeText, inactiveText) {
            const container = document.getElementById(elementId);
            const indicator = container.querySelector(".indicator");
            const textSpan = container.querySelector("span:last-child");

            indicator.className = "indicator " + (enabled ? "enabled" : "disabled");
            textSpan.textContent = enabled ? activeText : inactiveText;
        }

        async function submitConfig() {
            // Read options
            const installZsh = document.getElementById("installZsh").checked;
            const configureGit = document.getElementById("setupGit").checked;
            const gitName = document.getElementById("gitName").value;
            const gitEmail = document.getElementById("gitEmail").value;
            const applyTheme = document.getElementById("setupTheme").checked;
            const themeName = document.getElementById("themeSelect").value;
            const applyBackground = document.getElementById("setupBg").checked;

            let backgroundImage = "";
            const bgSrc = document.getElementById("bgSource").value;
            if (bgSrc === "2") {
                backgroundImage = document.getElementById("repoWpSelect").value; // we'll join it in Go
            } else if (bgSrc === "3") {
                backgroundImage = document.getElementById("customBgPath").value;
            }

            const enableDocker = document.getElementById("installDocker").checked;
            const enableZshDefault = document.getElementById("defaultShellZsh").checked;

            const payload = {
                install_oh_my_zsh: installZsh,
                configure_git: configureGit,
                git_name: gitName,
                git_email: gitEmail,
                apply_theme: applyTheme,
                theme_mode: selectedThemeMode,
                theme_name: themeName,
                apply_background: applyBackground,
                background_image: backgroundImage,
                enable_docker: enableDocker,
                enable_zsh_default: enableZshDefault
            };

            // Switch to console view
            document.getElementById("progressBarContainer").style.display = "none";
            document.getElementById("step6").style.display = "none";
            document.getElementById("btnRow").style.display = "none";
            document.getElementById("pulseContainer").style.display = "flex";
            document.getElementById("consoleCard").classList.add("active");

            try {
                // Post config
                const res = await fetch('/api/apply', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(payload)
                });

                if (!res.ok) {
                    throw new Error("Failed to post configuration payload.");
                }

                // Listen to SSE Stream
                const consoleDiv = document.getElementById("consoleCard");
                const exportSettings = document.getElementById("exportSettingsWeb").checked;
                const eventSource = new EventSource(` + "`" + `/api/stream?export=\${exportSettings}` + "`" + `);

                eventSource.onmessage = function(event) {
                    const line = document.createElement("div");
                    line.className = "console-line";
                    line.textContent = event.data;
                    consoleDiv.appendChild(line);
                    consoleDiv.scrollTop = consoleDiv.scrollHeight;
                };

                eventSource.onerror = function() {
                    // Stream closed or error
                    eventSource.close();
                    document.getElementById("pulseContainer").style.display = "none";
                    document.getElementById("consoleCard").style.borderColor = "var(--success-color)";
                    
                    // Show finish view
                    document.getElementById("finishView").classList.add("active");

                    // Trigger browser file download if export settings is checked
                    if (exportSettings) {
                        const dataStr = "data:text/json;charset=utf-8," + encodeURIComponent(JSON.stringify(payload, null, 2));
                        const downloadAnchor = document.createElement('a');
                        downloadAnchor.setAttribute("href", dataStr);
                        downloadAnchor.setAttribute("download", "config.json");
                        document.body.appendChild(downloadAnchor);
                        downloadAnchor.click();
                        downloadAnchor.remove();
                    }
                };

            } catch (err) {
                console.error(err);
                const consoleDiv = document.getElementById("consoleCard");
                const errLine = document.createElement("div");
                errLine.className = "console-line";
                errLine.style.color = "var(--error-color)";
                errLine.textContent = "Error occurred: " + err.message;
                consoleDiv.appendChild(errLine);
                document.getElementById("pulseContainer").style.display = "none";
            }
        }

        function showToast(message, isSuccess) {
            const toast = document.getElementById("toastNotification");
            toast.textContent = message;
            toast.style.display = "block";
            if (isSuccess) {
                toast.style.background = "rgba(16, 185, 129, 0.15)";
                toast.style.border = "1px solid var(--success-color)";
                toast.style.color = "#34d399";
            } else {
                toast.style.background = "rgba(239, 68, 68, 0.15)";
                toast.style.border = "1px solid var(--error-color)";
                toast.style.color = "#f87171";
            }
            setTimeout(() => {
                toast.style.display = "none";
            }, 3000);
        }

        async function manualImport() {
            try {
                const configRes = await fetch('/api/config');
                if (!configRes.ok) throw new Error("Failed to load config.");
                const configData = await configRes.json();
                importedSettings = true;
                applyConfigToUI(configData);
                showStep(1); // Force banner refresh
                showToast("Settings imported successfully!", true);
            } catch (err) {
                showToast("Failed to import settings: " + err.message, false);
            }
        }

        async function manualExport() {
            try {
                const installZsh = document.getElementById("installZsh").checked;
                const configureGit = document.getElementById("setupGit").checked;
                const gitName = document.getElementById("gitName").value;
                const gitEmail = document.getElementById("gitEmail").value;
                const applyTheme = document.getElementById("setupTheme").checked;
                const themeName = document.getElementById("themeSelect").value;
                const applyBackground = document.getElementById("setupBg").checked;

                let backgroundImage = "";
                const bgSrc = document.getElementById("bgSource").value;
                if (bgSrc === "2") {
                    backgroundImage = document.getElementById("repoWpSelect").value;
                } else if (bgSrc === "3") {
                    backgroundImage = document.getElementById("customBgPath").value;
                }

                const enableDocker = document.getElementById("installDocker").checked;
                const enableZshDefault = document.getElementById("defaultShellZsh").checked;

                const payload = {
                    install_oh_my_zsh: installZsh,
                    configure_git: configureGit,
                    git_name: gitName,
                    git_email: gitEmail,
                    apply_theme: applyTheme,
                    theme_mode: selectedThemeMode,
                    theme_name: themeName,
                    apply_background: applyBackground,
                    background_image: backgroundImage,
                    enable_docker: enableDocker,
                    enable_zsh_default: enableZshDefault
                };

                const res = await fetch('/api/export', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(payload)
                });

                if (res.ok) {
                    const dataStr = "data:text/json;charset=utf-8," + encodeURIComponent(JSON.stringify(payload, null, 2));
                    const downloadAnchor = document.createElement('a');
                    downloadAnchor.setAttribute("href", dataStr);
                    downloadAnchor.setAttribute("download", "config.json");
                    document.body.appendChild(downloadAnchor);
                    downloadAnchor.click();
                    downloadAnchor.remove();

                    showToast("Settings exported and config.json downloaded successfully!", true);
                } else {
                    throw new Error("Failed to write settings.");
                }
            } catch (err) {
                showToast("Failed to export settings: " + err.message, false);
            }
        }

        async function restartTerminal() {
            try {
                await fetch('/api/restart', { method: 'POST' });
            } catch (err) {
                console.error("Restart terminal request completed with standard connection close.");
            }
        }

        // Init
        document.getElementById("btnPrev").style.visibility = "hidden";
        toggleGitFields();
        toggleBgInputs();
        fetchSystemResources();
    </script>
</body>
</html>
`
