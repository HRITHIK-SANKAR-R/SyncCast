import { LitElement, html, css } from 'lit';
import { sharedStyles } from './styles.js';

export class ScNav extends LitElement {
  static properties = {
    active: { type: String },
  };

  static styles = [
    sharedStyles,
    css`
      :host {
        display: block;
      }

      nav {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 4px;
        padding: 6px;
        border-radius: var(--radius);
        background: rgba(15, 23, 42, 0.7);
        backdrop-filter: blur(12px);
        border: 1px solid rgba(255, 255, 255, 0.08);
        max-width: 340px;
        margin: 0 auto;
      }

      a {
        display: flex;
        align-items: center;
        gap: 6px;
        padding: 8px 16px;
        border-radius: 10px;
        font-size: 13px;
        font-weight: 500;
        color: var(--muted);
        text-decoration: none;
        transition: all var(--transition);
        -webkit-tap-highlight-color: transparent;
        white-space: nowrap;
      }

      a svg {
        width: 16px;
        height: 16px;
        fill: none;
        stroke: currentColor;
        stroke-width: 2;
        stroke-linecap: round;
        stroke-linejoin: round;
        flex-shrink: 0;
      }

      a:hover {
        color: var(--text);
        background: rgba(255, 255, 255, 0.06);
      }

      a.active {
        color: var(--text);
        background: rgba(20, 184, 166, 0.15);
        border: 1px solid rgba(20, 184, 166, 0.25);
      }

      a.active svg {
        stroke: var(--accent);
      }

      @media (max-width: 400px) {
        a span { display: none; }
        a { padding: 10px 14px; }
        nav { gap: 2px; }
      }
    `,
  ];

  constructor() {
    super();
    this.active = '';
  }

  render() {
    return html`
      <nav>
        <a href="/dashboard" class="${this.active === 'dashboard' ? 'active' : ''}">
          <svg viewBox="0 0 24 24"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/></svg>
          <span>Dashboard</span>
        </a>
        <a href="/remote" class="${this.active === 'remote' ? 'active' : ''}">
          <svg viewBox="0 0 24 24"><rect x="5" y="2" width="14" height="20" rx="2" ry="2"/><line x1="12" y1="18" x2="12.01" y2="18"/></svg>
          <span>Remote</span>
        </a>
        <a href="/player" class="${this.active === 'player' ? 'active' : ''}">
          <svg viewBox="0 0 24 24"><rect x="2" y="7" width="20" height="15" rx="2" ry="2"/><polyline points="17 2 12 7 7 2"/></svg>
          <span>Player</span>
        </a>
      </nav>
    `;
  }
}

customElements.define('sc-nav', ScNav);
