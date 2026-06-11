package web

// IndexTemplate is the embedded HTML/CSS/JS template for our local web wizard.
const IndexTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>ZoneRestore Dashboard</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600&family=Outfit:wght@400;600;800&display=swap" rel="stylesheet">
    <script src="/js/alpine.min.js" defer></script>
    <style>
        :root {
            --bg-color: #030303;      /* Dark obsidian background */
            --card-bg: rgba(18, 18, 22, 0.7); /* Translucent premium card */
            --border-color: rgba(255, 255, 255, 0.08); /* Minimal soft border */
            --primary-accent: #6366f1; /* Indigo violet accent */
            --secondary-accent: #3b82f6; /* Royal blue secondary */
            --text-main: #f8fafc;     /* Bright off-white */
            --text-muted: #94a3b8;    /* Muted cool gray */
            --success-color: #10b981; /* Emerald green */
            --error-color: #f43f5e;   /* Soft crimson red */
            --warning-color: #f59e0b; /* Warm amber */
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

        /* Ambient Glow Blobs */
        .glow-bg {
            position: fixed;
            border-radius: 50%;
            filter: blur(160px);
            opacity: 0.12;
            z-index: 1;
            pointer-events: none;
            transition: all 1s ease;
        }
        .glow-1 {
            width: 450px;
            height: 450px;
            background: radial-gradient(circle, var(--primary-accent), transparent 70%);
            top: -120px;
            left: -120px;
        }
        .glow-2 {
            width: 450px;
            height: 450px;
            background: radial-gradient(circle, var(--secondary-accent), transparent 70%);
            bottom: -120px;
            right: -120px;
        }

        .container {
            width: 100%;
            max-width: 680px;
            padding: 40px 24px;
            z-index: 10;
        }

        /* Header Styling */
        header {
            text-align: center;
            margin-bottom: 36px;
        }
        header h1 {
            font-family: 'Outfit', sans-serif;
            font-size: 2.6rem;
            font-weight: 800;
            letter-spacing: -1.2px;
            margin-bottom: 8px;
            background: linear-gradient(135deg, #ffffff 40%, #a5b4fc 70%, #3b82f6 100%);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
        }
        header p {
            color: var(--text-muted);
            font-size: 0.95rem;
            font-weight: 400;
        }

        /* Glassmorphism Card */
        .glass-card {
            background: var(--card-bg);
            border: 1px solid var(--border-color);
            border-radius: 20px;
            padding: 36px;
            box-shadow: 0 30px 60px rgba(0, 0, 0, 0.6), inset 0 1px 0 rgba(255, 255, 255, 0.05);
            backdrop-filter: blur(24px);
            -webkit-backdrop-filter: blur(24px);
            transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
            position: relative;
        }

        /* Stepper Progress Bar */
        .progress-container {
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin-bottom: 32px;
        }
        .progress-bar-wrapper {
            flex-grow: 1;
            height: 5px;
            background: rgba(255, 255, 255, 0.05);
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
            font-size: 0.8rem;
            color: var(--text-muted);
            font-weight: 600;
            letter-spacing: 0.5px;
            text-transform: uppercase;
            min-width: 75px;
        }

        /* Fade-in transitions */
        .step-content {
            display: none;
            animation: slideUpFade 0.35s cubic-bezier(0.16, 1, 0.3, 1) forwards;
        }
        .step-content.active {
            display: block;
        }

        @keyframes slideUpFade {
            from { opacity: 0; transform: translateY(12px); }
            to { opacity: 1; transform: translateY(0); }
        }

        h2 {
            font-family: 'Outfit', sans-serif;
            font-size: 1.4rem;
            font-weight: 600;
            margin-bottom: 24px;
            color: var(--text-main);
            letter-spacing: -0.4px;
        }

        /* Form Controls */
        .form-group {
            margin-bottom: 24px;
        }
        .form-group label {
            display: block;
            font-size: 0.85rem;
            color: var(--text-muted);
            margin-bottom: 8px;
            font-weight: 600;
            letter-spacing: 0.2px;
        }
        .text-input {
            width: 100%;
            background: rgba(0, 0, 0, 0.3);
            border: 1px solid var(--border-color);
            border-radius: 8px;
            padding: 12px 16px;
            color: var(--text-main);
            font-family: inherit;
            font-size: 0.95rem;
            transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
        }
        .text-input:focus {
            outline: none;
            border-color: var(--primary-accent);
            background: rgba(0, 0, 0, 0.5);
            box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.15);
        }

        /* Toggle Switches */
        .toggle-card {
            background: rgba(255, 255, 255, 0.01);
            border: 1px solid var(--border-color);
            border-radius: 12px;
            padding: 18px 20px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin-bottom: 16px;
            transition: all 0.25s cubic-bezier(0.16, 1, 0.3, 1);
        }
        .toggle-card:hover {
            border-color: rgba(255, 255, 255, 0.15);
            background: rgba(255, 255, 255, 0.03);
            transform: translateY(-1px);
        }
        .toggle-info {
            max-width: 80%;
        }
        .toggle-title {
            font-weight: 600;
            font-size: 0.95rem;
            margin-bottom: 3px;
        }
        .toggle-desc {
            font-size: 0.8rem;
            color: var(--text-muted);
            line-height: 1.45;
        }
        
        .switch {
            position: relative;
            display: inline-block;
            width: 46px;
            height: 26px;
            flex-shrink: 0;
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
            transition: .25s cubic-bezier(0.16, 1, 0.3, 1);
            border-radius: 26px;
            border: 1px solid rgba(255, 255, 255, 0.05);
        }
        .slider:before {
            position: absolute;
            content: "";
            height: 18px;
            width: 18px;
            left: 3px;
            bottom: 3px;
            background-color: #fff;
            transition: .25s cubic-bezier(0.16, 1, 0.3, 1);
            border-radius: 50%;
            box-shadow: 0 2px 4px rgba(0,0,0,0.2);
        }
        input:checked + .slider {
            background-color: var(--primary-accent);
            border-color: rgba(99, 102, 241, 0.2);
        }
        input:checked + .slider:before {
            transform: translateX(20px);
        }

        /* Choice Boxes for Grid Mode */
        .options-grid {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 16px;
            margin-bottom: 24px;
        }
        .choice-box {
            background: rgba(255, 255, 255, 0.01);
            border: 1px solid var(--border-color);
            border-radius: 12px;
            padding: 20px;
            text-align: center;
            cursor: pointer;
            transition: all 0.25s cubic-bezier(0.16, 1, 0.3, 1);
        }
        .choice-box:hover {
            border-color: rgba(255, 255, 255, 0.15);
            background: rgba(255, 255, 255, 0.03);
            transform: translateY(-2px);
        }
        .choice-box.selected {
            background: rgba(99, 102, 241, 0.1);
            border-color: var(--primary-accent);
            box-shadow: 0 0 25px rgba(99, 102, 241, 0.15);
        }
        .choice-icon {
            font-size: 1.6rem;
            margin-bottom: 8px;
        }
        .choice-title {
            font-weight: 600;
            font-size: 0.9rem;
        }

        .select-input {
            width: 100%;
            background: rgba(0, 0, 0, 0.3);
            border: 1px solid var(--border-color);
            border-radius: 8px;
            padding: 12px 16px;
            color: var(--text-main);
            font-family: inherit;
            font-size: 0.95rem;
            cursor: pointer;
            outline: none;
            transition: all 0.2s ease;
        }
        .select-input:focus {
            border-color: var(--primary-accent);
            box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.15);
        }

        /* Summary Setup Overview */
        .summary-list {
            background: rgba(255, 255, 255, 0.01);
            border: 1px solid var(--border-color);
            border-radius: 12px;
            padding: 20px;
            margin-bottom: 24px;
        }
        .summary-row {
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 12px 0;
            border-bottom: 1px solid rgba(255, 255, 255, 0.03);
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
            width: 8px;
            height: 8px;
            border-radius: 50%;
            display: inline-block;
            margin-right: 8px;
        }
        .indicator.enabled {
            background-color: var(--success-color);
            box-shadow: 0 0 8px var(--success-color);
        }
        .indicator.disabled {
            background-color: rgba(255,255,255,0.2);
        }

        /* Action Buttons */
        .btn-row {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-top: 32px;
        }
        .btn {
            font-family: inherit;
            font-size: 0.9rem;
            font-weight: 600;
            padding: 12px 26px;
            border-radius: 8px;
            cursor: pointer;
            transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
            border: none;
            display: inline-flex;
            align-items: center;
            justify-content: center;
            text-decoration: none;
        }
        .btn-prev {
            background: rgba(255, 255, 255, 0.03);
            color: var(--text-main);
            border: 1px solid var(--border-color);
        }
        .btn-prev:hover {
            background: rgba(255, 255, 255, 0.07);
            border-color: rgba(255, 255, 255, 0.2);
            transform: translateY(-1px);
        }
        .btn-next {
            background: linear-gradient(135deg, var(--primary-accent), var(--secondary-accent));
            color: white;
            box-shadow: 0 4px 15px rgba(99, 102, 241, 0.25);
        }
        .btn-next:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 20px rgba(99, 102, 241, 0.4);
            filter: brightness(1.1);
        }
        .btn-next:active {
            transform: translateY(0);
        }

        /* Floating Toast Alert Box */
        #toastNotification {
            position: fixed;
            top: 24px;
            right: 24px;
            z-index: 10000;
            max-width: 380px;
            padding: 16px 20px;
            border-radius: 12px;
            box-shadow: 0 16px 36px rgba(0, 0, 0, 0.5);
            font-size: 0.9rem;
            font-weight: 600;
            display: none;
            backdrop-filter: blur(12px);
            -webkit-backdrop-filter: blur(12px);
            animation: slideIn 0.35s cubic-bezier(0.16, 1, 0.3, 1) forwards;
        }
        @keyframes slideIn {
            from { transform: translateX(120%) translateY(-10px); opacity: 0; }
            to { transform: translateX(0) translateY(0); opacity: 1; }
        }

        /* macOS Console Style Terminal */
        .console-card {
            background: #09090b;
            padding: 20px;
            font-family: 'SFMono-Regular', Consolas, "Liberation Mono", Menlo, Courier, monospace;
            height: 340px;
            overflow-y: auto;
            color: #e4e4e7;
            box-shadow: inset 0 2px 10px rgba(0,0,0,0.8);
        }
        .console-line {
            margin-bottom: 6px;
            white-space: pre-wrap;
            line-height: 1.5;
            font-size: 0.85rem;
        }
        .console-pulse-container {
            display: flex;
            align-items: center;
            color: var(--text-muted);
            font-size: 0.85rem;
        }
        .pulse-dot {
            width: 8px;
            height: 8px;
            border-radius: 50%;
            background-color: var(--primary-accent);
            margin-right: 10px;
            animation: pulse 1.2s infinite;
            box-shadow: 0 0 6px var(--primary-accent);
        }
        @keyframes pulse {
            0% { opacity: 0.3; }
            50% { opacity: 1; }
            100% { opacity: 0.3; }
        }

        /* Modal Overlays & Cards */
        .modal-overlay {
            display: none;
            position: fixed;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background: rgba(3, 3, 3, 0.75);
            backdrop-filter: blur(18px);
            -webkit-backdrop-filter: blur(18px);
            z-index: 1000;
            align-items: center;
            justify-content: center;
            opacity: 0;
            transition: opacity 0.3s cubic-bezier(0.16, 1, 0.3, 1);
        }
        .modal-card {
            background: rgba(18, 18, 22, 0.8);
            border: 1px solid var(--border-color);
            border-radius: 24px;
            padding: 40px;
            max-width: 500px;
            width: 90%;
            text-align: center;
            box-shadow: 0 30px 60px rgba(0, 0, 0, 0.6), inset 0 1px 0 rgba(255,255,255,0.05);
            transform: scale(0.95);
            transition: transform 0.3s cubic-bezier(0.16, 1, 0.3, 1);
            backdrop-filter: blur(24px);
            -webkit-backdrop-filter: blur(24px);
        }

        /* Clean Finish View */
        .finish-view {
            display: none;
            text-align: center;
            padding: 24px 0;
        }
        .finish-view.active {
            display: block;
        }
        .finish-icon {
            font-size: 3.5rem;
            margin-bottom: 16px;
            filter: drop-shadow(0 0 10px rgba(99,102,241,0.3));
        }

        /* Custom Alias item card design */
        .alias-item {
            background: rgba(255, 255, 255, 0.02);
            border: 1px solid var(--border-color);
            border-radius: 10px;
            padding: 12px 16px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 16px;
            transition: all 0.2s ease;
        }
        .alias-item:hover {
            background: rgba(255, 255, 255, 0.04);
            border-color: rgba(255,255,255,0.15);
        }
    </style>
