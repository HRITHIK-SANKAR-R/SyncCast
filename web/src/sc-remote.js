import { LitElement, html, css } from 'lit';
import { sharedStyles } from './styles.js';
import './sc-status-badge.js';
import './sc-nav.js';

export class ScRemote extends LitElement {
  static properties = {
    _wsState: { state: true },
    _playing: { state: true },
    _position: { state: true },
    _duration: { state: true },
    _volume: { state: true },
    _muted: { state: true },
    _fileName: { state: true },
    _seeking: { state: true },
  };

  static styles = [
    sharedStyles,
    css`
      :host {
        display: block;
        min-height: 100vh;
        min-height: 100dvh;
        background:
          radial-gradient(600px 400px at 70% 10%, rgba(20, 184, 166, 0.12), transparent),
          radial-gradient(500px 400px at 20% 85%, rgba(99, 102, 241, 0.08), transparent),
          var(--bg);
      }

      .container {
        max-width: 420px;
        margin: 0 auto;
        padding: 20px 20px 40px;
        display: flex;
        flex-direction: column;
        min-height: 100vh;
        min-height: 100dvh;
      }

      sc-nav {
        margin-bottom: 20px;
      }

      /* Header */
      .header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 24px;
        animation: fadeIn 400ms ease;
      }

      .logo {
        display: flex;
        align-items: center;
        gap: 10px;
      }

      .logo-icon {
        width: 32px;
        height: 32px;
        border-radius: 8px;
        background: linear-gradient(135deg, var(--accent), var(--accent2));
        display: flex;
        align-items: center;
        justify-content: center;
      }

      .logo-icon svg {
        width: 18px;
        height: 18px;
        fill: #021a17;
      }

      .logo-text {
        font-size: 18px;
        font-weight: 600;
        letter-spacing: -0.02em;
      }

      /* Now Playing */
      .now-playing {
        padding: 16px 18px;
        margin-bottom: 20px;
        animation: fadeIn 450ms ease 50ms backwards;
      }

      .np-label {
        font-size: 11px;
        text-transform: uppercase;
        letter-spacing: 0.1em;
        color: var(--muted);
        margin-bottom: 6px;
      }

      .np-title {
        font-size: 15px;
        font-weight: 500;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }

      /* Progress */
      .progress-section {
        margin-bottom: 28px;
        animation: fadeIn 500ms ease 100ms backwards;
      }

      .progress-bar-track {
        width: 100%;
        height: 6px;
        border-radius: 3px;
        background: rgba(255, 255, 255, 0.08);
        cursor: pointer;
        position: relative;
        transition: height 120ms ease;
      }

      .progress-bar-track:hover {
        height: 10px;
      }

      .progress-bar-fill {
        height: 100%;
        border-radius: 3px;
        background: linear-gradient(90deg, var(--accent), var(--accent2));
        transition: width 300ms linear;
        position: relative;
      }

      .progress-bar-fill::after {
        content: '';
        position: absolute;
        right: -5px;
        top: 50%;
        transform: translateY(-50%) scale(0);
        width: 14px;
        height: 14px;
        border-radius: 50%;
        background: var(--text);
        box-shadow: 0 0 8px rgba(0,0,0,0.4);
        transition: transform 120ms ease;
      }

      .progress-bar-track:hover .progress-bar-fill::after {
        transform: translateY(-50%) scale(1);
      }

      .time-row {
        display: flex;
        justify-content: space-between;
        margin-top: 8px;
        font-size: 12px;
        color: var(--muted);
        font-variant-numeric: tabular-nums;
      }

      /* Controls */
      .controls {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 12px;
        margin-bottom: 32px;
        animation: fadeIn 550ms ease 150ms backwards;
      }

      .ctrl-btn {
        width: 52px;
        height: 52px;
        border-radius: 50%;
        font-size: 0;
      }

      .ctrl-btn svg {
        width: 22px;
        height: 22px;
        fill: none;
        stroke: var(--text);
        stroke-width: 2;
        stroke-linecap: round;
        stroke-linejoin: round;
      }

      .ctrl-btn.play-btn {
        width: 72px;
        height: 72px;
        background: var(--accent);
        box-shadow: 0 0 28px var(--accent-glow);
      }

      .ctrl-btn.play-btn:hover {
        background: #0fd9a4;
        box-shadow: 0 0 36px var(--accent-glow);
      }

      .ctrl-btn.play-btn svg {
        width: 30px;
        height: 30px;
        stroke: #021a17;
        fill: #021a17;
        stroke-width: 1;
      }

      .skip-label {
        font-size: 10px;
        color: var(--muted);
        position: absolute;
        bottom: -18px;
        width: 100%;
        text-align: center;
      }

      .ctrl-wrap {
        position: relative;
      }

      /* Volume */
      .volume-section {
        padding: 16px 18px;
        animation: fadeIn 600ms ease 200ms backwards;
      }

      .volume-row {
        display: flex;
        align-items: center;
        gap: 12px;
      }

      .vol-icon {
        width: 20px;
        height: 20px;
        fill: none;
        stroke: var(--muted);
        stroke-width: 2;
        stroke-linecap: round;
        stroke-linejoin: round;
        flex-shrink: 0;
        cursor: pointer;
        transition: stroke var(--transition);
      }

      .vol-icon:hover {
        stroke: var(--text);
      }

      .volume-slider {
        flex: 1;
        -webkit-appearance: none;
        appearance: none;
        height: 4px;
        border-radius: 2px;
        background: rgba(255, 255, 255, 0.1);
        outline: none;
      }

      .volume-slider::-webkit-slider-thumb {
        -webkit-appearance: none;
        width: 16px;
        height: 16px;
        border-radius: 50%;
        background: var(--accent);
        cursor: pointer;
        box-shadow: 0 0 8px var(--accent-glow);
      }

      .volume-slider::-moz-range-thumb {
        width: 16px;
        height: 16px;
        border-radius: 50%;
        background: var(--accent);
        cursor: pointer;
        border: none;
        box-shadow: 0 0 8px var(--accent-glow);
      }

      .vol-value {
        font-size: 12px;
        color: var(--muted);
        min-width: 32px;
        text-align: right;
        font-variant-numeric: tabular-nums;
      }

      /* Stop */
      .stop-section {
        margin-top: auto;
        padding-top: 20px;
        animation: fadeIn 650ms ease 250ms backwards;
      }

      .stop-btn {
        width: 100%;
        padding: 14px;
        border-radius: var(--radius-sm);
        font-size: 13px;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.05em;
        background: rgba(244, 63, 94, 0.1);
        border: 1px solid rgba(244, 63, 94, 0.2);
        color: #fda4af;
        cursor: pointer;
        transition: background var(--transition), border-color var(--transition);
        font-family: inherit;
        -webkit-tap-highlight-color: transparent;
      }

      .stop-btn:hover {
        background: rgba(244, 63, 94, 0.2);
        border-color: rgba(244, 63, 94, 0.4);
      }

      .stop-btn:active {
        transform: scale(0.98);
      }
    `,
  ];

