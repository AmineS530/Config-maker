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
            --bg-color: hsl(222, 47%, 6%);
            --card-bg: rgba(17, 24, 39, 0.7);
            --border-color: rgba(255, 255, 255, 0.08);
            --primary-accent: hsl(199, 100%, 50%);
            --secondary-accent: hsl(271, 91%, 65%);
            --text-main: hsl(0, 0%, 95%);
            --text-muted: hsl(215, 15%, 65%);
            --success-color: hsl(142, 71%, 45%);
            --error-color: hsl(350, 89%, 60%);
            --warning-color: hsl(37, 90%, 50%);
        }

        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }

        body {
            font-family: 'Inter', sans-serif;
            background: linear-gradient(135deg, var(--bg-color) 0%, hsl(223, 47%, 12%) 100%);
            color: var(--text-main);
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            overflow-x: hidden;
            position: relative;
        }

        /* Abstract glowing blobs for premium feel */
        .glow-blob {
            position: absolute;
            width: 500px;
            height: 500px;
            border-radius: 50%;
            background: radial-gradient(circle, rgba(14, 165, 233, 0.15) 0%, rgba(139, 92, 246, 0.05) 50%, rgba(0,0,0,0) 100%);
            filter: blur(80px);
            z-index: 0;
            pointer-events: none;
        }
        .glow-blob-1 { top: -100px; left: -100px; }
        .glow-blob-2 { bottom: -100px; right: -100px; }

        .container {
            width: 100%;
            max-width: 720px;
            padding: 24px;
            z-index: 10;
        }

        /* Header design */
        header {
            text-align: center;
            margin-bottom: 32px;
        }
        header h1 {
            font-family: 'Outfit', sans-serif;
            font-size: 2.8rem;
            font-weight: 800;
            background: linear-gradient(to right, var(--primary-accent), var(--secondary-accent));
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            letter-spacing: -0.5px;
            margin-bottom: 8px;
            filter: drop-shadow(0 0 10px rgba(14, 165, 233, 0.3));
        }
        header p {
            color: var(--text-muted);
            font-size: 1.05rem;
            font-weight: 400;
        }

        /* Card design system */
        .glass-card {
            background: var(--card-bg);
            border: 1px solid var(--border-color);
            border-radius: 24px;
            padding: 40px;
            backdrop-filter: blur(16px);
            box-shadow: 0 20px 50px rgba(0, 0, 0, 0.3);
            transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
            position: relative;
            overflow: hidden;
        }
        .glass-card::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            height: 3px;
            background: linear-gradient(90deg, var(--primary-accent), var(--secondary-accent));
            opacity: 0.8;
        }

        /* Progress Bar styles */
        .progress-container {
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin-bottom: 32px;
            background: rgba(255, 255, 255, 0.03);
            padding: 12px 24px;
            border-radius: 16px;
            border: 1px solid rgba(255, 255, 255, 0.02);
        }
        .progress-bar-wrapper {
            flex-grow: 1;
            height: 6px;
            background: rgba(255, 255, 255, 0.08);
            border-radius: 3px;
            margin: 0 16px;
            overflow: hidden;
            position: relative;
        }
        .progress-bar-fill {
            height: 100%;
            width: 0%;
            background: linear-gradient(90deg, var(--primary-accent), var(--secondary-accent));
            border-radius: 3px;
            transition: width 0.4s cubic-bezier(0.16, 1, 0.3, 1);
        }
        .progress-step-text {
            font-size: 0.85rem;
            color: var(--text-muted);
            font-weight: 500;
            min-width: 65px;
        }

        /* Wizard Step content classes */
        .step-content {
            display: none;
            animation: fadeInSlide 0.5s cubic-bezier(0.16, 1, 0.3, 1) forwards;
        }
        .step-content.active {
            display: block;
        }

        @keyframes fadeInSlide {
            from {
                opacity: 0;
                transform: translateY(12px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }

        h2 {
            font-family: 'Outfit', sans-serif;
            font-size: 1.6rem;
            font-weight: 600;
            margin-bottom: 24px;
            color: var(--text-main);
        }

        /* Inputs & Interactive controls */
        .form-group {
            margin-bottom: 24px;
        }
        .form-group label {
            display: block;
            font-size: 0.9rem;
            color: var(--text-muted);
            margin-bottom: 8px;
            font-weight: 500;
        }
        .text-input {
            width: 100%;
            background: rgba(255, 255, 255, 0.04);
            border: 1px solid rgba(255, 255, 255, 0.1);
            border-radius: 12px;
            padding: 14px 18px;
            color: var(--text-main);
            font-family: inherit;
            font-size: 1rem;
            transition: all 0.3s ease;
        }
        .text-input:focus {
            outline: none;
            border-color: var(--primary-accent);
            background: rgba(255, 255, 255, 0.07);
            box-shadow: 0 0 12px rgba(14, 165, 233, 0.15);
        }

        /* Custom Toggle Switch */
        .toggle-card {
            background: rgba(255, 255, 255, 0.02);
            border: 1px solid rgba(255, 255, 255, 0.04);
            border-radius: 16px;
            padding: 20px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin-bottom: 20px;
            transition: all 0.3s ease;
        }
        .toggle-card:hover {
            border-color: rgba(255, 255, 255, 0.08);
            background: rgba(255, 255, 255, 0.03);
            transform: translateY(-2px);
        }
        .toggle-info {
            max-width: 80%;
        }
        .toggle-title {
            font-weight: 600;
            font-size: 1.05rem;
            margin-bottom: 4px;
        }
        .toggle-desc {
            font-size: 0.85rem;
            color: var(--text-muted);
            line-height: 1.4;
        }
        .switch {
            position: relative;
            display: inline-block;
            width: 52px;
            height: 28px;
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
            background-color: rgba(255, 255, 255, 0.1);
            transition: .3s cubic-bezier(0.16, 1, 0.3, 1);
            border-radius: 34px;
            border: 1px solid rgba(255, 255, 255, 0.05);
        }
        .slider:before {
            position: absolute;
            content: "";
            height: 20px;
            width: 20px;
            left: 3px;
            bottom: 3px;
            background-color: white;
            transition: .3s cubic-bezier(0.16, 1, 0.3, 1);
            border-radius: 50%;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
        }
        input:checked + .slider {
            background-image: linear-gradient(to right, var(--primary-accent), var(--secondary-accent));
            border-color: transparent;
        }
        input:checked + .slider:before {
            transform: translateX(24px);
        }

        /* Choice options grid (e.g. Mode selections) */
        .options-grid {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 16px;
            margin-bottom: 24px;
        }
        .choice-box {
            background: rgba(255, 255, 255, 0.02);
            border: 1px solid rgba(255, 255, 255, 0.05);
            border-radius: 16px;
            padding: 20px;
            text-align: center;
            cursor: pointer;
            transition: all 0.3s ease;
        }
        .choice-box:hover {
            border-color: rgba(255, 255, 255, 0.1);
            background: rgba(255, 255, 255, 0.03);
        }
        .choice-box.selected {
            background: rgba(14, 165, 233, 0.08);
            border-color: var(--primary-accent);
            box-shadow: 0 0 16px rgba(14, 165, 233, 0.1);
        }
        .choice-icon {
            font-size: 1.8rem;
            margin-bottom: 8px;
        }
        .choice-title {
            font-weight: 600;
            font-size: 1rem;
        }

        .select-input {
            width: 100%;
            background: hsl(222, 47%, 9%);
            border: 1px solid rgba(255, 255, 255, 0.1);
            border-radius: 12px;
            padding: 14px 18px;
            color: var(--text-main);
            font-family: inherit;
            font-size: 1rem;
            cursor: pointer;
        }

        /* Summary item layout */
        .summary-list {
            background: rgba(255, 255, 255, 0.02);
            border: 1px solid rgba(255, 255, 255, 0.04);
            border-radius: 16px;
            padding: 24px;
            margin-bottom: 24px;
        }
        .summary-row {
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 12px 0;
            border-bottom: 1px solid rgba(255, 255, 255, 0.05);
        }
        .summary-row:last-child {
            border-bottom: none;
        }
        .summary-label {
            font-weight: 500;
            color: var(--text-muted);
        }
        .summary-value {
            display: flex;
            align-items: center;
            font-weight: 600;
        }
        .indicator {
            width: 16px;
            height: 16px;
            border-radius: 50%;
            display: inline-block;
            margin-right: 8px;
        }
        .indicator.enabled {
            background-color: var(--success-color);
            box-shadow: 0 0 8px var(--success-color);
        }
        .indicator.disabled {
            background-color: var(--error-color);
            box-shadow: 0 0 8px var(--error-color);
        }

        /* Buttons styles */
        .btn-row {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-top: 32px;
        }
        .btn {
            font-family: 'Outfit', sans-serif;
            font-size: 1rem;
            font-weight: 600;
            padding: 14px 28px;
            border-radius: 14px;
            cursor: pointer;
            transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
            border: none;
            display: inline-flex;
            align-items: center;
            justify-content: center;
            text-decoration: none;
        }
        .btn-prev {
            background: rgba(255, 255, 255, 0.05);
            color: var(--text-main);
            border: 1px solid rgba(255, 255, 255, 0.08);
        }
        .btn-prev:hover {
            background: rgba(255, 255, 255, 0.1);
        }
        .btn-next {
            background: linear-gradient(90deg, var(--primary-accent), var(--secondary-accent));
            color: white;
            box-shadow: 0 4px 20px rgba(14, 165, 233, 0.25);
        }
        .btn-next:hover {
            transform: translateY(-2px);
            box-shadow: 0 8px 24px rgba(14, 165, 233, 0.4);
        }

        /* Console Output styles */
        .console-card {
            display: none;
            background: hsl(224, 71%, 4%);
            border: 1px solid var(--border-color);
            border-radius: 20px;
            padding: 24px;
            font-family: 'Courier New', Courier, monospace;
            height: 380px;
            overflow-y: auto;
            color: #10b981;
            margin-bottom: 24px;
            box-shadow: inset 0 4px 24px rgba(0,0,0,0.8);
        }
        .console-card.active {
            display: block;
        }
        .console-line {
            margin-bottom: 8px;
            white-space: pre-wrap;
            line-height: 1.4;
            font-size: 0.9rem;
        }
        .console-pulse-container {
            display: flex;
            align-items: center;
            margin-bottom: 16px;
            color: var(--text-muted);
            font-size: 0.85rem;
        }
        .pulse-dot {
            width: 8px;
            height: 8px;
            border-radius: 50%;
            background-color: var(--success-color);
            margin-right: 8px;
            animation: pulse 1.5s infinite;
        }
        @keyframes pulse {
            0% { transform: scale(0.9); opacity: 0.5; }
            50% { transform: scale(1.2); opacity: 1; box-shadow: 0 0 10px var(--success-color); }
            100% { transform: scale(0.9); opacity: 0.5; }
        }

        /* Finish View Styles */
        .finish-view {
            display: none;
            text-align: center;
            padding: 20px 0;
        }
        .finish-view.active {
            display: block;
        }
        .finish-icon {
            font-size: 4rem;
            margin-bottom: 16px;
            animation: scaleBounce 0.6s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
        }
        @keyframes scaleBounce {
            0% { transform: scale(0.3); opacity: 0; }
            100% { transform: scale(1); opacity: 1; }
        }
    </style>
</head>
<body>
    <div class="glow-blob glow-blob-1"></div>
    <div class="glow-blob glow-blob-2"></div>

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
            } catch (err) {
                console.error("Failed to load resources:", err);
            }
        }

        function applyConfigToUI(cfg) {
            if (!cfg) return;

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
                const eventSource = new EventSource("/api/stream");

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
                applyConfigToUI(configData);
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
                    showToast("Settings exported successfully!", true);
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
