import { LitElement, html, css } from 'lit';
import { sharedStyles } from './styles.js';

export class ScStatusBadge extends LitElement {
  static properties = {
    state: { type: String },
  };

  static styles = [
    sharedStyles,
    css`
      :host {
        display: inline-block;
      }

      .badge {
        display: inline-flex;
        align-items: center;
        gap: 6px;
        padding: 4px 12px;
        border-radius: 999px;
        font-size: 11px;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.08em;
        transition: background var(--transition), color var(--transition), border-color var(--transition);
      }

      .dot {
        width: 7px;
        height: 7px;
        border-radius: 50%;
        flex-shrink: 0;
      }

      /* States */
      .badge[data-state="connected"] {
        background: rgba(20, 184, 166, 0.15);
        color: #5eead4;
        border: 1px solid rgba(20, 184, 166, 0.35);
      }
      .badge[data-state="connected"] .dot {
        background: #14b8a6;
        box-shadow: 0 0 6px #14b8a6;
      }

      .badge[data-state="reconnecting"] {
        background: rgba(251, 191, 36, 0.15);
        color: #fde68a;
        border: 1px solid rgba(251, 191, 36, 0.35);
      }
      .badge[data-state="reconnecting"] .dot {
        background: #fbbf24;
        animation: pulse 1.2s infinite;
      }

      .badge[data-state="idle"] {
        background: rgba(148, 163, 184, 0.1);
        color: #94a3b8;
        border: 1px solid rgba(148, 163, 184, 0.2);
      }
      .badge[data-state="idle"] .dot {
        background: #64748b;
      }

      .badge[data-state="waiting_remote"],
      .badge[data-state="waiting_player"] {
        background: rgba(99, 102, 241, 0.15);
        color: #a5b4fc;
        border: 1px solid rgba(99, 102, 241, 0.35);
      }
      .badge[data-state="waiting_remote"] .dot,
      .badge[data-state="waiting_player"] .dot {
        background: #818cf8;
        animation: pulse 1.5s infinite;
      }

      .badge[data-state="error"] {
        background: rgba(244, 63, 94, 0.15);
        color: #fda4af;
        border: 1px solid rgba(244, 63, 94, 0.35);
      }
      .badge[data-state="error"] .dot {
        background: #f43f5e;
      }
    `,
  ];

  constructor() {
    super();
    this.state = 'idle';
  }

  _label() {
    const labels = {
      idle: 'Idle',
      connected: 'Connected',
      reconnecting: 'Reconnecting',
      waiting_remote: 'Waiting for remote',
      waiting_player: 'Waiting for player',
      error: 'Error',
      connecting: 'Connecting',
    };
    return labels[this.state] || this.state;
  }

  render() {
    return html`
      <span class="badge" data-state="${this.state}">
        <span class="dot"></span>
        ${this._label()}
      </span>
    `;
  }
}

customElements.define('sc-status-badge', ScStatusBadge);
