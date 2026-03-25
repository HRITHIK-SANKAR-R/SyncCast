import { LitElement, html, css } from 'lit';
import { sharedStyles } from './styles.js';
import './sc-status-badge.js';
import './sc-nav.js';

export class ScPlayer extends LitElement {
  static properties = {
    _wsState: { state: true },
    _fileName: { state: true },
  };

  static styles = [
    sharedStyles,
    css`
      :host {
        display: block;
        min-height: 100vh;
        background:
          radial-gradient(1100px 500px at 85% 8%, #0e7490 0%, transparent 60%),
          radial-gradient(900px 500px at 10% 90%, #065f46 0%, transparent 55%),
          linear-gradient(160deg, var(--bg), var(--bg2));
        overflow: hidden;
        display: grid;
        place-items: center;
      }

      .shell {
        width: min(96vw, 1320px);
        padding: 18px;
        animation: fadeIn 420ms ease-out;
      }

      sc-nav {
        margin-bottom: 14px;
      }

      .head {
        display: flex;
        justify-content: space-between;
        align-items: center;
        gap: 12px;
        margin-bottom: 12px;
      }

      .title-row {
        display: flex;
        align-items: center;
        gap: 10px;
      }

      .logo-icon {
        width: 28px;
        height: 28px;
        border-radius: 7px;
        background: linear-gradient(135deg, var(--accent), var(--accent2));
        display: flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0;
      }

      .logo-icon svg {
        width: 15px;
        height: 15px;
        fill: #021a17;
      }

      .title {
        font-size: clamp(17px, 2.1vw, 26px);
        font-weight: 600;
        letter-spacing: -0.01em;
      }

      .meta {
        color: var(--muted);
        font-size: 13px;
        margin-bottom: 12px;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }

      video {
        width: 100%;
        max-height: 80vh;
        border-radius: 12px;
        background: #000;
        outline: none;
        border: 1px solid rgba(255, 255, 255, 0.12);
        box-shadow: 0 0 0 0 rgba(34, 211, 238, 0.45);
        transition: box-shadow 180ms ease;
      }

      video:focus {
        box-shadow: 0 0 0 4px rgba(34, 211, 238, 0.45);
      }
    `,
  ];

  constructor() {
    super();
    this._wsState = 'connecting';
    this._fileName = '';
    this._ws = null;
    this._wakeLockSentinel = null;
    this._fetchInfo();
  }

  connectedCallback() {
    super.connectedCallback();
    this._connectWS();
    document.addEventListener('visibilitychange', this._onVisibilityChange);
  }

  disconnectedCallback() {
    super.disconnectedCallback();
    if (this._ws) { this._ws.close(); this._ws = null; }
    this._releaseWakeLock();
    document.removeEventListener('visibilitychange', this._onVisibilityChange);
  }

  async _fetchInfo() {
    try {
      const res = await fetch('/api/info');
      const data = await res.json();
      this._fileName = data.file || '';
    } catch (_) {}
  }

  get _video() {
    return this.renderRoot?.querySelector('video');
  }

  _wsURL() {
    const proto = location.protocol === 'https:' ? 'wss' : 'ws';
    return `${proto}://${location.host}/ws?role=player`;
  }

  _connectWS() {
    this._wsState = 'connecting';
    this._ws = new WebSocket(this._wsURL());

    this._ws.onopen = () => {
      this._wsState = 'connected';
      const v = this._video;
      if (v) {
        this._sendStatus('ready', {
          duration: Number.isFinite(v.duration) ? v.duration : 0,
          muted: v.muted,
          volume: v.volume,
        });
      }
    };

    this._ws.onclose = () => {
      this._wsState = 'reconnecting';
      setTimeout(() => this._connectWS(), 1000);
    };

    this._ws.onerror = () => { this._wsState = 'error'; };

    this._ws.onmessage = (event) => {
      let msg;
      try { msg = JSON.parse(event.data); } catch (_) { return; }
      if (msg.type !== 'command') return;

      const v = this._video;
      if (!v) return;

      switch (msg.action) {
        case 'play': v.play().catch(() => {}); break;
        case 'pause': v.pause(); break;
        case 'seek':
          if (msg.data && typeof msg.data.position === 'number') {
            v.currentTime = Math.max(0, msg.data.position);
          }
          break;
        case 'back10': v.currentTime = Math.max(0, v.currentTime - 10); break;
        case 'forward10':
          if (Number.isFinite(v.duration) && v.duration > 0) {
            v.currentTime = Math.min(v.duration, v.currentTime + 10);
          } else {
            v.currentTime = Math.max(0, v.currentTime + 10);
          }
          break;
        case 'volume':
          if (msg.data && typeof msg.data.value === 'number') {
            v.volume = Math.min(1, Math.max(0, msg.data.value));
          }
          break;
        case 'stop': v.pause(); v.currentTime = 0; break;
      }
    };
  }

  _sendStatus(action, data) {
    if (!this._ws || this._ws.readyState !== WebSocket.OPEN) return;
    this._ws.send(JSON.stringify({ type: 'status', action, data }));
  }

  _onPlay() {
    const v = this._video;
    this._sendStatus('play', { position: v.currentTime });
    this._requestWakeLock();
  }

  _onPause() {
    const v = this._video;
    this._sendStatus('pause', { position: v.currentTime });
    this._releaseWakeLock();
  }

  _onEnded() {
    const v = this._video;
    this._sendStatus('ended', { position: v.currentTime });
    this._releaseWakeLock();
  }

  _onVolumeChange() {
    const v = this._video;
    this._sendStatus('volume', { volume: v.volume, muted: v.muted });
  }

  _onTimeUpdate() {
    const v = this._video;
    this._sendStatus('time', { position: v.currentTime });
  }

  _onError() {
    const v = this._video;
    this._sendStatus('error', { code: v.error ? v.error.code : 0 });
  }

  _onVisibilityChange = () => {
    const v = this._video;
    if (document.visibilityState === 'visible' && v && !v.paused) {
      this._requestWakeLock();
    }
  };

  async _requestWakeLock() {
    if (!('wakeLock' in navigator)) return;
    if (this._wakeLockSentinel && !this._wakeLockSentinel.released) return;
    try {
      this._wakeLockSentinel = await navigator.wakeLock.request('screen');
      this._wakeLockSentinel.addEventListener('release', () => { this._wakeLockSentinel = null; });
      this._sendStatus('wakelock', { active: true });
    } catch (_) {
      this._sendStatus('wakelock', { active: false });
    }
  }

  async _releaseWakeLock() {
    if (!this._wakeLockSentinel) return;
    try { await this._wakeLockSentinel.release(); } catch (_) {}
    this._wakeLockSentinel = null;
    this._sendStatus('wakelock', { active: false });
  }

  render() {
    return html`
      <main class="shell glass">
        <sc-nav active="player"></sc-nav>
        <header class="head">
          <div class="title-row">
            <div class="logo-icon">
              <svg viewBox="0 0 24 24"><path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/></svg>
            </div>
            <span class="title">SyncCast Player</span>
          </div>
          <sc-status-badge .state="${this._wsState}"></sc-status-badge>
        </header>
        <div class="meta">Now playing: ${this._fileName || 'Loading...'}</div>
        <video
          src="/media"
          controls
          autoplay
          playsinline
          preload="metadata"
          @play="${this._onPlay}"
          @pause="${this._onPause}"
          @ended="${this._onEnded}"
          @volumechange="${this._onVolumeChange}"
          @timeupdate="${this._onTimeUpdate}"
          @error="${this._onError}"
        ></video>
      </main>
    `;
  }
}

customElements.define('sc-player', ScPlayer);
