import { LitElement, html, css } from 'lit';
import { sharedStyles } from './styles.js';
import './sc-status-badge.js';
import './sc-device-card.js';
import './sc-nav.js';

export class ScDashboard extends LitElement {
  static properties = {
    _devices: { state: true },
    _info: { state: true },
    _connState: { state: true },
    _files: { state: true },
    _currentDir: { state: true },
    _loading: { state: true },
    _scanning: { state: true },
    _settingFile: { state: true },
  };

  static styles = [
    sharedStyles,
    css`
      :host {
        display: block;
        min-height: 100vh;
        background:
          radial-gradient(800px 500px at 80% 5%, rgba(20, 184, 166, 0.08), transparent),
          radial-gradient(600px 400px at 10% 90%, rgba(99, 102, 241, 0.06), transparent),
          var(--bg);
      }

      .container {
        max-width: 760px;
        margin: 0 auto;
        padding: 28px 24px 48px;
      }

      sc-nav {
        margin-bottom: 24px;
      }

      /* Header */
      .header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 28px;
        animation: fadeIn 400ms ease;
      }

      .logo {
        display: flex;
        align-items: center;
        gap: 12px;
      }

      .logo-icon {
        width: 38px;
        height: 38px;
        border-radius: 10px;
        background: linear-gradient(135deg, var(--accent), var(--accent2));
        display: flex;
        align-items: center;
        justify-content: center;
      }

      .logo-icon svg { width: 20px; height: 20px; fill: #021a17; }

      .logo-text { font-size: 22px; font-weight: 600; letter-spacing: -0.02em; }
      .logo-sub { font-size: 12px; color: var(--muted); margin-top: -2px; }

      /* Sections */
      .section {
        margin-bottom: 24px;
        animation: fadeIn 500ms ease backwards;
      }

      .section:nth-child(2) { animation-delay: 50ms; }
      .section:nth-child(3) { animation-delay: 100ms; }
      .section:nth-child(4) { animation-delay: 150ms; }
      .section:nth-child(5) { animation-delay: 200ms; }

      .section-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 12px;
      }

      .section-title {
        font-size: 11px;
        text-transform: uppercase;
        letter-spacing: 0.12em;
        color: var(--muted);
        font-weight: 600;
      }

      /* Now Streaming Banner */
      .stream-banner {
        padding: 16px 18px;
        margin-bottom: 24px;
        animation: fadeIn 400ms ease 50ms backwards;
      }

      .stream-row {
        display: flex;
        align-items: center;
        gap: 14px;
      }

      .stream-icon {
        width: 42px;
        height: 42px;
        border-radius: 10px;
        background: linear-gradient(135deg, rgba(20, 184, 166, 0.25), rgba(34, 211, 238, 0.15));
        display: flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0;
      }

      .stream-icon svg {
        width: 20px;
        height: 20px;
        fill: none;
        stroke: var(--accent);
        stroke-width: 2;
        stroke-linecap: round;
        stroke-linejoin: round;
      }

      .stream-info { flex: 1; min-width: 0; }

      .stream-label {
        font-size: 11px;
        text-transform: uppercase;
        letter-spacing: 0.08em;
        color: var(--muted);
        margin-bottom: 2px;
      }

      .stream-file {
        font-size: 15px;
        font-weight: 500;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }

      .stream-idle {
        color: var(--muted);
        font-style: italic;
      }

      /* Info grid */
      .info-panel { padding: 16px 18px; }

      .info-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 12px;
      }

      .info-item { display: flex; flex-direction: column; gap: 2px; }
      .info-label { font-size: 11px; color: var(--muted); text-transform: uppercase; letter-spacing: 0.06em; }
      .info-value { font-size: 13px; font-weight: 500; word-break: break-all; }
      .info-value.mono { font-variant-numeric: tabular-nums; font-family: 'SF Mono', 'Fira Code', monospace; font-size: 12px; }

      /* File Browser */
      .file-browser { max-height: 400px; overflow-y: auto; padding: 4px; }
      .file-browser::-webkit-scrollbar { width: 5px; }
      .file-browser::-webkit-scrollbar-track { background: transparent; }
      .file-browser::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.1); border-radius: 3px; }

      .dir-path {
        padding: 8px 12px;
        font-size: 12px;
        font-family: 'SF Mono', 'Fira Code', monospace;
        color: var(--muted);
        background: rgba(0,0,0,0.2);
        border-radius: var(--radius-sm) var(--radius-sm) 0 0;
        border-bottom: 1px solid rgba(255,255,255,0.04);
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }

      .file-item {
        display: flex;
        align-items: center;
        gap: 10px;
        padding: 10px 14px;
        cursor: pointer;
        border-bottom: 1px solid rgba(255,255,255,0.03);
        transition: background var(--transition);
        -webkit-tap-highlight-color: transparent;
      }

      .file-item:last-child { border-bottom: none; }
      .file-item:hover { background: rgba(255,255,255,0.05); }
      .file-item:active { background: rgba(20, 184, 166, 0.08); }

      .file-icon {
        width: 18px;
        height: 18px;
        fill: none;
        stroke: var(--muted);
        stroke-width: 2;
        stroke-linecap: round;
        stroke-linejoin: round;
        flex-shrink: 0;
      }

      .file-item[data-dir] .file-icon { stroke: var(--accent2); }
      .file-item[data-media] .file-icon { stroke: var(--accent); }

      .file-name {
        flex: 1;
        font-size: 13px;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }

      .file-size {
        font-size: 11px;
        color: var(--muted);
        flex-shrink: 0;
        font-variant-numeric: tabular-nums;
      }

      .file-loading {
        text-align: center;
        padding: 24px;
        color: var(--muted);
        font-size: 13px;
      }

      /* Devices */
      .device-list { display: flex; flex-direction: column; gap: 8px; }

      .empty-msg {
        text-align: center;
        padding: 24px 16px;
        color: var(--muted);
        font-size: 13px;
      }

      .empty-msg .hint { font-size: 12px; margin-top: 4px; opacity: 0.6; }

      /* Action buttons */
      .action-btn {
        padding: 8px 16px;
        font-size: 12px;
        font-weight: 600;
        gap: 6px;
      }

      .action-btn svg {
        width: 14px; height: 14px;
        fill: none; stroke: currentColor;
        stroke-width: 2; stroke-linecap: round; stroke-linejoin: round;
      }

      .spinning svg { animation: spin 800ms linear infinite; }

      @keyframes spin { to { transform: rotate(360deg); } }

      /* URLs */
      .url-item {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 8px 0;
        border-bottom: 1px solid rgba(255,255,255,0.04);
      }
      .url-item:last-child { border-bottom: none; }
      .url-label { font-size: 12px; color: var(--muted); min-width: 70px; }
      .url-value { font-size: 12px; font-family: 'SF Mono', 'Fira Code', monospace; color: var(--accent); word-break: break-all; }
    `,
  ];

