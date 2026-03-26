import { LitElement, html, css } from 'lit';

// Lazy-load views
import './sc-remote.js';
import './sc-dashboard.js';
import './sc-player.js';

export class ScApp extends LitElement {
  static properties = {
    _view: { state: true },
  };

  static styles = css`
    :host {
      display: block;
    }
  `;

  constructor() {
    super();
    this._view = this._detectView();
  }

  _detectView() {
    const path = window.location.pathname;
    if (path.startsWith('/remote')) return 'remote';
    if (path.startsWith('/dashboard')) return 'dashboard';
    if (path.startsWith('/player')) return 'player';
    // Default: detect by screen width
    if (window.innerWidth <= 768) return 'remote';
    return 'dashboard';
  }

  render() {
    switch (this._view) {
      case 'remote':
        return html`<sc-remote></sc-remote>`;
      case 'player':
        return html`<sc-player></sc-player>`;
      case 'dashboard':
      default:
        return html`<sc-dashboard></sc-dashboard>`;
    }
  }
}

customElements.define('sc-app', ScApp);