  constructor() {
    super();
    this._wsState = 'connecting';
    this._playing = false;
    this._position = 0;
    this._duration = 0;
    this._volume = 1;
    this._muted = false;
    this._fileName = '';
    this._seeking = false;
    this._ws = null;
    this._fetchInfo();
  }

  connectedCallback() {
    super.connectedCallback();
    this._connectWS();
  }

  disconnectedCallback() {
    super.disconnectedCallback();
    if (this._ws) {
      this._ws.close();
      this._ws = null;
    }
  }

  async _fetchInfo() {
    try {
      const res = await fetch('/api/info');
      const data = await res.json();
      this._fileName = data.file || '';
    } catch (_) {}
  }

  _wsURL() {
    const proto = location.protocol === 'https:' ? 'wss' : 'ws';
    return `${proto}://${location.host}/ws?role=remote`;
  }

  _connectWS() {
    this._wsState = 'connecting';
    this._ws = new WebSocket(this._wsURL());

    this._ws.onopen = () => {
      this._wsState = 'connected';
    };

    this._ws.onclose = () => {
      this._wsState = 'reconnecting';
      setTimeout(() => this._connectWS(), 1500);
    };

    this._ws.onerror = () => {
      this._wsState = 'error';
    };

    this._ws.onmessage = (event) => {
      let msg;
      try { msg = JSON.parse(event.data); } catch (_) { return; }
      if (msg.type !== 'status') return;

      switch (msg.action) {
        case 'time':
          if (!this._seeking) this._position = msg.data?.position ?? this._position;
          break;
        case 'play':
          this._playing = true;
          this._position = msg.data?.position ?? this._position;
          break;
        case 'pause':
          this._playing = false;
          this._position = msg.data?.position ?? this._position;
          break;
        case 'ended':
          this._playing = false;
          break;
        case 'volume':
          this._volume = msg.data?.volume ?? this._volume;
          this._muted = msg.data?.muted ?? this._muted;
          break;
        case 'ready':
          this._duration = msg.data?.duration ?? 0;
          this._volume = msg.data?.volume ?? 1;
          this._muted = msg.data?.muted ?? false;
          break;
      }
    };
  }

  _send(action, data = {}) {
    if (!this._ws || this._ws.readyState !== WebSocket.OPEN) return;
    this._ws.send(JSON.stringify({ type: 'command', action, data }));
  }

  _togglePlay() {
    this._send(this._playing ? 'pause' : 'play');
  }

  _back10() {
    this._send('back10');
  }

  _forward10() {
    this._send('forward10');
  }

