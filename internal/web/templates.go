package web

// IndexTemplate is the embedded HTML/CSS/JS template for our local web wizard.
const IndexTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>ZoneRestore — Configuration Wizard</title>
    <meta name="description" content="ZoneRestore web configuration wizard for automated GNOME workstation setup.">
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;500;600;700;800&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
    <script src="/js/alpine.min.js" defer></script>
    <style>
        :root {
            --bg:           #07090F;
            --sidebar-bg:   #090C13;
            --card-bg:      rgba(13, 18, 30, 0.85);
            --border:       rgba(255,255,255,0.07);
            --accent:       #00D4FF;
            --accent-dim:   rgba(0,212,255,0.12);
            --accent-glow:  rgba(0,212,255,0.3);
            --purple:       #9B59B6;
            --green:        #2ECC71;
            --red:          #E74C3C;
            --yellow:       #F39C12;
            --text:         #EEF2F7;
            --muted:        #7B8FA6;
            --r-card:       14px;
            --r-btn:        9px;
            --transition:   0.22s cubic-bezier(.4,0,.2,1);
        }

        *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

        html, body {
            height: 100%;
            font-family: 'Outfit', sans-serif;
            background: var(--bg);
            color: var(--text);
            overflow: hidden;
        }

        /* ─── Ambient glows ─── */
        .glow {
            position: fixed; border-radius: 50%; filter: blur(130px);
            opacity: .13; pointer-events: none; z-index: 0;
        }
        .glow-tl { width:560px; height:560px;
            background: radial-gradient(circle, var(--accent), transparent 70%);
            top:-140px; left:-140px; }
        .glow-br { width:560px; height:560px;
            background: radial-gradient(circle, var(--purple), transparent 70%);
            bottom:-140px; right:-140px; }

        /* ─── Root layout ─── */
        .root {
            display: flex;
            height: 100vh;
            position: relative;
            z-index: 1;
        }

        /* ─── Sidebar ─── */
        .sidebar {
            width: 260px;
            flex-shrink: 0;
            background: var(--sidebar-bg);
            border-right: 1px solid var(--border);
            display: flex;
            flex-direction: column;
            padding: 32px 20px;
            gap: 32px;
            overflow-y: auto;
        }

        .brand {
            display: flex; align-items: center; gap: 11px;
        }
        .brand-icon {
            width: 34px; height: 34px;
            background: linear-gradient(135deg, var(--accent), var(--purple));
            border-radius: 8px;
            display: flex; align-items: center; justify-content: center;
            font-weight: 800; font-size: 1rem; color: var(--bg);
            flex-shrink: 0;
            box-shadow: 0 4px 14px var(--accent-glow);
        }
        .brand-name {
            font-size: 1.1rem; font-weight: 700;
            background: linear-gradient(90deg, #fff 40%, var(--accent));
            -webkit-background-clip: text; -webkit-text-fill-color: transparent;
        }

        .nav-steps { list-style: none; display: flex; flex-direction: column; gap: 4px; }

        .nav-step {
            display: flex; align-items: center; gap: 10px;
            padding: 10px 12px; border-radius: var(--r-btn);
            cursor: pointer;
            font-size: .9rem; font-weight: 500;
            color: var(--muted);
            border: 1px solid transparent;
            transition: all var(--transition);
            user-select: none;
        }
        .nav-step:hover { background: rgba(255,255,255,.03); color: var(--text); }
        .nav-step.is-active {
            background: var(--accent-dim);
            border-color: rgba(0,212,255,.22);
            color: var(--accent);
        }
        .nav-step.is-done { color: var(--green); }
        .nav-step.is-done:hover { background: rgba(46,204,113,.06); }

        .step-num {
            width: 22px; height: 22px; border-radius: 50%;
            border: 1.5px solid currentColor;
            display: flex; align-items: center; justify-content: center;
            font-size: .72rem; font-weight: 700; flex-shrink: 0;
            transition: all var(--transition);
        }
        .nav-step.is-done .step-num {
            background: var(--green); border-color: var(--green);
            color: #07090F; font-size: .75rem;
        }
        .nav-step.is-active .step-num {
            background: var(--accent); border-color: var(--accent);
            color: var(--bg);
            box-shadow: 0 0 10px var(--accent-glow);
        }

        /* ─── Main content ─── */
        .main {
            flex: 1;
            display: flex;
            flex-direction: column;
            overflow: hidden;
        }

        .content-area {
            flex: 1;
            padding: 32px 40px;
            overflow-y: auto;
            display: flex;
            flex-direction: column;
        }

        /* ─── Glass card ─── */
        .card {
            background: var(--card-bg);
            border: 1px solid var(--border);
            border-radius: var(--r-card);
            padding: 36px 40px;
            backdrop-filter: blur(20px);
            -webkit-backdrop-filter: blur(20px);
            box-shadow: 0 20px 60px rgba(0,0,0,.5);
            flex: 1;
            display: flex;
            flex-direction: column;
        }

        .card-body {
            flex: 1;
            display: flex;
            flex-direction: column;
        }

        /* ─── Step transitions ─── */
        .step-panel {
            display: none;
            flex-direction: column;
            gap: 20px;
            flex: 1;
        }
        .step-panel.active {
            display: flex;
            animation: fadeSlideIn var(--transition) both;
        }

        @keyframes fadeSlideIn {
            from { opacity: 0; transform: translateY(12px); }
            to   { opacity: 1; transform: translateY(0); }
        }

        /* ─── Typography ─── */
        .step-title {
            font-size: 1.55rem; font-weight: 700;
            letter-spacing: -.4px;
            border-left: 3px solid var(--accent);
            padding-left: 12px;
            line-height: 1.2;
        }
        .step-sub {
            color: var(--muted); font-size: .92rem; line-height: 1.55;
            margin-top: -8px;
        }

        /* ─── Grid layout helpers ─── */
        .grid-2 { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }
        .grid-3 { display: grid; grid-template-columns: repeat(3, 1fr); gap: 12px; }
        .grid-auto { display: grid; grid-template-columns: repeat(auto-fill, minmax(170px,1fr)); gap: 14px; }

        /* ─── Choose cards (Yes/No) ─── */
        .choice-card {
            background: rgba(255,255,255,.02);
            border: 1.5px solid var(--border);
            border-radius: var(--r-card);
            padding: 22px;
            cursor: pointer;
            transition: all var(--transition);
            position: relative;
            overflow: hidden;
        }
        .choice-card:hover {
            background: rgba(255,255,255,.04);
            border-color: rgba(0,212,255,.2);
            transform: translateY(-2px);
        }
        .choice-card.chosen {
            background: rgba(0,212,255,.06);
            border-color: var(--accent);
            box-shadow: inset 0 0 0 1px rgba(0,212,255,.15), 0 0 20px rgba(0,212,255,.08);
        }
        .choice-card h3 { font-size: 1rem; font-weight: 600; margin-bottom: 6px; color: #fff; }
        .choice-card p  { font-size: .84rem; color: var(--muted); line-height: 1.4; }
        .chosen-badge {
            position: absolute; top: 12px; right: 12px;
            width: 20px; height: 20px; border-radius: 50%;
            background: var(--accent); color: var(--bg);
            font-size: .7rem; font-weight: 800;
            display: flex; align-items: center; justify-content: center;
        }

        /* ─── Toggle row (single enable/disable) ─── */
        .toggle-row {
            display: flex; align-items: center; justify-content: space-between;
            padding: 14px 18px;
            background: rgba(255,255,255,.02);
            border: 1px solid var(--border);
            border-radius: var(--r-btn);
            cursor: pointer;
            transition: all var(--transition);
        }
        .toggle-row:hover { background: rgba(255,255,255,.04); border-color: rgba(0,212,255,.2); }
        .toggle-row h4 { font-size: .95rem; font-weight: 600; }
        .toggle-badge {
            font-size: .78rem; font-weight: 700; padding: 3px 10px;
            border-radius: 20px; transition: all var(--transition);
        }
        .toggle-badge.on  { background: rgba(0,212,255,.15); color: var(--accent); }
        .toggle-badge.off { background: rgba(255,255,255,.07); color: var(--muted); }

        /* ─── Mode switcher (Dark/Light tabs) ─── */
        .mode-switch {
            display: inline-flex;
            background: rgba(0,0,0,.35);
            border: 1px solid var(--border);
            border-radius: var(--r-btn);
            padding: 4px;
            gap: 2px;
        }
        .mode-opt {
            padding: 7px 18px; border-radius: 6px;
            font-size: .87rem; font-weight: 600; cursor: pointer;
            color: var(--muted); transition: all var(--transition);
        }
        .mode-opt.active {
            background: var(--accent); color: var(--bg);
            box-shadow: 0 2px 10px var(--accent-glow);
        }

        /* ─── Scrollable list (themes) ─── */
        .scroll-list {
            max-height: 200px; overflow-y: auto;
            border: 1px solid var(--border);
            border-radius: var(--r-btn);
            background: rgba(0,0,0,.25);
            padding: 6px;
            list-style: none;
        }
        .scroll-list-item {
            padding: 9px 14px; border-radius: 6px;
            font-size: .91rem; cursor: pointer;
            display: flex; align-items: center; justify-content: space-between;
            transition: background var(--transition);
        }
        .scroll-list-item:hover { background: rgba(255,255,255,.04); }
        .scroll-list-item.chosen { background: rgba(0,212,255,.08); color: var(--accent); font-weight: 600; }

        /* ─── Font cards ─── */
        .font-card {
            background: rgba(255,255,255,.02);
            border: 1.5px solid var(--border);
            border-radius: var(--r-card);
            padding: 16px 18px;
            cursor: pointer;
            transition: all var(--transition);
        }
        .font-card:hover { background: rgba(255,255,255,.04); border-color: rgba(0,212,255,.2); transform: translateY(-1px); }
        .font-card.chosen { background: rgba(0,212,255,.05); border-color: var(--accent); }
        .font-card-name { font-size: .9rem; font-weight: 700; color: #fff; margin-bottom: 6px; }
        .font-card-sample {
            font-size: .8rem; color: var(--muted);
            white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
            padding-top: 6px; border-top: 1px dashed rgba(255,255,255,.06);
            font-family: monospace;
        }
        .font-empty {
            grid-column: 1/-1;
            text-align: center; color: var(--muted);
            font-size: .9rem; padding: 24px;
        }

        /* ─── Wallpaper gallery ─── */
        .wp-gallery {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
            gap: 12px;
            max-height: 300px;
            overflow-y: auto;
            padding: 4px 2px;
        }
        .wp-card {
            border-radius: var(--r-card);
            overflow: hidden;
            aspect-ratio: 16/10;
            cursor: pointer;
            border: 2px solid transparent;
            position: relative;
            background: rgba(0,0,0,.6);
            transition: all var(--transition);
        }
        .wp-card img {
            width: 100%; height: 100%; object-fit: cover;
            opacity: .72; transition: opacity var(--transition), transform .35s ease;
            display: block;
        }
        .wp-card:hover img { opacity: .95; transform: scale(1.04); }
        .wp-card.chosen { border-color: var(--accent); box-shadow: 0 0 0 1px var(--accent), 0 0 18px var(--accent-glow); }
        .wp-card.chosen img { opacity: 1; }
        .wp-overlay {
            position: absolute; bottom: 0; left: 0; right: 0;
            background: linear-gradient(transparent, rgba(0,0,0,.8));
            padding: 6px 10px 8px;
            font-size: .72rem; color: #fff; font-weight: 500;
            white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
        }
        .wp-check {
            position: absolute; top: 8px; right: 8px;
            width: 20px; height: 20px; border-radius: 50%;
            background: var(--accent); color: var(--bg);
            font-size: .65rem; font-weight: 900;
            display: flex; align-items: center; justify-content: center;
            box-shadow: 0 2px 8px var(--accent-glow);
        }
        .wp-loading {
            grid-column: 1/-1; text-align: center;
            color: var(--muted); font-size: .9rem; padding: 30px;
            display: flex; align-items: center; justify-content: center; gap: 10px;
        }

        /* ─── Form inputs ─── */
        .form-label {
            font-size: .8rem; font-weight: 700; color: var(--muted);
            text-transform: uppercase; letter-spacing: .5px;
            margin-bottom: 7px; display: block;
        }
        .text-input {
            width: 100%;
            background: rgba(0,0,0,.35);
            border: 1px solid var(--border);
            border-radius: var(--r-btn);
            padding: 12px 16px;
            color: #fff;
            font-family: inherit; font-size: .93rem;
            transition: all var(--transition);
        }
        .text-input:focus {
            outline: none;
            border-color: var(--accent);
            background: rgba(0,0,0,.5);
            box-shadow: 0 0 0 3px rgba(0,212,255,.12);
        }
        .input-row { display: flex; gap: 10px; }

        /* ─── Summary list ─── */
        .summary-list {
            display: flex; flex-direction: column; gap: 8px;
            max-height: 340px; overflow-y: auto;
            background: rgba(0,0,0,.25);
            border: 1px solid var(--border);
            border-radius: var(--r-card);
            padding: 14px;
        }
        .s-row {
            display: flex; align-items: center; justify-content: space-between;
            padding: 8px 10px; border-radius: 7px;
            font-size: .9rem;
            background: rgba(255,255,255,.015);
        }
        .s-label { display: flex; align-items: center; gap: 10px; }
        .s-badge {
            font-size: .75rem; font-weight: 700;
            padding: 2px 9px; border-radius: 12px;
        }
        .s-badge.on  { background: rgba(46,204,113,.15); color: var(--green); }
        .s-badge.off { background: rgba(231,76,60,.12);  color: var(--red); }
        .s-sub {
            font-size: .82rem; color: var(--muted); padding: 4px 10px 4px 20px;
            overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
        }

        /* ─── Console ─── */
        .console {
            background: #040608; border: 1px solid var(--border);
            border-radius: var(--r-card);
            padding: 18px 20px;
            font-family: 'JetBrains Mono', monospace;
            font-size: .82rem;
            height: 260px; overflow-y: auto;
            display: flex; flex-direction: column; gap: 4px;
            box-shadow: inset 0 4px 24px rgba(0,0,0,.7);
        }
        .c-line { line-height: 1.45; color: #9EB3CC; }
        .c-info    { color: #7AA2F7; }
        .c-success { color: var(--green); }
        .c-error   { color: var(--red); }
        .c-warn    { color: var(--yellow); }

        .spinner-row {
            display: flex; align-items: center; gap: 10px;
            color: var(--accent); font-weight: 600; font-size: .9rem;
            margin-top: 10px;
        }
        .pulse {
            width: 9px; height: 9px; border-radius: 50%;
            background: var(--accent);
            animation: pulse-anim 1.2s ease-in-out infinite;
        }
        @keyframes pulse-anim {
            0%,100% { transform: scale(.8); opacity: .5; box-shadow: none; }
            50%      { transform: scale(1.15); opacity: 1; box-shadow: 0 0 10px var(--accent-glow); }
        }

        /* ─── Buttons ─── */
        .btn { 
            padding: 11px 22px; border-radius: var(--r-btn);
            font-family: inherit; font-size: .92rem; font-weight: 600;
            cursor: pointer; border: 1px solid transparent;
            transition: all var(--transition);
        }
        .btn-primary {
            background: var(--accent); color: var(--bg);
            box-shadow: 0 4px 14px rgba(0,212,255,.2);
        }
        .btn-primary:hover { background: #1DDDFF; box-shadow: 0 6px 20px rgba(0,212,255,.35); }
        .btn-secondary {
            background: rgba(255,255,255,.05);
            color: var(--text); border-color: var(--border);
        }
        .btn-secondary:hover { background: rgba(255,255,255,.09); border-color: rgba(255,255,255,.15); }
        .btn-danger {
            background: rgba(231,76,60,.12); color: var(--red); border-color: rgba(231,76,60,.2);
        }
        .btn-danger:hover { background: rgba(231,76,60,.2); }
        .btn:disabled { opacity: .45; cursor: not-allowed; }

        /* ─── Nav actions row ─── */
        .nav-actions {
            display: flex; justify-content: space-between; align-items: center;
            padding: 20px 0 0;
            border-top: 1px solid var(--border);
            margin-top: auto;
        }

        /* ─── Bottom status bar ─── */
        .status-bar {
            height: 50px; flex-shrink: 0;
            background: var(--sidebar-bg); border-top: 1px solid var(--border);
            display: flex; align-items: center;
            padding: 0 28px; gap: 10px; justify-content: space-between;
            font-size: .8rem;
        }
        .pills { display: flex; gap: 8px; flex-wrap: wrap; }
        .pill {
            display: flex; align-items: center; gap: 6px;
            padding: 3px 11px; border-radius: 20px;
            border: 1px solid var(--border);
            color: var(--muted); font-size: .78rem; font-weight: 500;
            transition: all var(--transition);
        }
        .pill.on { border-color: rgba(0,212,255,.25); color: var(--accent); background: rgba(0,212,255,.04); }
        .pill-dot { width: 5px; height: 5px; border-radius: 50%; background: currentColor; }

        /* ─── Toast ─── */
        .toast {
            position: fixed; bottom: 70px; right: 24px;
            background: rgba(13,18,30,.97); border: 1px solid var(--accent);
            border-radius: var(--r-btn); padding: 12px 22px;
            color: #fff; font-weight: 600; font-size: .9rem;
            z-index: 9999; display: flex; align-items: center; gap: 10px;
            box-shadow: 0 8px 30px rgba(0,212,255,.25);
            transform: translateY(30px); opacity: 0;
            transition: all .3s cubic-bezier(.16,1,.3,1);
            pointer-events: none;
        }
        .toast.show { transform: translateY(0); opacity: 1; }

        /* ─── Scrollbars ─── */
        ::-webkit-scrollbar { width: 5px; }
        ::-webkit-scrollbar-track { background: transparent; }
        ::-webkit-scrollbar-thumb { background: rgba(255,255,255,.09); border-radius: 4px; }
        ::-webkit-scrollbar-thumb:hover { background: rgba(255,255,255,.18); }

        /* ─── Divider ─── */
        .divider { height: 1px; background: var(--border); margin: 4px 0; }

        /* ─── Section label ─── */
        .sec-label {
            font-size: .75rem; font-weight: 700;
            color: var(--muted); text-transform: uppercase; letter-spacing: .6px;
            margin-bottom: 8px;
        }

        /* ─── img error fallback ─── */
        img.broken { opacity: 0; }
    </style>
</head>
<body x-data="wizard()" x-init="init()">

    <div class="glow glow-tl"></div>
    <div class="glow glow-br"></div>

    <div class="root">

        <!-- ═══════ SIDEBAR ═══════ -->
        <aside class="sidebar">
            <div class="brand">
                <div class="brand-icon">ZR</div>
                <div class="brand-name">ZoneRestore</div>
            </div>

            <ul class="nav-steps">
                <template x-for="(s, i) in steps" :key="i">
                    <li class="nav-step"
                        :class="{'is-active': step===i, 'is-done': step>i}"
                        @click="goTo(i)">
                        <div class="step-num">
                            <span x-show="step <= i" x-text="i+1"></span>
                            <span x-show="step > i">✓</span>
                        </div>
                        <span x-text="s"></span>
                    </li>
                </template>
            </ul>
        </aside>

        <!-- ═══════ MAIN ═══════ -->
        <div class="main">
            <div class="content-area">
                <div class="card">
                    <div class="card-body">

                        <!-- ── STEP 0 – Import ── -->
                        <div id="step-0" class="step-panel" :class="{'active': step===0}">
                            <div>
                                <h1 class="step-title">Configuration Profile</h1>
                                <p class="step-sub">Start from your saved settings or configure everything from scratch.</p>
                            </div>
                            <div class="grid-2">
                                <div class="choice-card" :class="{'chosen': importSettings}" @click="importSettings=true">
                                    <div class="chosen-badge" x-show="importSettings">✓</div>
                                    <h3>Import saved profile</h3>
                                    <p>Load previously configured preferences from config.json instantly.</p>
                                </div>
                                <div class="choice-card" :class="{'chosen': !importSettings}" @click="importSettings=false">
                                    <div class="chosen-badge" x-show="!importSettings">✓</div>
                                    <h3>Fresh defaults</h3>
                                    <p>Reset and configure every setting from recommended defaults.</p>
                                </div>
                            </div>
                            <div style="display:flex; gap:10px; flex-wrap:wrap; align-items:center;">
                                <button class="btn btn-secondary" @click="loadDefault()">↺ Reset Defaults</button>
                                <button class="btn btn-secondary" @click="loadFromDisk()">💾 Load Saved</button>
                                <!-- Hidden file input for JSON upload -->
                                <label class="btn btn-secondary" style="cursor:pointer;" title="Upload a config JSON file from your computer">
                                    📂 Upload JSON
                                    <input type="file" id="config-file-input" accept=".json,application/json"
                                           style="display:none;"
                                           @change="uploadConfig($event)">
                                </label>
                            </div>
                        </div>

                        <!-- ── STEP 1 – Oh-My-Zsh ── -->
                        <div id="step-1" class="step-panel" :class="{'active': step===1}">
                            <div>
                                <h1 class="step-title">Zsh Environment</h1>
                                <p class="step-sub">Install and configure Oh-My-Zsh shell framework unattended on your workspace.</p>
                            </div>
                            <div class="grid-2">
                                <div class="choice-card" :class="{'chosen': cfg.zsh.install_oh_my_zsh}" @click="cfg.zsh.install_oh_my_zsh=true">
                                    <div class="chosen-badge" x-show="cfg.zsh.install_oh_my_zsh">✓</div>
                                    <h3>Install Oh-My-Zsh</h3>
                                    <p>Runs the official unattended install script and configures plugins.</p>
                                </div>
                                <div class="choice-card" :class="{'chosen': !cfg.zsh.install_oh_my_zsh}" @click="cfg.zsh.install_oh_my_zsh=false">
                                    <div class="chosen-badge" x-show="!cfg.zsh.install_oh_my_zsh">✓</div>
                                    <h3>Skip shell setup</h3>
                                    <p>Leave your existing shell configuration untouched.</p>
                                </div>
                            </div>
                        </div>

                        <!-- ── STEP 2 – Git ── -->
                        <div id="step-2" class="step-panel" :class="{'active': step===2}">
                            <div>
                                <h1 class="step-title">Git Global Credentials</h1>
                                <p class="step-sub">Configure your identity for all repository interactions on this machine.</p>
                            </div>
                            <div class="grid-2" style="margin-bottom:4px;">
                                <div class="choice-card" :class="{'chosen': cfg.git.configure_git}" @click="cfg.git.configure_git=true">
                                    <div class="chosen-badge" x-show="cfg.git.configure_git">✓</div>
                                    <h3>Configure Git globally</h3>
                                    <p>Apply name and email to <code>~/.gitconfig</code>.</p>
                                </div>
                                <div class="choice-card" :class="{'chosen': !cfg.git.configure_git}" @click="cfg.git.configure_git=false">
                                    <div class="chosen-badge" x-show="!cfg.git.configure_git">✓</div>
                                    <h3>Skip Git setup</h3>
                                    <p>Keep existing global git identity unchanged.</p>
                                </div>
                            </div>
                            <div x-show="cfg.git.configure_git" x-transition style="display:flex;flex-direction:column;gap:14px;">
                                <div class="grid-2">
                                    <div>
                                        <label class="form-label" for="git-name">Full name</label>
                                        <input id="git-name" class="text-input" type="text" x-model="cfg.git.git_name" placeholder="e.g. John Doe">
                                    </div>
                                    <div>
                                        <label class="form-label" for="git-email">Email address</label>
                                        <input id="git-email" class="text-input" type="email" x-model="cfg.git.git_email" placeholder="e.g. john@example.com">
                                    </div>
                                </div>
                            </div>
                        </div>

                        <!-- ── STEP 3 – Theme ── -->
                        <div id="step-3" class="step-panel" :class="{'active': step===3}">
                            <div>
                                <h1 class="step-title">Gnome Desktop Theme</h1>
                                <p class="step-sub">Customize the GTK window decorations and interface color scheme.</p>
                            </div>
                            <div class="toggle-row" @click="cfg.theme.apply_theme=!cfg.theme.apply_theme">
                                <h4>Apply Gnome theme settings</h4>
                                <span class="toggle-badge" :class="cfg.theme.apply_theme?'on':'off'" x-text="cfg.theme.apply_theme?'Enabled':'Disabled'"></span>
                            </div>
                            <div x-show="cfg.theme.apply_theme" x-transition style="display:flex;flex-direction:column;gap:14px;">
                                <div>
                                    <p class="sec-label">Color mode</p>
                                    <div class="mode-switch">
                                        <div class="mode-opt" :class="{'active': cfg.theme.theme_mode==='1'}" @click="switchThemeMode('1')">🌙 Dark</div>
                                        <div class="mode-opt" :class="{'active': cfg.theme.theme_mode==='2'}" @click="switchThemeMode('2')">☀ Light</div>
                                    </div>
                                </div>
                                <div>
                                    <p class="sec-label">Available system themes
                                        <span x-show="themeList.length === 0" style="color:var(--red)"> — none found on system</span>
                                    </p>
                                    <ul class="scroll-list">
                                        <template x-for="t in themeList" :key="t">
                                            <li class="scroll-list-item" :class="{'chosen': cfg.theme.theme_name===t}" @click="cfg.theme.theme_name=t">
                                                <span x-text="t"></span>
                                                <span x-show="cfg.theme.theme_name===t">✓</span>
                                            </li>
                                        </template>
                                        <li x-show="themeList.length===0" class="scroll-list-item" style="color:var(--muted);cursor:default;">No themes detected</li>
                                    </ul>
                                </div>
                            </div>
                        </div>

                        <!-- ── STEP 4 – Fonts ── -->
                        <div id="step-4" class="step-panel" :class="{'active': step===4}">
                            <div>
                                <h1 class="step-title">Fonts</h1>
                                <p class="step-sub">Install fonts from the themes repository and choose a display font for your GNOME desktop.</p>
                            </div>

                            <!-- Locked terminal font notice -->
                            <div style="display:flex;align-items:center;gap:12px;padding:12px 16px;background:rgba(0,212,255,0.07);border:1px solid rgba(0,212,255,0.2);border-radius:10px;">
                                <span style="font-size:1.3rem;">🔒</span>
                                <div>
                                    <div style="font-size:.85rem;color:var(--muted);margin-bottom:2px;">Terminal font — always applied</div>
                                    <div style="font-family:'MesloLGS NF Regular','JetBrains Mono',monospace;font-size:1rem;color:var(--accent);font-weight:600;">MesloLGS NF Regular 12</div>
                                    <div style="font-size:.75rem;color:var(--muted);margin-top:2px;">Applied to all GNOME Terminal profiles and system monospace font.</div>
                                </div>
                            </div>

                            <div class="toggle-row" @click="cfg.fonts.configure_fonts=!cfg.fonts.configure_fonts">
                                <h4>Install &amp; configure fonts</h4>
                                <span class="toggle-badge" :class="cfg.fonts.configure_fonts?'on':'off'" x-text="cfg.fonts.configure_fonts?'Enabled':'Disabled'"></span>
                            </div>
                            <div x-show="cfg.fonts.configure_fonts" x-transition style="display:flex;flex-direction:column;gap:12px;">
                                <p class="sec-label">Choose a display font for GNOME interface
                                    <span x-show="fonts.length === 0 && !fontsLoading" style="color:var(--yellow)"> — none found in themes/fonts/</span>
                                </p>
                                <div x-show="fontsLoading" style="color:var(--muted);font-size:.9rem;padding:16px 0;">Loading fonts from repository…</div>
                                <div class="grid-auto" x-show="!fontsLoading">
                                    <!-- "System Default" option -->
                                    <div class="font-card" :class="{'chosen': cfg.fonts.display_font_name===''}" @click="cfg.fonts.display_font_name=''">
                                        <div class="font-card-name">System Default</div>
                                        <div class="font-card-sample" style="font-family:sans-serif;">Ubuntu — System UI font</div>
                                    </div>
                                    <template x-if="fonts.length === 0">
                                        <div class="font-empty">No additional fonts found in <code>themes/fonts/</code>.</div>
                                    </template>
                                    <template x-for="f in fonts" :key="f.name">
                                        <div class="font-card" :class="{'chosen': cfg.fonts.display_font_name===f.name}" @click="cfg.fonts.display_font_name=f.name">
                                            <div class="font-card-name" x-text="f.name"></div>
                                            <div class="font-card-sample"
                                                 :style="'font-family:' + JSON.stringify(f.name) + ',sans-serif'"
                                                 x-text="f.name + ' — The quick brown fox'"></div>
                                        </div>
                                    </template>
                                </div>
                            </div>
                        </div>


                        <!-- ── STEP 5 – Wallpaper ── -->
                        <div id="step-5" class="step-panel" :class="{'active': step===5}">
                            <div>
                                <h1 class="step-title">Desktop Wallpaper</h1>
                                <p class="step-sub">Choose a background image loaded live from the ZoneRestoreThemes repository.</p>
                            </div>
                            <div class="toggle-row" @click="cfg.wallpaper.apply_background=!cfg.wallpaper.apply_background">
                                <h4>Apply desktop wallpaper</h4>
                                <span class="toggle-badge" :class="cfg.wallpaper.apply_background?'on':'off'" x-text="cfg.wallpaper.apply_background?'Enabled':'Disabled'"></span>
                            </div>
                            <div x-show="cfg.wallpaper.apply_background" x-transition style="display:flex;flex-direction:column;gap:12px;">
                                <div class="mode-switch">
                                    <div class="mode-opt" :class="{'active': wpSource==='repo'}" @click="wpSource='repo'">🖼 Repository gallery</div>
                                    <div class="mode-opt" :class="{'active': wpSource==='custom'}" @click="wpSource='custom'">📁 Custom path</div>
                                </div>

                                <!-- Gallery view -->
                                <div x-show="wpSource==='repo'" x-transition>
                                    <p class="sec-label" x-text="'Wallpapers (' + wallpapers.length + ' found)'"></p>
                                    <div x-show="wallpapers.length===0 && !wpLoading" style="color:var(--muted);font-size:.88rem;padding:8px 0;">No wallpapers found in <code>themes/wallpapers/</code>.</div>
                                    <div x-show="wpLoading" class="wp-loading"><div class="pulse"></div> Loading wallpapers…</div>
                                    <div class="wp-gallery" x-show="!wpLoading">
                                        <template x-for="wp in wallpapers" :key="wp">
                                            <div class="wp-card" :class="{'chosen': cfg.wallpaper.background_image===wp}" @click="cfg.wallpaper.background_image=wp">
                                                <img :src="'/api/wallpaper/preview?name=' + encodeURIComponent(wp)"
                                                     :alt="wp"
                                                     loading="lazy"
                                                     @error="$el.style.display='none'">
                                                <div class="wp-overlay" x-text="wp"></div>
                                                <div class="wp-check" x-show="cfg.wallpaper.background_image===wp">✓</div>
                                            </div>
                                        </template>
                                    </div>
                                </div>

                                <!-- Custom path view -->
                                <div x-show="wpSource==='custom'" x-transition>
                                    <label class="form-label" for="custom-wp">Absolute image file path</label>
                                    <div class="input-row">
                                        <input id="custom-wp" class="text-input" type="text" x-model="cfg.wallpaper.background_image"
                                               placeholder="/home/user/Pictures/wallpaper.jpg">
                                        <button class="btn btn-secondary" @click="browseWallpaper()">Browse…</button>
                                    </div>
                                </div>

                                <!-- ── Persistent preview: shown for any selected wallpaper ── -->
                                <div x-show="cfg.wallpaper.background_image" x-transition style="margin-top:4px;">
                                    <p class="sec-label" style="margin-bottom:6px;">Preview</p>
                                    <div style="position:relative;width:100%;max-width:480px;border-radius:10px;overflow:hidden;border:2px solid var(--border);box-shadow:0 4px 24px rgba(0,0,0,.4);">
                                        <img
                                            :src="wpSource==='repo'
                                                ? '/api/wallpaper/preview?name=' + encodeURIComponent(cfg.wallpaper.background_image)
                                                : '/api/wallpaper/preview?path=' + encodeURIComponent(cfg.wallpaper.background_image)"
                                            :alt="cfg.wallpaper.background_image"
                                            style="width:100%;height:220px;object-fit:cover;display:block;"
                                            @error="$el.style.opacity='0.2'"
                                            @load="$el.style.opacity='1'"
                                            style="transition:opacity .3s;">
                                        <div style="position:absolute;bottom:0;left:0;right:0;padding:6px 10px;background:linear-gradient(transparent,rgba(0,0,0,.7));font-size:.78rem;color:#fff;" x-text="cfg.wallpaper.background_image.split('/').pop()"></div>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <!-- ── STEP 6 – System defaults ── -->
                        <div id="step-6" class="step-panel" :class="{'active': step===6}">
                            <div>
                                <h1 class="step-title">Environments &amp; Defaults</h1>
                                <p class="step-sub">Enable rootless Docker and configure shell preferences.</p>
                            </div>
                            <div style="display:flex;flex-direction:column;gap:10px;">
                                <div class="toggle-row" @click="cfg.docker.enable_docker=!cfg.docker.enable_docker">
                                    <div style="display:flex;align-items:center;gap:10px;">
                                        <span>🐳</span>
                                        <div>
                                            <h4>Docker Rootless Mode</h4>
                                            <div style="font-size:.8rem;color:var(--muted);margin-top:2px;">Installs Docker Engine in non-root user mode via get.docker.com/rootless</div>
                                        </div>
                                    </div>
                                    <span class="toggle-badge" :class="cfg.docker.enable_docker?'on':'off'" x-text="cfg.docker.enable_docker?'On':'Off'"></span>
                                </div>
                                <div class="toggle-row" @click="cfg.shell.enable_zsh_default=!cfg.shell.enable_zsh_default">
                                    <div style="display:flex;align-items:center;gap:10px;">
                                        <span>🐚</span>
                                        <div>
                                            <h4>Default shell → Zsh</h4>
                                            <div style="font-size:.8rem;color:var(--muted);margin-top:2px;">Adds Zsh launch commands to <code>~/.bashrc</code></div>
                                        </div>
                                    </div>
                                    <span class="toggle-badge" :class="cfg.shell.enable_zsh_default?'on':'off'" x-text="cfg.shell.enable_zsh_default?'On':'Off'"></span>
                                </div>
                                <div class="toggle-row" @click="cfg.dock.pin_discord=!cfg.dock.pin_discord">
                                    <div style="display:flex;align-items:center;gap:10px;">
                                        <span>⚓</span>
                                        <div>
                                            <h4>Pin Discord to favorites in Dock</h4>
                                            <div style="font-size:.8rem;color:var(--muted);margin-top:2px;">Pins Discord launcher into GNOME favorites apps dock panel</div>
                                        </div>
                                    </div>
                                    <span class="toggle-badge" :class="cfg.dock.pin_discord?'on':'off'" x-text="cfg.dock.pin_discord?'On':'Off'"></span>
                                </div>
                                <div class="toggle-row" @click="cfg.keyboard.configure_keyboard=!cfg.keyboard.configure_keyboard">
                                    <div style="display:flex;align-items:center;gap:10px;">
                                        <span>⌨️</span>
                                        <div>
                                            <h4>US + FR keyboard layouts</h4>
                                            <div style="font-size:.8rem;color:var(--muted);margin-top:2px;">Applies layout configuration to system input sources</div>
                                        </div>
                                    </div>
                                    <span class="toggle-badge" :class="cfg.keyboard.configure_keyboard?'on':'off'" x-text="cfg.keyboard.configure_keyboard?'On':'Off'"></span>
                                </div>
                                <div class="toggle-row" x-show="cfg.keyboard.configure_keyboard" x-transition @click="cfg.keyboard.add_arabic=!cfg.keyboard.add_arabic">
                                    <div style="display:flex;align-items:center;gap:10px;">
                                        <span>🌍</span>
                                        <div>
                                            <h4>Add Arabic layout to keyboard list</h4>
                                            <div style="font-size:.8rem;color:var(--muted);margin-top:2px;">Applies dual English/French layout plus standard Arabic layout</div>
                                        </div>
                                    </div>
                                    <span class="toggle-badge" :class="cfg.keyboard.add_arabic?'on':'off'" x-text="cfg.keyboard.add_arabic?'On':'Off'"></span>
                                </div>
                                <div class="toggle-row" @click="cfg.power.configure_power=!cfg.power.configure_power">
                                    <div style="display:flex;align-items:center;gap:10px;">
                                        <span>🔋</span>
                                        <div>
                                            <h4>Power management</h4>
                                            <div style="font-size:.8rem;color:var(--muted);margin-top:2px;">Sets sleep timeout to 1.5 hours via gsettings</div>
                                        </div>
                                    </div>
                                    <span class="toggle-badge" :class="cfg.power.configure_power?'on':'off'" x-text="cfg.power.configure_power?'On':'Off'"></span>
                                </div>
                            </div>
                        </div>

                        <!-- ── STEP 7 – Summary ── -->
                        <div id="step-7" class="step-panel" :class="{'active': step===7}">
                            <div>
                                <h1 class="step-title">Review &amp; Execute</h1>
                                <p class="step-sub">Confirm your selections below, then click <strong>Run Setup</strong> to apply everything.</p>
                            </div>
                            <div class="summary-list">
                                <div class="s-row">
                                    <div class="s-label">🛠 Install Oh-My-Zsh</div>
                                    <span class="s-badge" :class="cfg.zsh.install_oh_my_zsh?'on':'off'" x-text="cfg.zsh.install_oh_my_zsh?'Yes':'No'"></span>
                                </div>
                                <div class="s-row">
                                    <div class="s-label">👤 Configure Git globally</div>
                                    <span class="s-badge" :class="cfg.git.configure_git?'on':'off'" x-text="cfg.git.configure_git?'Yes':'No'"></span>
                                </div>
                                <div class="s-sub" x-show="cfg.git.configure_git" x-text="cfg.git.git_name + ' &lt;' + cfg.git.git_email + '&gt;'"></div>
                                <div class="s-row">
                                    <div class="s-label">🎨 Apply Gnome theme</div>
                                    <span class="s-badge" :class="cfg.theme.apply_theme?'on':'off'" x-text="cfg.theme.apply_theme?'Yes':'No'"></span>
                                </div>
                                <div class="s-sub" x-show="cfg.theme.apply_theme" x-text="(cfg.theme.theme_mode==='1'?'Dark':'Light') + ' — ' + cfg.theme.theme_name"></div>
                                <div class="s-row">
                                    <div class="s-label">🔤 Install repo fonts</div>
                                    <span class="s-badge" :class="cfg.fonts.configure_fonts?'on':'off'" x-text="cfg.fonts.configure_fonts?'Yes':'No'"></span>
                                </div>
                                <div class="s-sub" x-show="cfg.fonts.configure_fonts">
                                    <span style="color:var(--accent);">Terminal:</span> MesloLGS NF Regular 12 (locked)
                                    <span x-show="cfg.fonts.display_font_name" x-text="' · Display: ' + cfg.fonts.display_font_name"></span>
                                </div>
                                <div class="s-row">
                                    <div class="s-label">🖼 Apply wallpaper</div>
                                    <span class="s-badge" :class="cfg.wallpaper.apply_background?'on':'off'" x-text="cfg.wallpaper.apply_background?'Yes':'No'"></span>
                                </div>
                                <div class="s-sub" x-show="cfg.wallpaper.apply_background && cfg.wallpaper.background_image" x-text="cfg.wallpaper.background_image"></div>
                                <div class="s-row">
                                    <div class="s-label">🐳 Docker rootless</div>
                                    <span class="s-badge" :class="cfg.docker.enable_docker?'on':'off'" x-text="cfg.docker.enable_docker?'Yes':'No'"></span>
                                </div>
                                <div class="s-row">
                                    <div class="s-label">🐚 Default shell → Zsh</div>
                                    <span class="s-badge" :class="cfg.shell.enable_zsh_default?'on':'off'" x-text="cfg.shell.enable_zsh_default?'Yes':'No'"></span>
                                </div>
                                <div class="s-row">
                                    <div class="s-label">⚓ Pin Discord to Dock</div>
                                    <span class="s-badge" :class="cfg.dock.pin_discord?'on':'off'" x-text="cfg.dock.pin_discord?'Yes':'No'"></span>
                                </div>
                                <div class="s-row">
                                    <div class="s-label">⌨️ US+FR keyboard</div>
                                    <span class="s-badge" :class="cfg.keyboard.configure_keyboard?'on':'off'" x-text="cfg.keyboard.configure_keyboard?'Yes':'No'"></span>
                                </div>
                                <div class="s-sub" x-show="cfg.keyboard.configure_keyboard" x-text="'Layouts: US + FR' + (cfg.keyboard.add_arabic ? ' + AR (Arabic)' : '')"></div>
                                <div class="s-row">
                                    <div class="s-label">🔋 Power management</div>
                                    <span class="s-badge" :class="cfg.power.configure_power?'on':'off'" x-text="cfg.power.configure_power?'Yes':'No'"></span>
                                </div>
                            </div>
                        </div>

                        <!-- ── STEP 8 – Execution ── -->
                        <div id="step-8" class="step-panel" :class="{'active': step===8}">
                            <div>
                                <h1 class="step-title">Applying Configuration</h1>
                                <p class="step-sub">Setup is running live. Please do not close this window.</p>
                            </div>
                            <div class="console" id="console-output">
                                <template x-for="(l, i) in logs" :key="i">
                                    <div class="c-line" :class="l.cls" x-text="l.text"></div>
                                </template>
                            </div>
                            <div class="spinner-row" x-show="!execFinished">
                                <div class="pulse"></div>
                                <span>Executing installer tasks… please wait.</span>
                            </div>
                            <div x-show="execFinished && !execErr" x-transition style="display:flex;flex-direction:column;gap:10px;margin-top:8px;">
                                <p style="color:var(--green);font-weight:700;">✔ Setup completed successfully!</p>
                                <!-- Action buttons — hidden once closing starts -->
                                <div x-show="!isClosing" style="display:flex;gap:10px;flex-wrap:wrap;">
                                    <button class="btn btn-primary" @click="finish(true)">💾 Save to disk &amp; Download</button>
                                    <button class="btn btn-secondary" @click="finish(false)">🔄 Restart Terminal</button>
                                </div>
                                <!-- Closing screen shown after finish() is triggered -->
                                <div x-show="isClosing" x-transition style="display:flex;flex-direction:column;align-items:center;gap:14px;padding:24px 0;">
                                    <div style="font-size:2.4rem;">✅</div>
                                    <p style="font-size:1.1rem;font-weight:700;color:var(--accent);">All done! Closing ZoneRestore…</p>
                                    <p style="color:var(--muted);font-size:.9rem;">Restarting your terminal. You can close this tab now.</p>
                                    <div class="pulse" style="width:10px;height:10px;"></div>
                                </div>
                            </div>
                            <div x-show="execFinished && execErr" x-transition style="margin-top:10px;">
                                <p style="color:var(--red);font-weight:700;margin-bottom:8px;">✖ An error occurred.</p>
                                <button class="btn btn-secondary" @click="goTo(7)">← Back to Summary</button>
                            </div>
                        </div>

                        <!-- ── NAV ACTIONS ── -->
                        <div class="nav-actions" x-show="step < 8">
                            <button class="btn btn-secondary"
                                    @click="prev()"
                                    :style="step===0 ? 'opacity:0;pointer-events:none;' : ''">
                                ← Back
                            </button>
                            <button class="btn btn-primary" @click="next()"
                                    x-text="step===7 ? '▶ Run Setup' : 'Next →'">
                            </button>
                        </div>

                    </div><!-- card-body -->
                </div><!-- card -->
            </div><!-- content-area -->

            <!-- ══ BOTTOM STATUS BAR ══ -->
            <div class="status-bar">
                <div class="pills">
                    <div class="pill" :class="{'on': cfg.zsh.install_oh_my_zsh}">
                        <div class="pill-dot"></div> Oh-My-Zsh
                    </div>
                    <div class="pill" :class="{'on': cfg.git.configure_git}">
                        <div class="pill-dot"></div> Git
                    </div>
                    <div class="pill" :class="{'on': cfg.theme.apply_theme}">
                        <div class="pill-dot"></div> Theme
                    </div>
                    <div class="pill" :class="{'on': cfg.fonts.configure_fonts}">
                        <div class="pill-dot"></div>
                        <span x-text="cfg.fonts.configure_fonts ? (cfg.fonts.display_font_name || 'MesloLGS NF') : 'Fonts'"></span>
                    </div>
                    <div class="pill" :class="{'on': cfg.wallpaper.apply_background}">
                        <div class="pill-dot"></div> Wallpaper
                    </div>
                    <div class="pill" :class="{'on': cfg.docker.enable_docker}">
                        <div class="pill-dot"></div> Docker
                    </div>
                    <div class="pill" :class="{'on': cfg.shell.enable_zsh_default}">
                        <div class="pill-dot"></div> Zsh
                    </div>
                    <div class="pill" :class="{'on': cfg.dock.pin_discord}">
                        <div class="pill-dot"></div> Dock Favorites
                    </div>
                </div>
                <span style="color:var(--muted);">ZoneRestore v2 · step <span x-text="step+1"></span> / 9</span>
            </div>
        </div><!-- main -->

    </div><!-- root -->

    <!-- Toast -->
    <div class="toast" id="toast"><span id="toast-msg">✔ Done!</span></div>

    <script>
    document.addEventListener('alpine:init', () => {
        Alpine.data('wizard', () => ({
            step: 0,
            steps: [
                'Import Profile',
                'Oh-My-Zsh',
                'Git Credentials',
                'Gnome Theme',
                'Fonts',
                'Wallpaper',
                'System Defaults',
                'Summary'
            ],

            importSettings: true,
            wpSource: 'repo',

            cfg: {
                zsh: {
                    install_oh_my_zsh: true
                },
                git: {
                    configure_git: true,
                    git_name: '',
                    git_email: ''
                },
                theme: {
                    apply_theme: true,
                    theme_mode: '1',
                    theme_name: ''
                },
                fonts: {
                    configure_fonts: true,
                    font_name: 'MesloLGS NF',  // locked terminal font — always MesloLGS
                    display_font_name: ''       // user-selectable display/UI font
                },
                wallpaper: {
                    apply_background: true,
                    background_image: ''
                },
                docker: {
                    enable_docker: true
                },
                dock: {
                    pin_discord: true
                },
                keyboard: {
                    configure_keyboard: true,
                    add_arabic: false
                },
                power: {
                    configure_power: true
                },
                shell: {
                    enable_zsh_default: false
                },
                custom_username: '',
                aliases: []
            },

            // Resources loaded from API
            darkThemes:  [],
            lightThemes: [],
            fonts:       [],
            wallpapers:  [],

            fontsLoading:   true,
            wpLoading:      true,

            // Computed shortcut
            get themeList() {
                return this.cfg.theme.theme_mode === '1' ? this.darkThemes : this.lightThemes;
            },

            // Execution
            logs:         [],
            execFinished: false,
            execErr:      false,
            isClosing:    false,

            // ── INIT ──
            async init() {
                // Load stored cfg
                try {
                    const r = await fetch('/api/config');
                    if (r.ok) this.applyRemoteCfg(await r.json());
                } catch(e) {}

                // Load resources (themes, fonts, wallpapers)
                this.fetchResources();
            },

            async fetchResources() {
                this.fontsLoading = true;
                this.wpLoading    = true;
                try {
                    const r = await fetch('/api/resources');
                    if (!r.ok) return;
                    const d = await r.json();
                    this.darkThemes  = d.dark_themes  || [];
                    this.lightThemes = d.light_themes || [];

                    // Fonts: array of { name, files[] }
                    this.fonts = d.fonts || [];
                    this.injectFontFaces(this.fonts);

                    // Wallpapers: array of strings
                    this.wallpapers = d.wallpapers || [];

                    // Pick sensible defaults if not set
                    if (!this.cfg.theme.theme_name) {
                        this.cfg.theme.theme_name = this.themeList[0] || 'Yaru-dark';
                    }
                    // display_font_name: don't auto-select; leave empty = system default
                    // font_name (terminal) is always MesloLGS NF — never changed by the picker
                    if (!this.cfg.wallpaper.background_image && this.wallpapers.length > 0) {
                        this.cfg.wallpaper.background_image = this.wallpapers[0];
                    }
                } catch(e) {
                    console.error('fetchResources error:', e);
                } finally {
                    this.fontsLoading = false;
                    this.wpLoading    = false;
                }
            },

            injectFontFaces(fonts) {
                let css = '';
                fonts.forEach(f => {
                    if (f.files && f.files.length > 0) {
                        const file = f.files[0];
                        // flat_file fonts live directly in themes/fonts/, not in a subdirectory
                        let url = '/api/fonts/file?file=' + encodeURIComponent(file);
                        if (!f.flat_file) {
                            url += '&font=' + encodeURIComponent(f.name);
                        } else {
                            url += '&flat=1';
                        }
                        css += '@font-face { font-family: "' + f.name + '"; src: url("' + url + '"); }\n';
                    }
                });
                if (css) {
                    const s = document.createElement('style');
                    s.textContent = css;
                    document.head.appendChild(s);
                }
            },

            applyRemoteCfg(d) {
                this.cfg.zsh.install_oh_my_zsh  = d.zsh?.install_oh_my_zsh  ?? true;
                this.cfg.git.configure_git      = d.git?.configure_git      ?? true;
                this.cfg.git.git_name           = d.git?.git_name           || '';
                this.cfg.git.git_email          = d.git?.git_email          || '';
                this.cfg.theme.apply_theme      = d.theme?.apply_theme      ?? true;
                this.cfg.theme.theme_mode       = d.theme?.theme_mode       || '1';
                this.cfg.theme.theme_name       = d.theme?.theme_name       || '';
                this.cfg.fonts.configure_fonts  = d.fonts?.configure_fonts  ?? true;
                this.cfg.fonts.font_name         = 'MesloLGS NF'; // always locked
                this.cfg.fonts.display_font_name = d.fonts?.display_font_name || '';
                this.cfg.wallpaper.apply_background = d.wallpaper?.apply_background ?? true;
                this.cfg.docker.enable_docker   = d.docker?.enable_docker   ?? true;
                this.cfg.dock.pin_discord       = d.dock?.pin_discord       ?? true;
                this.cfg.keyboard.configure_keyboard = d.keyboard?.configure_keyboard ?? true;
                this.cfg.keyboard.add_arabic    = d.keyboard?.add_arabic    ?? false;
                this.cfg.power.configure_power  = d.power?.configure_power  ?? true;
                this.cfg.shell.enable_zsh_default = d.shell?.enable_zsh_default ?? false;
                this.cfg.custom_username        = d.custom_username         || '';
                this.cfg.aliases                = d.aliases                 || [];

                // Resolve background image: extract filename from absolute path
                if (d.wallpaper?.background_image) {
                    const name = d.wallpaper.background_image.split('/').pop();
                    if (d.wallpaper.background_image.startsWith('/') && !d.wallpaper.background_image.includes('themes/wallpapers/')) {
                        this.cfg.wallpaper.background_image = d.wallpaper.background_image;
                        this.wpSource = 'custom';
                    } else {
                        this.cfg.wallpaper.background_image = name;
                        this.wpSource = 'repo';
                    }
                }
            },

            // ── NAVIGATION ──
            goTo(i) {
                if (i < 8) this.step = i;
            },
            prev() {
                if (this.step > 0) this.step--;
            },
            next() {
                if (this.step === 7) {
                    this.runSetup();
                } else {
                    this.step++;
                }
            },

            // ── THEME MODE SWITCH ──
            switchThemeMode(mode) {
                this.cfg.theme.theme_mode = mode;
                const list = mode === '1' ? this.darkThemes : this.lightThemes;
                if (list.length && !list.includes(this.cfg.theme.theme_name)) {
                    this.cfg.theme.theme_name = list[0];
                }
            },

            // ── ACTIONS ──
            async loadDefault() {
                try {
                    const r = await fetch('/api/config/default');
                    if (r.ok) {
                        this.applyRemoteCfg(await r.json());
                        this.toast('↺ Reset to defaults');
                    }
                } catch(e) {}
            },

            // Upload a JSON config file from the user's computer via <input type="file">
            async uploadConfig(event) {
                const file = event.target.files[0];
                if (!file) return;
                const form = new FormData();
                form.append('config', file);
                try {
                    const r = await fetch('/api/config/upload', { method: 'POST', body: form });
                    if (r.ok) {
                        const d = await r.json();
                        if (d.status === 'success') {
                            this.applyRemoteCfg(d.config);
                            this.toast('📂 Config imported from file!');
                            this.step = 1;
                        } else {
                            this.toast('⚠ Import failed: ' + (d.message || d.status));
                        }
                    }
                } catch(e) {
                    this.toast('⚠ Upload error: ' + e.message);
                }
                // Reset input so the same file can be re-selected if needed
                event.target.value = '';
            },

            // Load the previously saved config from ~/.config/zonerestore/config.json on the server
            async loadFromDisk() {
                try {
                    const r = await fetch('/api/config/import');
                    if (r.ok) {
                        const d = await r.json();
                        if (d.status === 'success') {
                            this.applyRemoteCfg(d.config);
                            this.toast('💾 Loaded saved settings from disk');
                            this.step = 1;
                        } else if (d.status === 'not_found') {
                            this.toast('ℹ No saved config found on disk');
                        } else {
                            this.toast('⚠ Load failed: ' + (d.message || d.status));
                        }
                    }
                } catch(e) {
                    this.toast('⚠ Error loading from disk: ' + e.message);
                }
            },

            async browseWallpaper() {
                try {
                    const r = await fetch('/api/select-wallpaper');
                    if (r.ok) {
                        const d = await r.json();
                        if (d.status === 'success') {
                            this.cfg.wallpaper.background_image = d.path;
                            this.wpSource = 'custom';
                            this.toast('🖼 Custom wallpaper selected');
                        }
                    }
                } catch(e) {}
            },

            // ── RUN SETUP ──
            async runSetup() {
                this.step = 8;
                this.logs = [];
                this.execFinished = false;
                this.execErr = false;

                // POST config
                try {
                    const body = Object.assign({}, this.cfg);
                    // If using repo gallery, keep just the filename (server resolves absolute path)
                    const res = await fetch('/api/apply', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify(body)
                    });
                    if (!res.ok) throw new Error('Config POST failed: ' + res.status);
                } catch(e) {
                    this.logLine('Error: ' + e.message, 'c-error');
                    this.execFinished = true;
                    this.execErr = true;
                    return;
                }

                // Stream logs via SSE
                const es = new EventSource('/api/stream?export=false');
                const el = document.getElementById('console-output');

                es.onmessage = (ev) => {
                    const text = ev.data;
                    let cls = 'c-line';
                    if (text.includes('[SUCCESS]') || text.startsWith('✔') || text.toLowerCase().includes('success')) cls = 'c-success';
                    else if (text.includes('[ERROR]')   || text.startsWith('✖') || text.toLowerCase().includes('error'))   cls = 'c-error';
                    else if (text.includes('[WARNING]') || text.toLowerCase().includes('warning')) cls = 'c-warn';
                    else if (text.includes('[INFO]')    || text.startsWith('Cloning') || text.startsWith('Applying')) cls = 'c-info';
                    this.logLine(text, cls);
                };

                es.onerror = () => {
                    es.close();
                    this.execFinished = true;
                    const last = this.logs[this.logs.length - 1];
                    this.execErr = !(last && last.text.toLowerCase().includes('finished'));
                };
            },

            logLine(text, cls) {
                this.logs.push({ text, cls: cls || 'c-line' });
                this.$nextTick(() => {
                    const el = document.getElementById('console-output');
                    if (el) el.scrollTop = el.scrollHeight;
                });
            },

            // ── FINISH ──
            async finish(save) {
                this.isClosing = true;

                if (save) {
                    try {
                        // 1. Persist to ~/.config/zonerestore/config.json on disk
                        const saveRes = await fetch('/api/save', {
                            method: 'POST',
                            headers: { 'Content-Type': 'application/json' },
                            body: JSON.stringify(this.cfg)
                        });
                        if (saveRes.ok) {
                            this.toast('💾 Config saved!');
                        }
                        // 2. Trigger browser file download of config JSON
                        const a = document.createElement('a');
                        a.href = '/api/config/download';
                        a.download = 'zonerestore-config.json';
                        document.body.appendChild(a);
                        a.click();
                        document.body.removeChild(a);
                    } catch(e) {
                        this.toast('⚠ Could not save config: ' + e.message);
                    }
                }

                // Signal server to run FinishSetup and shut down
                await fetch('/api/restart', { method: 'POST' }).catch(() => {});

                // Give the user 2s to see the closing screen, then close the tab
                setTimeout(() => {
                    window.close();
                    // Fallback: if window.close() is blocked, show a manual close hint
                    setTimeout(() => {
                        this.toast('ℹ Setup complete — you can close this tab now.');
                    }, 500);
                }, 2000);
            },

            // ── TOAST ──
            toast(msg) {
                const el = document.getElementById('toast');
                document.getElementById('toast-msg').textContent = msg;
                el.classList.add('show');
                setTimeout(() => el.classList.remove('show'), 3200);
            }
        }));
    });
    </script>
</body>
</html>
`