  constructor() {
    super();
    this._devices = [];
    this._info = {};
    this._connState = { state: 'idle', remote_clients: 0, player_clients: 0 };
    this._files = [];
    this._currentDir = '';
    this._loading = false;
    this._scanning = false;
    this._settingFile = '';
    this._pollTimer = null;
  }

  connectedCallback() {
    super.connectedCallback();
    this._fetchAll();
    this._pollTimer = setInterval(() => this._fetchState(), 3000);
  }

  disconnectedCallback() {
    super.disconnectedCallback();
    if (this._pollTimer) clearInterval(this._pollTimer);
  }

  async _fetchAll() {
    this._loading = true;
    await Promise.all([this._fetchInfo(), this._fetchDevices(), this._fetchState()]);
    // Load home directory for file browser
    if (!this._currentDir && this._info.home) {
      await this._fetchFiles(this._info.home);
    }
    this._loading = false;
  }

  async _fetchInfo() {
    try {
      const res = await fetch('/api/info');
      this._info = await res.json();
    } catch (_) {}
  }

  async _fetchDevices() {
    try {
      const res = await fetch('/api/devices');
      this._devices = await res.json();
    } catch (_) {}
  }

  async _fetchState() {
    try {
      const res = await fetch('/api/state');
      this._connState = await res.json();

      // Also refresh info to catch file changes
      await this._fetchInfo();
    } catch (_) {}
  }

  async _fetchFiles(dir) {
    try {
      const res = await fetch(`/api/files?dir=${encodeURIComponent(dir)}`);
      const data = await res.json();
      if (data.error) return;
      this._currentDir = data.dir;
      this._files = data.files || [];
    } catch (_) {}
  }