</head>
<body>
    <!-- Ambient Blur Background Blobs -->
    <div class="glow-bg glow-1"></div>
    <div class="glow-bg glow-2"></div>

    <!-- Floating Toast Notifications -->
    <div id="toastNotification"></div>

    <div class="container">
        <header>
            <h1>ZoneRestore</h1>
            <p>Desktop Configuration Wizard</p>
            <div style="margin-top: 20px; display: flex; justify-content: center; gap: 16px;">
                <button class="btn btn-prev" onclick="manualImport()" style="padding: 10px 22px; font-size: 0.85rem; border-radius: 10px; margin-top: 0; background: rgba(99, 102, 241, 0.1); border-color: rgba(99, 102, 241, 0.25); color: #a5b4fc;">📥 Import Settings</button>
            </div>
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
                        <input type="text" id="gitName" class="text-input" placeholder="e.g. 3elal">
                    </div>
                    <div class="form-group">
                        <label for="gitEmail">Email Address</label>
                        <input type="email" id="gitEmail" class="text-input" placeholder="e.g. 3elal@example.com">
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
                            <option value="3">Select custom image via GNOME</option>
                        </select>
                    </div>
                    <div class="form-group" id="repoWallpapersWrapper" style="display: none;">
                        <label for="repoWpSelect">Available Wallpaper</label>
                        <select id="repoWpSelect" class="select-input" onchange="updateWallpaperPreview()">
                            <!-- Populated from system values -->
                        </select>
                    </div>
                    <div class="form-group" id="customPathWrapper" style="display: none;">
                        <label>Custom Selected Image</label>
                        <div style="display: flex; gap: 12px; align-items: center; margin-bottom: 12px;">
                            <button type="button" class="btn" onclick="selectGnomeImage()" style="background: rgba(139, 92, 246, 0.2); border: 1px solid var(--secondary-accent); color: var(--text-light); padding: 8px 16px; border-radius: 8px; cursor: pointer; font-size: 0.9rem; transition: background 0.2s;">
                                📂 Choose Image...
                            </button>
                            <span id="selectedImagePath" style="font-size: 0.85rem; color: var(--text-muted); word-break: break-all;">No image selected</span>
                        </div>
                    </div>
                    <div id="wallpaperPreviewWrapper" style="margin-top: 20px; display: none;">
                        <label style="display: block; margin-bottom: 8px; font-weight: 500; font-size: 0.9rem; color: var(--text-muted);">Wallpaper Preview</label>
                        <div class="preview-card" style="border: 1px solid rgba(255, 255, 255, 0.1); border-radius: 16px; overflow: hidden; background: rgba(255, 255, 255, 0.03); aspect-ratio: 16/9; display: flex; align-items: center; justify-content: center; position: relative; box-shadow: inset 0 0 20px rgba(0,0,0,0.4);">
                            <img id="wallpaperPreviewImg" src="" style="width: 100%; height: 100%; object-fit: cover; display: none; transition: opacity 0.3s ease;">
                            <div id="wallpaperPreviewPlaceholder" style="color: var(--text-muted); font-size: 0.9rem; display: flex; flex-direction: column; align-items: center; gap: 8px;">
                                <span>🖼️</span>
                                <span>No preview available</span>
                            </div>
                        </div>
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
                <div class="toggle-card">
                    <div class="toggle-info">
                        <div class="toggle-title">Configure Keyboard Layouts</div>
                        <div class="toggle-desc">Applies US and French keyboard layouts in Gnome.</div>
                    </div>
                    <label class="switch">
                        <input type="checkbox" id="configureKeyboard" checked>
                        <span class="slider"></span>
                    </label>
                </div>
                <div class="toggle-card">
                    <div class="toggle-info">
                        <div class="toggle-title">Configure GNOME Power Settings</div>
                        <div class="toggle-desc">Sets desktop logout/idle sleep timeout to 1.5 hours.</div>
                    </div>
                    <label class="switch">
                        <input type="checkbox" id="configurePower" checked>
                        <span class="slider"></span>
                    </label>
                </div>
                <div class="toggle-card">
                    <div class="toggle-info">
                        <div class="toggle-title">Install Custom Fonts</div>
                        <div class="toggle-desc">Installs custom MPLUS/Meslo fonts and sets terminal default font.</div>
                    </div>
                    <label class="switch">
                        <input type="checkbox" id="configureFonts" checked>
                        <span class="slider"></span>
                    </label>
                </div>

                <!-- Alpine JS Shell customization wrapper -->
                <div x-data="{
                    customUsername: '',
                    aliases: [],
                    newAliasName: '',
                    newAliasCommand: '',
                    init() {
                        document.addEventListener('load-aliases', (e) => {
                            this.customUsername = e.detail.customUsername || '';
                            this.aliases = e.detail.aliases || [];
                        });
                    },
                    addAlias() {
                        if (this.newAliasName.trim() && this.newAliasCommand.trim()) {
                            this.aliases.push({
                                name: this.newAliasName.trim(),
                                command: this.newAliasCommand.trim(),
                                enabled: true
                            });
                            this.newAliasName = '';
                            this.newAliasCommand = '';
                        }
                    },
                    removeAlias(index) {
                        this.aliases.splice(index, 1);
                    }
                }" @get-aliases.window="$event.detail.callback({ customUsername: customUsername, aliases: aliases })" style="margin-top: 24px; padding-top: 24px; border-top: 1px dashed var(--border-color);">
                    
                    <!-- Custom Username for Prompt (PS1) -->
                    <div class="form-group">
                        <label for="customUsernameInput" style="font-weight: 600; color: var(--text-main); display: block; margin-bottom: 8px;">Prompt Display Name (PS1)</label>
                        <input type="text" id="customUsernameInput" class="text-input" x-model="customUsername" placeholder="e.g. myname" style="width: 100%; max-width: 100%;">
                        <div style="font-size: 0.8rem; color: var(--text-muted); margin-top: 4px;">Leaves segment hidden if empty. Defaults to system username.</div>
                    </div>

                    <!-- Zsh Aliases Section -->
                    <div style="margin-top: 28px;">
                        <label style="font-weight: 600; color: var(--text-main); display: block; margin-bottom: 12px;">Configure Zsh Aliases</label>
                        
                        <!-- Default & Custom Aliases List -->
                        <div style="display: flex; flex-direction: column; gap: 12px;">
                            <template x-for="(alias, index) in aliases" :key="index">
                                <div style="background: rgba(255, 255, 255, 0.02); border: 1px solid var(--border-color); border-radius: 12px; padding: 12px 16px; display: flex; align-items: center; justify-content: space-between; gap: 16px;">
                                    <div style="display: flex; align-items: center; gap: 12px; flex: 1;">
                                        <label class="switch" style="flex-shrink: 0;">
                                            <input type="checkbox" x-model="alias.enabled">
                                            <span class="slider"></span>
                                        </label>
                                        <div style="word-break: break-all;">
                                            <strong style="color: var(--primary-accent); font-family: monospace; font-size: 0.95rem;" x-text="alias.name"></strong>
                                            <span style="color: var(--text-muted); font-size: 0.85rem; margin-left: 8px;" x-text="'= ' + alias.command"></span>
                                        </div>
                                    </div>
                                    <button type="button" @click="removeAlias(index)" style="background: none; border: none; color: var(--error-color); cursor: pointer; font-size: 1rem; padding: 4px 8px; border-radius: 6px; transition: background 0.2s;" onmouseover="this.style.background='rgba(239, 68, 68, 0.1)'" onmouseout="this.style.background='none'">
                                        ❌
                                    </button>
                                </div>
                            </template>
                        </div>

                        <!-- Add New Alias Form -->
                        <div style="margin-top: 16px; background: rgba(255, 255, 255, 0.01); border: 1px dashed var(--border-color); border-radius: 16px; padding: 16px; display: flex; flex-direction: column; gap: 12px;">
                            <div style="font-size: 0.85rem; font-weight: 500; color: var(--text-muted);">➕ Add Custom Alias</div>
                            <div style="display: flex; gap: 12px; flex-wrap: wrap;">
                                <input type="text" x-model="newAliasName" placeholder="alias name (e.g. gp)" class="text-input" style="flex: 1; min-width: 150px; font-family: monospace; font-size: 0.9rem; margin-top: 0;">
                                <input type="text" x-model="newAliasCommand" placeholder="command (e.g. git push)" class="text-input" style="flex: 2; min-width: 200px; font-family: monospace; font-size: 0.9rem; margin-top: 0;">
                            </div>
                            <button type="button" @click="addAlias()" class="btn" style="align-self: flex-start; background: rgba(139, 92, 246, 0.2); border: 1px solid var(--secondary-accent); color: var(--text-light); padding: 8px 16px; border-radius: 8px; cursor: pointer; font-size: 0.85rem; transition: background 0.2s;">
                                Add Alias
                            </button>
                        </div>
                    </div>
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
                    <div class="summary-row">
                        <span class="summary-label">Configure Keyboard Layouts</span>
                        <span class="summary-value" id="sumKeyboard"><span class="indicator"></span><span></span></span>
                    </div>
                    <div class="summary-row">
                        <span class="summary-label">Configure Power Settings</span>
                        <span class="summary-value" id="sumPower"><span class="indicator"></span><span></span></span>
                    </div>
                    <div class="summary-row">
                        <span class="summary-label">Install Custom Fonts</span>
                        <span class="summary-value" id="sumFonts"><span class="indicator"></span><span></span></span>
                    </div>
                    <div class="summary-row">
                        <span class="summary-label">Prompt Display Name</span>
                        <span class="summary-value" id="sumPromptName"><span></span></span>
                    </div>
                    <div class="summary-row">
                        <span class="summary-label">Active Zsh Aliases</span>
                        <span class="summary-value" id="sumAliases"><span></span></span>
                    </div>
                </div>
                <div class="toggle-card" style="margin-top: 24px;">
                    <div class="toggle-info" style="text-align: left;">
                        <div class="toggle-title">Save & Export Settings</div>
                        <div class="toggle-desc">Automatically save these choices to ~/.config/zonerestore/config.json.</div>
                    </div>
                    <label class="switch">
                        <input type="checkbox" id="exportSettingsWeb" checked>
                        <span class="slider"></span>
                    </label>
                </div>
            </div>

            <!-- Streaming Console View -->
            <div class="console-pulse-container" id="pulseContainer" style="display: none; margin-bottom: 16px; justify-content: center;">
                <div class="pulse-dot"></div>
                <span id="consoleStatusText" style="font-weight: 500;">Applying changes... Please wait.</span>
            </div>
            
            <div class="terminal-window" id="terminalWindow" style="display: none; margin-bottom: 24px; box-shadow: 0 20px 50px rgba(0,0,0,0.6); border-radius: 12px; overflow: hidden; border: 1px solid var(--border-color);">
                <div class="terminal-header" style="background: #141416; border-bottom: 1px solid var(--border-color); padding: 12px 16px; display: flex; align-items: center; justify-content: space-between;">
                    <div style="display: flex; gap: 8px;">
                        <span style="width: 12px; height: 12px; border-radius: 50%; background: #ef4444; display: inline-block;"></span>
                        <span style="width: 12px; height: 12px; border-radius: 50%; background: #f59e0b; display: inline-block;"></span>
                        <span style="width: 12px; height: 12px; border-radius: 50%; background: #10b981; display: inline-block;"></span>
                    </div>
                    <div style="font-size: 0.8rem; color: var(--text-muted); font-family: monospace; font-weight: 500;">zsh - ZoneRestore</div>
                    <div style="width: 52px;"></div>
                </div>
                <div class="console-card" id="consoleCard" style="border-radius: 0; border: none; margin-bottom: 0; height: 320px; display: block; background: #09090b;">
                    <!-- Exec logs shown live -->
                </div>
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

    <!-- Startup Dialog Modal -->
    <div id="startupModal" class="modal-overlay" style="display: none; position: fixed; top: 0; left: 0; width: 100%; height: 100%; background: rgba(9, 9, 11, 0.85); backdrop-filter: blur(12px); z-index: 1000; align-items: center; justify-content: center; opacity: 0; transition: opacity 0.3s ease;">
        <div class="modal-card" style="background: var(--card-bg); border: 1px solid var(--border-color); border-radius: 24px; padding: 40px; max-width: 500px; width: 90%; text-align: center; box-shadow: 0 20px 40px rgba(0, 0, 0, 0.5); transform: scale(0.95); transition: transform 0.3s ease, opacity 0.3s ease;">
            <div style="font-size: 3rem; margin-bottom: 20px;">⚙️</div>
            <h2 style="font-family: 'Outfit', sans-serif; font-size: 1.8rem; margin-bottom: 12px; color: var(--text-main);">Welcome to ZoneRestore</h2>
            <p style="color: var(--text-muted); font-size: 1.05rem; margin-bottom: 32px; line-height: 1.5;">Configure your shell, tools, and theme preferences to build a dream workspace.</p>
            <div style="display: flex; flex-direction: column; gap: 16px;">
                <button onclick="selectStartupOption('fresh')" style="background: linear-gradient(90deg, var(--primary-accent), var(--secondary-accent)); border: none; color: var(--text-main); padding: 16px; border-radius: 12px; font-weight: 600; cursor: pointer; font-size: 1rem; transition: transform 0.2s, box-shadow 0.2s; box-shadow: 0 4px 12px rgba(59, 130, 246, 0.35);">
                    🚀 First-Time Setup (Start Fresh)
                </button>
                <button onclick="selectStartupOption('import')" style="background: rgba(255, 255, 255, 0.04); border: 1px solid var(--border-color); color: var(--text-main); padding: 16px; border-radius: 12px; font-weight: 600; cursor: pointer; font-size: 1rem; transition: background 0.2s, transform 0.2s;">
                    💾 Import Saved Settings (config.json)
                </button>
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

                // Open the startup modal overlay
                const modal = document.getElementById("startupModal");
                modal.style.display = "flex";
                setTimeout(() => {
                    modal.style.opacity = "1";
                    modal.querySelector(".modal-card").style.transform = "scale(1)";
                }, 50);
            } catch (err) {
                console.error("Failed to load resources:", err);
            }
        }

        function applyConfigToUI(cfg, isImported) {
            if (!cfg) return;

            importedSettings = !!isImported;
            const banner = document.getElementById("importBanner");
            if (banner) {
                banner.style.display = importedSettings ? "flex" : "none";
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
            const isRepoWp = wallpapers.some(wp => bgImage.endsWith(wp));
            if (bgImage === "" || bgImage.endsWith("Background.jpeg")) {
                document.getElementById("bgSource").value = "1";
            } else if (isRepoWp || bgImage.includes("wallpapers/")) {
                document.getElementById("bgSource").value = "2";
                const parts = bgImage.split("/");
                const filename = parts[parts.length - 1];
                document.getElementById("repoWpSelect").value = filename;
            } else {
                document.getElementById("bgSource").value = "3";
                selectedGnomeImagePath = bgImage;
                document.getElementById("selectedImagePath").textContent = bgImage || "No image selected";
            }
            toggleBgInputs();
            toggleBgFields();
            updateWallpaperPreview();

            // Docker & Default Shell
            document.getElementById("installDocker").checked = cfg.enable_docker === true;
            document.getElementById("defaultShellZsh").checked = cfg.enable_zsh_default !== false;
            document.getElementById("configureKeyboard").checked = cfg.configure_keyboard !== false;
            document.getElementById("configurePower").checked = cfg.configure_power !== false;
            document.getElementById("configureFonts").checked = cfg.configure_fonts !== false;

            // Dispatch to Alpine
            document.dispatchEvent(new CustomEvent('load-aliases', { detail: {
                customUsername: cfg.custom_username,
                aliases: cfg.aliases
            }}));
        }

        async function selectStartupOption(option) {
            const modal = document.getElementById("startupModal");

            if (option === 'fresh') {
                try {
                    const res = await fetch('/api/config/default');
                    const defaultCfg = await res.json();
                    applyConfigToUI(defaultCfg, false);
                    showToast("Loaded fresh default settings.", true);

                    // Close modal
                    modal.style.opacity = "0";
                    modal.querySelector(".modal-card").style.transform = "scale(0.95)";
                    setTimeout(() => { modal.style.display = "none"; }, 300);
                } catch (err) {
                    showToast("Failed to load default configuration: " + err.message, false);
                }
            } else if (option === 'import') {
                try {
                    const res = await fetch('/api/config/import');
                    const data = await res.json();
                    if (data.status === 'success') {
                        applyConfigToUI(data.config, true);
                        showToast("Settings imported successfully!", true);

                        // Close modal
                        modal.style.opacity = "0";
                        modal.querySelector(".modal-card").style.transform = "scale(0.95)";
                        setTimeout(() => { modal.style.display = "none"; }, 300);
                    } else if (data.status === 'canceled') {
                        showToast("Import canceled", false);
                    } else if (data.status === 'error') {
                        showToast("Failed to import: " + data.message, false);
                    }
                } catch (err) {
                    showToast("Failed to load saved configuration: " + err.message, false);
                }
            }
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

        let selectedGnomeImagePath = "";

        async function selectGnomeImage() {
            try {
                const res = await fetch('/api/select-wallpaper');
                if (!res.ok) throw new Error("Server error");
                const data = await res.json();
                if (data.status === "success" && data.path) {
                    selectedGnomeImagePath = data.path;
                    document.getElementById("selectedImagePath").textContent = data.path;
                    updateWallpaperPreview();
                } else if (data.status === "canceled") {
                    showToast("Selection canceled", false);
                }
            } catch (err) {
                showToast("Failed to select image: " + err.message, false);
            }
        }

        function updateWallpaperPreview() {
            const bgEnabled = document.getElementById("setupBg").checked;
            const previewWrapper = document.getElementById("wallpaperPreviewWrapper");
            const previewImg = document.getElementById("wallpaperPreviewImg");
            const previewPlaceholder = document.getElementById("wallpaperPreviewPlaceholder");

            if (!bgEnabled) {
                previewWrapper.style.display = "none";
                return;
            }

            const src = document.getElementById("bgSource").value;
            let previewUrl = "";

            if (src === "1") {
                previewUrl = '/api/wallpaper/preview?name=Background.jpeg';
            } else if (src === "2") {
                const wpName = document.getElementById("repoWpSelect").value;
                if (wpName) {
                    previewUrl = '/api/wallpaper/preview?name=' + encodeURIComponent(wpName);
                }
            } else if (src === "3") {
                if (selectedGnomeImagePath) {
                    previewUrl = '/api/wallpaper/preview?path=' + encodeURIComponent(selectedGnomeImagePath);
                }
            }

            if (previewUrl) {
                previewImg.src = previewUrl;
                previewImg.style.display = "block";
                previewPlaceholder.style.display = "none";
                previewWrapper.style.display = "block";
            } else {
                previewImg.src = "";
                previewImg.style.display = "none";
                previewPlaceholder.style.display = "flex";
                previewWrapper.style.display = "block";
            }
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
            updateWallpaperPreview();
        }

        function toggleBgInputs() {
            const src = document.getElementById("bgSource").value;
            document.getElementById("repoWallpapersWrapper").style.display = src === "2" ? "block" : "none";
            document.getElementById("customPathWrapper").style.display = src === "3" ? "block" : "none";
            updateWallpaperPreview();
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
            else if (bgSrc === "3") bgVal = "Custom Image: " + (selectedGnomeImagePath || "None");
            updateSummaryRow("sumBg", bgEnabled, bgEnabled ? ` + "`" + `Set Wallpaper (${bgVal})` + "`" + ` : "Set Background", "Skip Background");

            updateSummaryRow("sumDocker", document.getElementById("installDocker").checked, "Install Docker Rootless", "Skip Docker");
            updateSummaryRow("sumShell", document.getElementById("defaultShellZsh").checked, "Set Zsh as default", "Keep Bash");
            updateSummaryRow("sumKeyboard", document.getElementById("configureKeyboard").checked, "Configure Keyboard Layouts", "Skip Keyboard Layouts");
            updateSummaryRow("sumPower", document.getElementById("configurePower").checked, "Configure Power Settings", "Skip Power Settings");
            updateSummaryRow("sumFonts", document.getElementById("configureFonts").checked, "Install Custom Fonts", "Skip Fonts");

            let aliasesData = { customUsername: '', aliases: [] };
            document.dispatchEvent(new CustomEvent('get-aliases', {
                detail: { callback: (data) => { aliasesData = data; } }
            }));
            const customUsername = aliasesData.customUsername || "Default (None)";
            const activeAliasesCount = (aliasesData.aliases || []).filter(a => a.enabled).length;

            document.getElementById("sumPromptName").querySelector("span").textContent = customUsername;
            document.getElementById("sumAliases").querySelector("span").textContent = activeAliasesCount + " Enabled";
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
                if (applyBackground && !selectedGnomeImagePath) {
                    showToast("Please select a custom wallpaper image first.", false);
                    return;
                }
                backgroundImage = selectedGnomeImagePath;
            }

            const enableDocker = document.getElementById("installDocker").checked;
            const enableZshDefault = document.getElementById("defaultShellZsh").checked;
            const configureKeyboard = document.getElementById("configureKeyboard").checked;
            const configurePower = document.getElementById("configurePower").checked;
            const configureFonts = document.getElementById("configureFonts").checked;

            let aliasesData = { customUsername: '', aliases: [] };
            document.dispatchEvent(new CustomEvent('get-aliases', {
                detail: { callback: (data) => { aliasesData = data; } }
            }));
            const customUsername = aliasesData.customUsername;
            const aliases = aliasesData.aliases;

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
                enable_zsh_default: enableZshDefault,
                configure_keyboard: configureKeyboard,
                configure_power: configurePower,
                configure_fonts: configureFonts,
                custom_username: customUsername,
                aliases: aliases
            };

            // Switch to console view
            document.getElementById("progressBarContainer").style.display = "none";
            document.getElementById("step6").style.display = "none";
            document.getElementById("btnRow").style.display = "none";
            document.getElementById("pulseContainer").style.display = "flex";
            document.getElementById("terminalWindow").style.display = "block";

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
                    document.getElementById("terminalWindow").style.borderColor = "var(--success-color)";
                    
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
                    if (applyBackground && !selectedGnomeImagePath) {
                        showToast("Please select a custom wallpaper image first.", false);
                        return;
                    }
                    backgroundImage = selectedGnomeImagePath;
                }

                const enableDocker = document.getElementById("installDocker").checked;
                const enableZshDefault = document.getElementById("defaultShellZsh").checked;
                const configureKeyboard = document.getElementById("configureKeyboard").checked;
                const configurePower = document.getElementById("configurePower").checked;
                const configureFonts = document.getElementById("configureFonts").checked;

                let aliasesData = { customUsername: '', aliases: [] };
                document.dispatchEvent(new CustomEvent('get-aliases', {
                    detail: { callback: (data) => { aliasesData = data; } }
                }));
                const customUsername = aliasesData.customUsername;
                const aliases = aliasesData.aliases;

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
                    enable_zsh_default: enableZshDefault,
                    configure_keyboard: configureKeyboard,
                    configure_power: configurePower,
                    configure_fonts: configureFonts,
                    custom_username: customUsername,
                    aliases: aliases
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
