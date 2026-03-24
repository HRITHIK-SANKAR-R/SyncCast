import { LitElement, html, css } from 'lit';
import { sharedStyles } from './styles.js';

export class ScDeviceCard extends LitElement {
  static properties = {
    name: { type: String },
    ip: { type: String },
  };

  static styles = [
    sharedStyles,
    css`
      :host {
        display: block;
      }

      .card {
        display: flex;
        align-items: center;
        gap: 14px;
        padding: 14px 18px;
        border-radius: var(--radius-sm);
        background: rgba(255, 255, 255, 0.04);
        border: 1px solid rgba(255, 255, 255, 0.06);
        transition: background var(--transition), border-color var(--transition);
        animation: fadeIn 350ms ease backwards;
      }

      .card:hover {
        background: rgba(255, 255, 255, 0.07);
        border-color: rgba(20, 184, 166, 0.25);
      }

      .icon {
        width: 40px;
        height: 40px;
        border-radius: 10px;
        background: linear-gradient(135deg, rgba(20, 184, 166, 0.2), rgba(34, 211, 238, 0.1));
        display: flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0;
      }

      .icon svg {
        width: 20px;
        height: 20px;
        fill: none;
        stroke: var(--accent);
        stroke-width: 1.8;
        stroke-linecap: round;
        stroke-linejoin: round;
      }

      .info {
        flex: 1;
        min-width: 0;
      }

      .name {
        font-size: 14px;
        font-weight: 500;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }

      .ip {
        font-size: 12px;
        color: var(--muted);
        margin-top: 2px;
        font-variant-numeric: tabular-nums;
      }
    `,
  ];

  constructor() {
    super();
    this.name = 'Unknown Device';
    this.ip = '';
  }

  render() {
    return html`
      <div class="card">
        <div class="icon">
          <svg viewBox="0 0 24 24">
            <rect x="2" y="3" width="20" height="14" rx="2" ry="2"/>
            <line x1="8" y1="21" x2="16" y2="21"/>
            <line x1="12" y1="17" x2="12" y2="21"/>
          </svg>
        </div>
        <div class="info">
          <div class="name">${this.name}</div>
          <div class="ip">${this.ip}</div>
        </div>
      </div>
    `;
  }
}

customElements.define('sc-device-card', ScDeviceCard);