  _onSeek(e) {
    const track = e.currentTarget;
    const rect = track.getBoundingClientRect();
    const ratio = Math.max(0, Math.min(1, (e.clientX - rect.left) / rect.width));
    const pos = ratio * this._duration;
    this._position = pos;
    this._send('seek', { position: pos });
  }

  _onVolumeInput(e) {
    const val = parseFloat(e.target.value);
    this._volume = val;
    this._muted = false;
    this._send('volume', { value: val });
  }

  _toggleMute() {
    this._muted = !this._muted;
    this._send('volume', { value: this._muted ? 0 : this._volume });
  }

  _stop() {
    this._send('stop');
    this._playing = false;
    this._position = 0;
  }

  _fmt(seconds) {
    if (!seconds || !Number.isFinite(seconds)) return '0:00';
    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  }

  _progressPct() {
    if (!this._duration || this._duration <= 0) return 0;
    return Math.min(100, (this._position / this._duration) * 100);
  }

  render() {
    return html`
      <div class="container">
        <!-- Header -->
        <div class="header">
          <div class="logo">
            <div class="logo-icon">
              <svg viewBox="0 0 24 24"><path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/></svg>
            </div>
            <span class="logo-text">SyncCast</span>
          </div>
          <sc-status-badge .state="${this._wsState}"></sc-status-badge>
        </div>

        <sc-nav active="remote"></sc-nav>

        <!-- Now Playing -->
        <div class="now-playing glass">
          <div class="np-label">Now Playing</div>
          <div class="np-title">${this._fileName || 'No media loaded'}</div>
        </div>

        <!-- Progress -->
        <div class="progress-section">
          <div class="progress-bar-track" @click="${this._onSeek}">
            <div class="progress-bar-fill" style="width: ${this._progressPct()}%"></div>
          </div>
          <div class="time-row">
            <span>${this._fmt(this._position)}</span>
            <span>${this._fmt(this._duration)}</span>
          </div>
        </div>

        <!-- Controls -->
        <div class="controls">
          <div class="ctrl-wrap">
            <button class="btn ctrl-btn" @click="${this._back10}" title="Back 10s">
              <svg viewBox="0 0 24 24">
                <path d="M12.5 8V4L6 9.5l6.5 5.5v-4c3.6 0 6.5 2.9 6.5 6.5S16.1 23 12.5 23 6 20.1 6 16.5"/>
                <text x="10" y="15" fill="currentColor" stroke="none" font-size="7" font-weight="700" font-family="Inter, sans-serif">10</text>
              </svg>
            </button>
          </div>

          <button class="btn ctrl-btn play-btn" @click="${this._togglePlay}" title="${this._playing ? 'Pause' : 'Play'}">
            ${this._playing
              ? html`<svg viewBox="0 0 24 24"><rect x="6" y="4" width="4" height="16" rx="1"/><rect x="14" y="4" width="4" height="16" rx="1"/></svg>`
              : html`<svg viewBox="0 0 24 24"><polygon points="6,3 20,12 6,21"/></svg>`
            }
          </button>

          <div class="ctrl-wrap">
            <button class="btn ctrl-btn" @click="${this._forward10}" title="Forward 10s">
              <svg viewBox="0 0 24 24">
                <path d="M11.5 8V4l6.5 5.5-6.5 5.5v-4C8 11 5 13.9 5 17.5S8 24 11.5 24s6.5-2.9 6.5-6.5"/>
                <text x="8.5" y="15" fill="currentColor" stroke="none" font-size="7" font-weight="700" font-family="Inter, sans-serif">10</text>
              </svg>
            </button>
          </div>
        </div>

        <!-- Volume -->
        <div class="volume-section glass">
          <div class="volume-row">
            <svg class="vol-icon" viewBox="0 0 24 24" @click="${this._toggleMute}">
              ${this._muted || this._volume === 0
                ? html`<polygon points="11,5 6,9 2,9 2,15 6,15 11,19"/><line x1="23" y1="9" x2="17" y2="15"/><line x1="17" y1="9" x2="23" y2="15"/>`
                : html`<polygon points="11,5 6,9 2,9 2,15 6,15 11,19"/><path d="M19.07 4.93a10 10 0 0 1 0 14.14"/><path d="M15.54 8.46a5 5 0 0 1 0 7.07"/>`
              }
            </svg>
            <input
              type="range"
              class="volume-slider"
              min="0"
              max="1"
              step="0.02"
              .value="${String(this._muted ? 0 : this._volume)}"
              @input="${this._onVolumeInput}"
            />
            <span class="vol-value">${Math.round((this._muted ? 0 : this._volume) * 100)}%</span>
          </div>
        </div>

        <!-- Stop -->
        <div class="stop-section">
          <button class="stop-btn" @click="${this._stop}">Stop Playback</button>
        </div>
      </div>
    `;
  }
}

customElements.define('sc-remote', ScRemote);
