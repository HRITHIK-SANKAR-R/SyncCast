import { css } from 'lit';

export const sharedStyles = css`
  :host {
    --bg: #0f172a;
    --bg2: #111827;
    --panel: rgba(17, 24, 39, 0.82);
    --panel-border: rgba(255, 255, 255, 0.08);
    --text: #f8fafc;
    --muted: #94a3b8;
    --accent: #14b8a6;
    --accent-glow: rgba(20, 184, 166, 0.25);
    --accent2: #22d3ee;
    --danger: #f43f5e;
    --radius: 14px;
    --radius-sm: 8px;
    --transition: 180ms ease;

    font-family: 'Inter', system-ui, -apple-system, sans-serif;
    color: var(--text);
  }

  .glass {
    background: var(--panel);
    backdrop-filter: blur(18px);
    -webkit-backdrop-filter: blur(18px);
    border: 1px solid var(--panel-border);
    border-radius: var(--radius);
  }

  .btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border: none;
    background: rgba(255, 255, 255, 0.06);
    color: var(--text);
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: background var(--transition), transform var(--transition), box-shadow var(--transition);
    font-family: inherit;
    -webkit-tap-highlight-color: transparent;
    user-select: none;
  }

  .btn:hover {
    background: rgba(255, 255, 255, 0.12);
  }

  .btn:active {
    transform: scale(0.94);
  }

  .btn-accent {
    background: var(--accent);
    color: #021a17;
    font-weight: 600;
  }

  .btn-accent:hover {
    background: #0fd9a4;
    box-shadow: 0 0 20px var(--accent-glow);
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
  }

  @keyframes pulse {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.5; }
  }
`;
