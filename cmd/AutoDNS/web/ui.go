package web

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>AutoDNS</title>
<link rel="icon" type="image/svg+xml" href="data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0ibm9uZSIgc3Ryb2tlPSIjM2I4MmY2IiBzdHJva2Utd2lkdGg9IjIiIHN0cm9rZS1saW5lY2FwPSJyb3VuZCIgc3Ryb2tlLWxpbmVqb2luPSJyb3VuZCI+PGNpcmNsZSBjeD0iMTIiIGN5PSIxMiIgcj0iMTAiLz48bGluZSB4MT0iMiIgeTE9IjEyIiB4Mj0iMjIiIHkyPSIxMiIvPjxwYXRoIGQ9Ik0xMiAyYTE1LjMgMTUuMyAwIDAgMSA0IDEwIDE1LjMgMTUuMyAwIDAgMS00IDEwIDE1LjMgMTUuMyAwIDAgMS00LTEwIDE1LjMgMTUuMyAwIDAgMSA0LTEweiIvPjwvc3ZnPg==">
<style>
*{box-sizing:border-box;margin:0;padding:0}
body{font-family:system-ui,-apple-system,sans-serif;background:#0f1117;color:#e2e8f0;min-height:100vh}
.topbar{background:#1a1d27;border-bottom:1px solid #2d3148;padding:14px 24px;display:flex;align-items:center;gap:12px}
.topbar h1{font-size:18px;font-weight:600;color:#fff}
.topbar .ip-badge{margin-left:auto;background:#1e2235;border:1px solid #2d3148;border-radius:6px;padding:6px 14px;font-size:13px;color:#94a3b8}
.topbar .ip-badge span{color:#60a5fa;font-weight:600}
.container{max-width:860px;margin:32px auto;padding:0 20px}
.card{background:#1a1d27;border:1px solid #2d3148;border-radius:10px;padding:20px;margin-bottom:16px}
.card h2{font-size:14px;font-weight:600;color:#94a3b8;text-transform:uppercase;letter-spacing:.06em;margin-bottom:14px}
.btn{display:inline-flex;align-items:center;gap:6px;padding:7px 14px;border-radius:6px;border:none;cursor:pointer;font-size:13px;font-weight:500;transition:background .15s}
.btn-primary{background:#3b82f6;color:#fff}.btn-primary:hover{background:#2563eb}
.btn-danger{background:#ef4444;color:#fff}.btn-danger:hover{background:#dc2626}
.btn-ghost{background:#2d3148;color:#e2e8f0}.btn-ghost:hover{background:#374168}
.btn-sm{padding:4px 10px;font-size:12px}
.form-row{display:grid;gap:10px;margin-bottom:10px}
.form-row.cols2{grid-template-columns:1fr 1fr}
label{display:block;font-size:12px;color:#94a3b8;margin-bottom:4px}
input,select{width:100%;background:#0f1117;border:1px solid #2d3148;border-radius:6px;padding:8px 10px;color:#e2e8f0;font-size:13px;outline:none}
input:focus,select:focus{border-color:#3b82f6}
select option{background:#1a1d27}
.entry-list{display:flex;flex-direction:column;gap:10px}
.entry{background:#0f1117;border:1px solid #2d3148;border-radius:8px;padding:14px 16px;display:flex;align-items:center;gap:12px}
.entry-info{flex:1;min-width:0}
.entry-name{font-weight:600;font-size:14px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.entry-meta{font-size:12px;color:#64748b;margin-top:2px}
.entry-meta span{margin-right:10px}
.badge{display:inline-block;padding:2px 8px;border-radius:4px;font-size:11px;font-weight:600}
.badge-ok{background:#14532d;color:#4ade80}
.badge-error{background:#450a0a;color:#f87171}
.badge-pending{background:#1e293b;color:#94a3b8}
.badge-disabled{background:#1e293b;color:#475569}
.toggle{width:36px;height:20px;background:#374151;border-radius:10px;position:relative;cursor:pointer;border:none;flex-shrink:0;transition:background .2s}
.toggle.on{background:#3b82f6}
.toggle::after{content:'';position:absolute;width:14px;height:14px;background:#fff;border-radius:50%;top:3px;left:3px;transition:left .2s}
.toggle.on::after{left:19px}
.divider{border:none;border-top:1px solid #2d3148;margin:14px 0}
.settings-row{display:flex;align-items:center;gap:10px}
.settings-row label{margin:0;white-space:nowrap;color:#e2e8f0;font-size:13px}
.action-bar{display:flex;gap:8px;margin-bottom:16px}
.provider-params{display:flex;flex-direction:column;gap:10px;margin-top:10px}
#add-form{display:none}
.empty{text-align:center;padding:32px;color:#475569;font-size:14px}
</style>
</head>
<body>
<div class="topbar">
  <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#3b82f6" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="2" y1="12" x2="22" y2="12"/><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/></svg>
  <h1>AutoDNS</h1>
  <div class="ip-badge">Public IP: <span id="current-ip">—</span></div>
</div>

<div class="container">
  <div class="action-bar">
    <button class="btn btn-primary" onclick="toggleAddForm()">+ Add Entry</button>
    <button class="btn btn-ghost" onclick="forceUpdate()">↻ Update Now</button>
  </div>

  <div class="card" id="add-form">
    <h2>New DNS Entry</h2>
    <div class="form-row cols2">
      <div>
        <label>Name</label>
        <input id="f-name" type="text" placeholder="My Home">
      </div>
      <div>
        <label>Provider</label>
        <select id="f-provider" onchange="loadParams()">
          <option value="">Select provider…</option>
        </select>
      </div>
    </div>
    <div id="f-params" class="provider-params"></div>
    <hr class="divider">
    <div style="display:flex;gap:8px">
      <button class="btn btn-primary" onclick="submitEntry()">Save</button>
      <button class="btn btn-ghost" onclick="toggleAddForm()">Cancel</button>
    </div>
  </div>

  <div class="card">
    <h2>DNS Entries</h2>
    <div class="entry-list" id="entries"><div class="empty">No entries yet. Click "Add Entry" to get started.</div></div>
  </div>

  <div class="card">
    <h2>Settings</h2>
    <div class="form-row">
      <div>
        <label>Check interval (minutes)</label>
        <input id="s-interval" type="number" min="1" value="5">
      </div>
    </div>
    <div class="form-row">
      <div>
        <label>IP detector URL</label>
        <input id="s-detector" type="text" placeholder="https://api.ipify.org">
      </div>
    </div>
    <button class="btn btn-primary btn-sm" onclick="saveSettings()">Save Settings</button>
  </div>
</div>

<footer style="text-align:center;padding:24px 20px;color:#475569;font-size:12px;border-top:1px solid #2d3148;margin-top:8px;display:flex;justify-content:space-between;max-width:860px;margin:8px auto 0">
  <span>© 2026 AutoDNS. All rights reserved.</span>
  <span>{{VERSION}}</span>
</footer>

<script>
let providers = {};
let entries = [];

async function api(method, path, body) {
  const opts = {method, headers:{'Content-Type':'application/json'}};
  if (body !== undefined) opts.body = JSON.stringify(body);
  const r = await fetch(path, opts);
  return r.json();
}

async function load() {
  const [status, cfg, provList] = await Promise.all([
    api('GET', 'api/status'),
    api('GET', 'api/config'),
    api('GET', 'api/providers'),
  ]);

  document.getElementById('current-ip').textContent = status.current_ip || '—';
  document.getElementById('s-interval').value = cfg.check_interval_min;
  document.getElementById('s-detector').value = cfg.ip_detector_url;

  providers = {};
  const sel = document.getElementById('f-provider');
  sel.innerHTML = '<option value="">Select provider…</option>';
  provList.forEach(p => {
    providers[p.name] = p;
    const opt = document.createElement('option');
    opt.value = p.name;
    opt.textContent = p.label;
    sel.appendChild(opt);
  });

  entries = status.entries || [];
  renderEntries();
}

function renderEntries() {
  const el = document.getElementById('entries');
  if (!entries.length) {
    el.innerHTML = '<div class="empty">No entries yet. Click "Add Entry" to get started.</div>';
    return;
  }
  el.innerHTML = entries.map(e => {
    const badge = e.enabled
      ? (e.last_status === 'ok' ? '<span class="badge badge-ok">OK</span>'
        : e.last_status === 'error' ? '<span class="badge badge-error">Error</span>'
        : '<span class="badge badge-pending">Pending</span>')
      : '<span class="badge badge-disabled">Disabled</span>';

    const lastUpdate = e.last_update && e.last_update !== '0001-01-01T00:00:00Z'
      ? 'Updated ' + relTime(e.last_update) : 'Never updated';

    const errMsg = e.last_error ? '<span style="color:#f87171">' + esc(e.last_error) + '</span>' : '';

    return ` + "`" + `
    <div class="entry">
      <div class="entry-info">
        <div class="entry-name">${esc(e.name)}</div>
        <div class="entry-meta">
          <span>${provLabel(e.provider)}</span>
          <span>${badge}</span>
          <span>${lastUpdate}</span>
          ${errMsg}
        </div>
        <div class="entry-meta" style="margin-top:4px;color:#60a5fa;font-size:11px">${e.last_ip || ''}</div>
      </div>
      <button class="toggle ${e.enabled ? 'on' : ''}" onclick="toggle('${e.id}')" title="${e.enabled ? 'Disable' : 'Enable'}"></button>
      <button class="btn btn-danger btn-sm" onclick="deleteEntry('${e.id}')">Delete</button>
    </div>` + "`" + `;
  }).join('');
}

function provLabel(name) {
  return (providers[name] && providers[name].label) || name;
}

function esc(s) {
  return String(s || '').replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;');
}

function relTime(iso) {
  const d = new Date(iso), now = Date.now(), diff = Math.floor((now - d) / 1000);
  if (diff < 60) return diff + 's ago';
  if (diff < 3600) return Math.floor(diff/60) + 'm ago';
  if (diff < 86400) return Math.floor(diff/3600) + 'h ago';
  return Math.floor(diff/86400) + 'd ago';
}

function loadParams() {
  const name = document.getElementById('f-provider').value;
  const container = document.getElementById('f-params');
  if (!name || !providers[name]) { container.innerHTML = ''; return; }
  container.innerHTML = providers[name].param_defs.map(p => ` + "`" + `
    <div>
      <label>${esc(p.label)}</label>
      <input type="${p.secret ? 'password' : 'text'}" id="p-${p.key}" placeholder="${esc(p.placeholder)}">
    </div>` + "`" + `).join('');
}

function toggleAddForm() {
  const f = document.getElementById('add-form');
  f.style.display = f.style.display === 'block' ? 'none' : 'block';
}

async function submitEntry() {
  const name = document.getElementById('f-name').value.trim();
  const provider = document.getElementById('f-provider').value;
  if (!name || !provider) { alert('Name and provider are required'); return; }

  const params = {};
  if (providers[provider]) {
    providers[provider].param_defs.forEach(p => {
      params[p.key] = document.getElementById('p-' + p.key).value.trim();
    });
  }

  const res = await api('POST', 'api/entries', {name, provider, params});
  if (res.error) { alert(res.error); return; }

  document.getElementById('f-name').value = '';
  document.getElementById('f-provider').value = '';
  document.getElementById('f-params').innerHTML = '';
  toggleAddForm();
  await load();
}

async function deleteEntry(id) {
  if (!confirm('Delete this entry?')) return;
  await api('DELETE', 'api/entries?id=' + id);
  await load();
}

async function toggle(id) {
  await api('POST', 'api/entries/toggle?id=' + id);
  await load();
}

async function forceUpdate() {
  await api('POST', 'api/update');
  setTimeout(load, 2000);
}

async function saveSettings() {
  const interval = parseInt(document.getElementById('s-interval').value);
  const url = document.getElementById('s-detector').value.trim();
  if (interval < 1) { alert('Interval must be at least 1 minute'); return; }
  await api('POST', 'api/config', {check_interval_min: interval, ip_detector_url: url});
  alert('Settings saved');
}

load();
setInterval(load, 30000);
</script>
</body>
</html>`
