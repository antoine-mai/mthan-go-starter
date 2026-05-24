import React, { useState, useEffect } from 'react';
import { 
  Server, 
  Terminal, 
  UserCheck, 
  CheckCircle, 
  AlertTriangle, 
  ArrowRight,
  RefreshCw,
  Send,
  Lock
} from 'lucide-react';

const BACKEND_URL = 'http://localhost:8080';

export default function App() {
  const [backendStatus, setBackendStatus] = useState<'checking' | 'online' | 'offline'>('checking');
  
  // State for hello api
  const [helloName, setHelloName] = useState('World');
  const [helloResponse, setHelloResponse] = useState<string | null>(null);
  const [helloLoading, setHelloLoading] = useState(false);
  const [helloError, setHelloError] = useState<string | null>(null);

  // State for post action
  const [actionName, setActionName] = useState('send_report');
  const [actionResponse, setActionResponse] = useState<string | null>(null);
  const [actionLoading, setActionLoading] = useState(false);
  const [actionError, setActionError] = useState<string | null>(null);

  // Checks the root API status
  const checkStatus = async () => {
    setBackendStatus('checking');
    try {
      const res = await fetch(`${BACKEND_URL}/api`);
      if (res.ok) {
        setBackendStatus('online');
      } else {
        setBackendStatus('offline');
      }
    } catch (err) {
      setBackendStatus('offline');
    }
  };

  useEffect(() => {
    checkStatus();
  }, []);

  const handleHelloSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setHelloLoading(true);
    setHelloError(null);
    setHelloResponse(null);
    try {
      const res = await fetch(`${BACKEND_URL}/api/hello?name=${encodeURIComponent(helloName)}`);
      const data = await res.json();
      if (data.success) {
        setHelloResponse(data.data.message);
      } else {
        setHelloError(data.error?.message || 'Error occurred');
      }
    } catch (err) {
      setHelloError('Failed to connect to Go server.');
    } finally {
      setHelloLoading(false);
    }
  };

  const handleActionSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setActionLoading(true);
    setActionError(null);
    setActionResponse(null);
    try {
      const res = await fetch(`${BACKEND_URL}/post/action`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ action: actionName })
      });
      const data = await res.json();
      if (data.success) {
        setActionResponse(data.data.message);
      } else {
        setActionError(data.error?.message || 'Error occurred');
      }
    } catch (err) {
      setActionError('Failed to connect to Go server.');
    } finally {
      setActionLoading(false);
    }
  };

  return (
    <div className="app-container">
      {/* Decorative Blur Orbs */}
      <div className="blur-orb orb-1"></div>
      <div className="blur-orb orb-2"></div>

      <header className="app-header">
        <div className="header-logo">
          <Server className="icon-main" />
          <div>
            <h1>MTHAN <span>App</span></h1>
            <p>Go Starter Kit Client Panel</p>
          </div>
        </div>

        <div className="status-badge-container">
          <div className={`status-badge ${backendStatus}`}>
            <span className="pulse-dot"></span>
            {backendStatus === 'checking' && 'Checking Server...'}
            {backendStatus === 'online' && 'Go Backend Online'}
            {backendStatus === 'offline' && 'Go Backend Offline'}
          </div>
          <button onClick={checkStatus} className="btn-refresh" title="Check connection">
            <RefreshCw size={16} />
          </button>
        </div>
      </header>

      <main className="dashboard-grid">
        {/* Card 1: API GET Greeting */}
        <section className="glass-card">
          <div className="card-header">
            <UserCheck className="card-icon blue" />
            <h2>Public Client API</h2>
          </div>
          <p className="card-desc">Tests the public GET greeting endpoint `/api/hello` to request service-level manipulation.</p>
          
          <form onSubmit={handleHelloSubmit} className="card-form">
            <div className="input-group">
              <label htmlFor="name-input">Query Parameter: name</label>
              <input 
                id="name-input"
                type="text" 
                value={helloName} 
                onChange={(e) => setHelloName(e.target.value)} 
                placeholder="Enter name..."
                required
              />
            </div>
            <button type="submit" className="btn-primary" disabled={helloLoading}>
              {helloLoading ? 'Requesting...' : 'Send Request'}
              <ArrowRight size={16} />
            </button>
          </form>

          {/* Response Box */}
          {(helloResponse || helloError) && (
            <div className={`response-box ${helloError ? 'error' : 'success'}`}>
              <div className="response-header">
                {helloError ? <AlertTriangle size={14} /> : <CheckCircle size={14} />}
                <span>{helloError ? 'API Error' : 'API Success'}</span>
              </div>
              <pre>{helloResponse || helloError}</pre>
            </div>
          )}
        </section>

        {/* Card 2: POST Action Simulation */}
        <section className="glass-card">
          <div className="card-header">
            <Terminal className="card-icon purple" />
            <h2>Built-in Post Action</h2>
          </div>
          <p className="card-desc">Tests the internal POST endpoint `/post/action` to execute mock database transactions or reports.</p>
          
          <form onSubmit={handleActionSubmit} className="card-form">
            <div className="input-group">
              <label htmlFor="action-input">JSON Payload: action</label>
              <select 
                id="action-input"
                value={actionName} 
                onChange={(e) => setActionName(e.target.value)}
              >
                <option value="send_report">send_report</option>
                <option value="sync_db">sync_db</option>
                <option value="clear_cache">clear_cache</option>
                <option value="trigger_error">trigger_error (simulates err)</option>
              </select>
            </div>
            <button type="submit" className="btn-primary purple" disabled={actionLoading}>
              {actionLoading ? 'Executing...' : 'Post Payload'}
              <Send size={16} />
            </button>
          </form>

          {/* Response Box */}
          {(actionResponse || actionError) && (
            <div className={`response-box ${actionError ? 'error' : 'success'}`}>
              <div className="response-header">
                {actionError ? <AlertTriangle size={14} /> : <CheckCircle size={14} />}
                <span>{actionError ? 'Execution Error' : 'Execution Success'}</span>
              </div>
              <pre>{actionResponse || actionError}</pre>
            </div>
          )}
        </section>

        {/* Card 3: Admin Panel Portal */}
        <section className="glass-card full-width">
          <div className="card-header">
            <Lock className="card-icon gold" />
            <h2>Administration Panel Gateway</h2>
          </div>
          <div className="admin-portal-content">
            <div className="admin-info">
              <p>The backend conditionally enables an admin control center. Make sure to set configuration credentials in your active environment configuration file (`.env`).</p>
              <div className="admin-credentials">
                <div className="cred-item">
                  <strong>Expected Path:</strong> <code>/admin</code> (or custom ADMIN_PATH)
                </div>
                <div className="cred-item">
                  <strong>Auth Method:</strong> HTTP Basic Authentication
                </div>
              </div>
            </div>
            <a 
              href={`${BACKEND_URL}/admin`} 
              target="_blank" 
              rel="noopener noreferrer" 
              className="btn-admin"
            >
              Open Admin Panel
              <ArrowRight size={18} />
            </a>
          </div>
        </section>
      </main>

      <footer className="app-footer">
        <p>Vietnamese Golang Starter Kit Design. Crafted with Vite, React & Go.</p>
      </footer>
    </div>
  );
}