  async _onFileClick(file) {
    if (file.is_dir) {
      await this._fetchFiles(file.path);
      return;
    }
    // It's a media file — start streaming
    this._settingFile = file.path;
    try {
      const res = await fetch('/api/stream', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ file: file.path }),
      });
      const data = await res.json();
      if (data.ok) {
        await this._fetchInfo();
      }
    } catch (_) {}
    this._settingFile = '';
  }

  async _onScan() {
    this._scanning = true;
    try {
      const res = await fetch('/api/discover', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({}),
      });
      const data = await res.json();
      if (data.devices) {
        this._devices = data.devices;
      }
    } catch (_) {}
    this._scanning = false;
  }

  _buildURL(path) {
    const host = this._info.host || location.hostname;
    const port = this._info.port || location.port;
    return `http://${host}:${port}${path}`;
  }

  _formatSize(bytes) {
    if (!bytes) return '';
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1048576) return (bytes / 1024).toFixed(0) + ' KB';
    if (bytes < 1073741824) return (bytes / 1048576).toFixed(1) + ' MB';
    return (bytes / 1073741824).toFixed(2) + ' GB';
  }

  _dirIcon() {
    return html`<svg class="file-icon" viewBox="0 0 24 24"><path d="M22 19a2 2 0 01-2 2H4a2 2 0 01-2-2V5a2 2 0 012-2h5l2 3h9a2 2 0 012 2z"/></svg>`;
  }

  _mediaIcon() {
    return html`<svg class="file-icon" viewBox="0 0 24 24"><polygon points="5 3 19 12 5 21 5 3"/></svg>`;
  }

  _upIcon() {
    return html`<svg class="file-icon" viewBox="0 0 24 24"><polyline points="15 18 9 12 15 6"/></svg>`;
  }

  render() {
    const streaming = this._info.streaming;
    const fileName = this._info.file || '';

    return html`
      <div class="container">
        <!-- Header -->
        <div class="header">
          <div class="logo">
            <div class="logo-icon">
              <svg viewBox="0 0 24 24"><path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/></svg>
            </div>
            <div>
              <div class="logo-text">SyncCast</div>
              <div class="logo-sub">Dashboard</div>
            </div>
          </div>
          <sc-status-badge .state="${this._connState.state || 'idle'}"></sc-status-badge>
        </div>

        <sc-nav active="dashboard"></sc-nav>

        <!-- Now Streaming Banner -->
        <div class="stream-banner glass">
          <div class="stream-row">
            <div class="stream-icon">
              ${streaming
                ? html`<svg viewBox="0 0 24 24"><polygon points="5 3 19 12 5 21 5 3"/></svg>`
                : html`<svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>`
              }
            </div>
            <div class="stream-info">
              <div class="stream-label">${streaming ? 'Now Streaming' : 'No File Selected'}</div>
              <div class="stream-file ${streaming ? '' : 'stream-idle'}">
                ${streaming ? fileName : 'Pick a video below to start streaming'}
              </div>
            </div>
          </div>
        </div>

        <!-- File Browser -->
        <div class="section">
          <div class="section-header">
            <div class="section-title">Select a Video</div>
          </div>
          <div class="glass" style="overflow:hidden;">
            ${this._currentDir ? html`<div class="dir-path">${this._currentDir}</div>` : ''}
            <div class="file-browser">
              ${this._files.length === 0
                ? html`<div class="file-loading">${this._loading ? 'Loading...' : 'No media files here'}</div>`
                : this._files.map(f => html`
                  <div
                    class="file-item"
                    ?data-dir="${f.is_dir}"
                    ?data-media="${!f.is_dir}"
                    @click="${() => this._onFileClick(f)}"
                  >
                    ${f.is_dir
                      ? (f.name === '..' ? this._upIcon() : this._dirIcon())
                      : this._mediaIcon()
                    }
                    <span class="file-name">${f.name === '..' ? 'Back' : f.name}</span>
                    ${!f.is_dir ? html`<span class="file-size">${this._formatSize(f.size)}</span>` : ''}
                    ${this._settingFile === f.path ? html`<span class="file-size" style="color:var(--accent)">Setting...</span>` : ''}
                  </div>
                `)
              }
            </div>
          </div>
        </div>

        <!-- Devices -->
        <div class="section">
          <div class="section-header">
            <div class="section-title">Devices (${this._devices.length})</div>
            <button class="btn action-btn ${this._scanning ? 'spinning' : ''}" @click="${this._onScan}">
              <svg viewBox="0 0 24 24"><path d="M1 1l22 22"/><path d="M16.72 11.06A10.94 10.94 0 0119 12.55"/><path d="M5 12.55a10.94 10.94 0 015.17-2.39"/><path d="M10.71 5.05A16 16 0 0122.56 9"/><path d="M1.42 9a15.91 15.91 0 014.7-2.88"/><path d="M8.53 16.11a6 6 0 016.95 0"/><line x1="12" y1="20" x2="12.01" y2="20"/></svg>
              ${this._scanning ? 'Scanning...' : 'Scan Network'}
            </button>
          </div>
          ${this._devices.length > 0
            ? html`<div class="device-list">${this._devices.map((d, i) =>
                html`<sc-device-card style="animation-delay:${i * 60}ms" .name="${d.name}" .ip="${d.ip}"></sc-device-card>`
              )}</div>`
            : html`<div class="empty-msg">
                No devices discovered yet
                <div class="hint">Click "Scan Network" or start with --ip flag</div>
              </div>`
          }
        </div>

        <!-- Server Info -->
        <div class="section">
          <div class="section-header">
            <div class="section-title">Server & URLs</div>
          </div>
          <div class="info-panel glass">
            <div class="info-grid" style="margin-bottom:12px;">
              <div class="info-item">
                <span class="info-label">Host</span>
                <span class="info-value mono">${this._info.host || '—'}:${this._info.port || '—'}</span>
              </div>
              <div class="info-item">
                <span class="info-label">Clients</span>
                <span class="info-value">${this._connState.remote_clients ?? 0} remote, ${this._connState.player_clients ?? 0} player</span>
              </div>
            </div>
            <div class="url-item">
              <span class="url-label">Remote</span>
              <span class="url-value">${this._buildURL('/remote')}</span>
            </div>
            <div class="url-item">
              <span class="url-label">Player</span>
              <span class="url-value">${this._buildURL('/player')}</span>
            </div>
          </div>
        </div>
      </div>
    `;
  }
}

customElements.define('sc-dashboard', ScDashboard);
