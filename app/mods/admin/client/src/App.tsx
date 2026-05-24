import { useState, useEffect } from 'react';
import { 
  Database, 
  Activity, 
  Terminal, 
  Settings, 
  ShieldCheck, 
  Server, 
  Code2, 
  Cpu, 
  HardDrive, 
  Layers, 
  Play,
  RotateCw
} from 'lucide-react';

interface RouteInfo {
  path: string;
  methods: string[];
  description: string;
}

interface LogEntry {
  timestamp: string;
  level: 'INFO' | 'WARN' | 'ERROR';
  message: string;
}

export default function App() {
  const [activeTab, setActiveTab] = useState<'overview' | 'routes' | 'database' | 'logs'>('overview');
  const [dbStatus, setDbStatus] = useState<'connected' | 'disconnected'>('connected');
  const [cpuUsage, setCpuUsage] = useState<number>(14);
  const [memoryUsage, setMemoryUsage] = useState<number>(38);
  const [telemetryLogs, setTelemetryLogs] = useState<LogEntry[]>([]);
  const [isDiagnosticRunning, setIsDiagnosticRunning] = useState(false);

  // Initial mock routes
  const routes: RouteInfo[] = [
    { path: '/api', methods: ['GET', 'POST'], description: 'Base API endpoint providing gateway status.' },
    { path: '/api/hello', methods: ['GET', 'POST'], description: 'Sample endpoint returning dynamic hello greetings.' },
    { path: '/post/action', methods: ['POST'], description: 'Core endpoint processing administrative & post events.' },
    { path: '/admin', methods: ['GET'], description: 'Administrative control portal (Basic Auth Protected).' }
  ];

  // Load initial logs
  useEffect(() => {
    const initialLogs: LogEntry[] = [
      { timestamp: new Date(Date.now() - 300000).toLocaleTimeString(), level: 'INFO', message: 'LoggerService initialized with development context' },
      { timestamp: new Date(Date.now() - 280000).toLocaleTimeString(), level: 'INFO', message: 'SQLite database connection pool opened' },
      { timestamp: new Date(Date.now() - 250000).toLocaleTimeString(), level: 'INFO', message: 'Conditional client file server mapping resolved' },
      { timestamp: new Date(Date.now() - 200000).toLocaleTimeString(), level: 'INFO', message: 'CORS Preflight rules applied to /api/* and /post/*' },
      { timestamp: new Date(Date.now() - 150000).toLocaleTimeString(), level: 'WARN', message: 'High CPU load detected during client asset compilation' }
    ];
    setTelemetryLogs(initialLogs);
  }, []);

  // System metrics ticker
  useEffect(() => {
    const interval = setInterval(() => {
      setCpuUsage(prev => {
        const change = Math.floor(Math.random() * 9) - 4; // -4% to +4%
        const newVal = prev + change;
        return Math.max(5, Math.min(95, newVal));
      });
      setMemoryUsage(prev => {
        const change = Math.floor(Math.random() * 3) - 1; // -1% to +1%
        const newVal = prev + change;
        return Math.max(30, Math.min(60, newVal));
      });
    }, 3000);

    return () => clearInterval(interval);
  }, []);

  // Trigger diagnostic test
  const runDiagnostics = () => {
    setIsDiagnosticRunning(true);
    
    // Add starting log
    const newLog: LogEntry = {
      timestamp: new Date().toLocaleTimeString(),
      level: 'INFO',
      message: 'Starting full system diagnostics check...'
    };
    setTelemetryLogs(prev => [newLog, ...prev]);

    setTimeout(() => {
      setIsDiagnosticRunning(false);
      setDbStatus('connected');
      
      const completeLog: LogEntry = {
        timestamp: new Date().toLocaleTimeString(),
        level: 'INFO',
        message: 'Diagnostics completed. SQLite pool verified. Memory leaks: None.'
      };
      setTelemetryLogs(prev => [completeLog, ...prev]);
    }, 1500);
  };

  return (
    <div className="flex min-height-screen bg-slate-950 text-slate-100 overflow-hidden" style={{ minHeight: '100vh' }}>
      
      {/* Sidebar Navigation */}
      <aside className="w-64 border-r border-slate-800 bg-slate-900/50 backdrop-blur flex flex-col justify-between">
        <div>
          {/* Brand Header */}
          <div className="h-16 flex items-center gap-3 px-6 border-b border-slate-800">
            <div className="w-8 h-8 rounded bg-gradient-to-tr from-sky-500 to-indigo-600 flex items-center justify-center text-white font-bold text-lg shadow-lg shadow-sky-500/20">
              M
            </div>
            <div>
              <h1 className="font-semibold text-sm tracking-wide text-slate-200 uppercase">Mthan Portal</h1>
              <span className="text-[10px] text-sky-400 font-medium px-1.5 py-0.5 rounded bg-sky-950/50 border border-sky-900/50">Admin UI</span>
            </div>
          </div>

          {/* Navigation Links */}
          <nav className="p-4 space-y-1">
            <button 
              onClick={() => setActiveTab('overview')}
              className={`w-full flex items-center gap-3 px-4 py-2.5 rounded-lg text-sm font-medium transition-all ${activeTab === 'overview' ? 'bg-sky-500/10 text-sky-400 border border-sky-500/20' : 'text-slate-400 hover:bg-slate-800/50 hover:text-slate-200 border border-transparent'}`}
            >
              <Activity size={18} />
              Overview
            </button>
            <button 
              onClick={() => setActiveTab('routes')}
              className={`w-full flex items-center gap-3 px-4 py-2.5 rounded-lg text-sm font-medium transition-all ${activeTab === 'routes' ? 'bg-sky-500/10 text-sky-400 border border-sky-500/20' : 'text-slate-400 hover:bg-slate-800/50 hover:text-slate-200 border border-transparent'}`}
            >
              <Layers size={18} />
              Route Explorer
            </button>
            <button 
              onClick={() => setActiveTab('database')}
              className={`w-full flex items-center gap-3 px-4 py-2.5 rounded-lg text-sm font-medium transition-all ${activeTab === 'database' ? 'bg-sky-500/10 text-sky-400 border border-sky-500/20' : 'text-slate-400 hover:bg-slate-800/50 hover:text-slate-200 border border-transparent'}`}
            >
              <Database size={18} />
              Database Config
            </button>
            <button 
              onClick={() => setActiveTab('logs')}
              className={`w-full flex items-center gap-3 px-4 py-2.5 rounded-lg text-sm font-medium transition-all ${activeTab === 'logs' ? 'bg-sky-500/10 text-sky-400 border border-sky-500/20' : 'text-slate-400 hover:bg-slate-800/50 hover:text-slate-200 border border-transparent'}`}
            >
              <Terminal size={18} />
              Live Telemetry Logs
            </button>
          </nav>
        </div>

        {/* User context footer */}
        <div className="p-4 border-t border-slate-800 bg-slate-950/30 flex items-center gap-3">
          <div className="w-9 h-9 rounded-full bg-slate-800 flex items-center justify-center border border-slate-700">
            <ShieldCheck size={18} className="text-emerald-400" />
          </div>
          <div>
            <p className="text-xs font-semibold text-slate-300">Authorized Session</p>
            <p className="text-[10px] text-slate-500">Mode: Administrator</p>
          </div>
        </div>
      </aside>

      {/* Main Content Area */}
      <main className="flex-1 flex flex-col overflow-y-auto">
        {/* Top Header */}
        <header className="h-16 border-b border-slate-800 bg-slate-950/40 backdrop-blur flex items-center justify-between px-8">
          <h2 className="text-lg font-semibold tracking-tight text-slate-200 capitalize">
            {activeTab === 'overview' ? 'System Diagnostic Center' : activeTab === 'routes' ? 'Server API Directory' : activeTab === 'database' ? 'Database Engine Settings' : 'Real-time Console Stream'}
          </h2>
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2">
              <span className={`w-2.5 h-2.5 rounded-full ${dbStatus === 'connected' ? 'bg-emerald-500 shadow-lg shadow-emerald-500/30' : 'bg-red-500 shadow-lg shadow-red-500/30'}`}></span>
              <span className="text-xs font-medium text-slate-400">DB Connectivity</span>
            </div>
            <button 
              disabled={isDiagnosticRunning}
              onClick={runDiagnostics}
              className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-sky-600 hover:bg-sky-500 text-white text-xs font-medium transition-all shadow-md shadow-sky-500/10 active:scale-95 disabled:opacity-50"
            >
              <Play size={12} className={isDiagnosticRunning ? 'animate-spin' : ''} />
              Run Diagnostics
            </button>
          </div>
        </header>

        {/* Tab Contents */}
        <div className="p-8 flex-1 max-w-7xl w-full mx-auto space-y-6">
          
          {/* Tab 1: OVERVIEW */}
          {activeTab === 'overview' && (
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              
              {/* Telemetry Card: CPU */}
              <div className="p-6 rounded-xl border border-slate-800 bg-slate-900/30 backdrop-blur flex flex-col justify-between h-40">
                <div className="flex justify-between items-start">
                  <div>
                    <h3 className="text-xs font-semibold text-slate-500 uppercase tracking-wider">CPU Processor Core</h3>
                    <p className="text-3xl font-bold tracking-tight text-slate-100 mt-1">{cpuUsage}%</p>
                  </div>
                  <div className="p-2 rounded bg-sky-950/40 border border-sky-900/50">
                    <Cpu className="text-sky-400" size={20} />
                  </div>
                </div>
                <div className="w-full bg-slate-800 rounded-full h-1.5 mt-4">
                  <div className="bg-sky-500 h-1.5 rounded-full transition-all duration-1000" style={{ width: `${cpuUsage}%` }}></div>
                </div>
              </div>

              {/* Telemetry Card: MEMORY */}
              <div className="p-6 rounded-xl border border-slate-800 bg-slate-900/30 backdrop-blur flex flex-col justify-between h-40">
                <div className="flex justify-between items-start">
                  <div>
                    <h3 className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Resident RAM Allocation</h3>
                    <p className="text-3xl font-bold tracking-tight text-slate-100 mt-1">{memoryUsage}%</p>
                  </div>
                  <div className="p-2 rounded bg-indigo-950/40 border border-indigo-900/50">
                    <HardDrive className="text-indigo-400" size={20} />
                  </div>
                </div>
                <div className="w-full bg-slate-800 rounded-full h-1.5 mt-4">
                  <div className="bg-indigo-500 h-1.5 rounded-full transition-all duration-1000" style={{ width: `${memoryUsage}%` }}></div>
                </div>
              </div>

              {/* Telemetry Card: STATE */}
              <div className="p-6 rounded-xl border border-slate-800 bg-slate-900/30 backdrop-blur flex flex-col justify-between h-40">
                <div className="flex justify-between items-start">
                  <div>
                    <h3 className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Gateway Server Context</h3>
                    <p className="text-2xl font-bold tracking-tight text-emerald-400 mt-1.5 uppercase">Online</p>
                  </div>
                  <div className="p-2 rounded bg-emerald-950/40 border border-emerald-900/50">
                    <Server className="text-emerald-400" size={20} />
                  </div>
                </div>
                <div className="flex justify-between text-xs text-slate-500">
                  <span>Engine: Go HTTP Mux</span>
                  <span>Port: 8080</span>
                </div>
              </div>

              {/* System Details and Info */}
              <div className="md:col-span-2 p-6 rounded-xl border border-slate-800 bg-slate-900/20 space-y-4">
                <h3 className="font-semibold text-sm tracking-wider uppercase text-slate-400">Environment Metadata</h3>
                <div className="grid grid-cols-2 gap-4 text-sm">
                  <div className="p-4 rounded-lg bg-slate-950/40 border border-slate-900">
                    <span className="text-slate-500 text-xs">Admin Server Endpoint</span>
                    <p className="font-mono text-slate-300 mt-0.5">/admin</p>
                  </div>
                  <div className="p-4 rounded-lg bg-slate-950/40 border border-slate-900">
                    <span className="text-slate-500 text-xs">React UI Hosting</span>
                    <p className="font-mono text-slate-300 mt-0.5">Stand-alone Build</p>
                  </div>
                  <div className="p-4 rounded-lg bg-slate-950/40 border border-slate-900">
                    <span className="text-slate-500 text-xs">CORS Configuration</span>
                    <p className="font-mono text-slate-300 mt-0.5">Allow: localhost:3000/3001</p>
                  </div>
                  <div className="p-4 rounded-lg bg-slate-950/40 border border-slate-900">
                    <span className="text-slate-500 text-xs">Working Directory</span>
                    <p className="font-mono text-slate-300 mt-0.5 truncate">/app/mods/admin/client</p>
                  </div>
                </div>
              </div>

              {/* Panel Quick Actions */}
              <div className="p-6 rounded-xl border border-slate-800 bg-slate-900/20 flex flex-col justify-between">
                <div>
                  <h3 className="font-semibold text-sm tracking-wider uppercase text-slate-400 mb-3">Gateway Operations</h3>
                  <p className="text-xs text-slate-500">Quickly toggle DB settings or trigger logging traces directly from the browser window.</p>
                </div>
                <div className="space-y-2 mt-4">
                  <button 
                    onClick={() => setDbStatus(prev => prev === 'connected' ? 'disconnected' : 'connected')}
                    className="w-full py-2 rounded bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-semibold flex items-center justify-center gap-2 border border-slate-700 transition"
                  >
                    <RotateCw size={14} />
                    Toggle DB Pool Status
                  </button>
                  <button 
                    onClick={() => {
                      const triggerLog: LogEntry = {
                        timestamp: new Date().toLocaleTimeString(),
                        level: 'WARN',
                        message: 'Forced telemetry check initiated by administrator'
                      };
                      setTelemetryLogs(prev => [triggerLog, ...prev]);
                    }}
                    className="w-full py-2 rounded bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-semibold flex items-center justify-center gap-2 border border-slate-700 transition"
                  >
                    <Code2 size={14} />
                    Inject Warning Trace
                  </button>
                </div>
              </div>

            </div>
          )}

          {/* Tab 2: ROUTES */}
          {activeTab === 'routes' && (
            <div className="space-y-4">
              <div className="overflow-hidden border border-slate-800 rounded-xl bg-slate-900/20">
                <table className="w-full text-left border-collapse">
                  <thead>
                    <tr className="bg-slate-900/50 text-slate-400 text-xs font-semibold uppercase border-b border-slate-800">
                      <th className="px-6 py-4">URL Route Path</th>
                      <th className="px-6 py-4">HTTP Methods</th>
                      <th className="px-6 py-4">Gateway Scoping / Authorization</th>
                      <th className="px-6 py-4">Process Target</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-800 text-sm">
                    {routes.map((route, idx) => (
                      <tr key={idx} className="hover:bg-slate-900/10">
                        <td className="px-6 py-4 font-mono font-medium text-sky-400">{route.path}</td>
                        <td className="px-6 py-4">
                          <div className="flex gap-1.5">
                            {route.methods.map((method, mIdx) => (
                              <span key={mIdx} className={`px-2 py-0.5 rounded text-[10px] font-bold ${method === 'GET' ? 'bg-sky-950 text-sky-400 border border-sky-900/50' : 'bg-emerald-950 text-emerald-400 border border-emerald-900/50'}`}>
                                {method}
                              </span>
                            ))}
                          </div>
                        </td>
                        <td className="px-6 py-4 text-slate-400">{route.description}</td>
                        <td className="px-6 py-4 text-slate-500 font-mono">
                          {route.path === '/admin' ? 'app/mods/admin' : 'app/routes' + route.path}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          )}

          {/* Tab 3: DATABASE */}
          {activeTab === 'database' && (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              
              {/* Connection configuration */}
              <div className="p-6 rounded-xl border border-slate-800 bg-slate-900/20 space-y-4">
                <h3 className="font-semibold text-sm tracking-wider uppercase text-slate-400">Driver Mapping</h3>
                
                <div className="space-y-3">
                  <div>
                    <label className="text-xs text-slate-500 block mb-1">Active Driver</label>
                    <select className="w-full bg-slate-950 border border-slate-800 rounded-lg p-2 text-sm text-slate-300 font-mono outline-none">
                      <option value="sqlite">SQLite (sqlite3 / embed)</option>
                      <option value="postgres">PostgreSQL (pgx / tcp)</option>
                    </select>
                  </div>
                  <div>
                    <label className="text-xs text-slate-500 block mb-1">Connection DSN / URL string</label>
                    <input 
                      type="text" 
                      readOnly 
                      value="storage/db.sqlite" 
                      className="w-full bg-slate-950 border border-slate-800 rounded-lg p-2 text-sm text-slate-300 font-mono outline-none" 
                    />
                    <span className="text-[10px] text-slate-500 block mt-1">If blank, default SQLite path is auto-mapped under ~/.mthan-app/storage/db.sqlite.</span>
                  </div>
                </div>
              </div>

              {/* Status and telemetry */}
              <div className="p-6 rounded-xl border border-slate-800 bg-slate-900/20 space-y-4">
                <h3 className="font-semibold text-sm tracking-wider uppercase text-slate-400">Engine Statistics</h3>
                <div className="space-y-4 text-sm">
                  <div className="flex justify-between py-2 border-b border-slate-900">
                    <span className="text-slate-500">Database Driver</span>
                    <span className="font-mono text-slate-300 font-medium">sqlite</span>
                  </div>
                  <div className="flex justify-between py-2 border-b border-slate-900">
                    <span className="text-slate-500">Connection Pools Active</span>
                    <span className="font-mono text-slate-300 font-medium">1</span>
                  </div>
                  <div className="flex justify-between py-2 border-b border-slate-900">
                    <span className="text-slate-500">Max Open Connections</span>
                    <span className="font-mono text-slate-300 font-medium">100</span>
                  </div>
                  <div className="flex justify-between py-2">
                    <span className="text-slate-500">Database File Size</span>
                    <span className="font-mono text-slate-300 font-medium">0 KB (Empty)</span>
                  </div>
                </div>
              </div>

            </div>
          )}

          {/* Tab 4: TELEMETRY LOGS */}
          {activeTab === 'logs' && (
            <div className="space-y-4">
              <div className="p-4 bg-slate-950 border border-slate-800 rounded-xl flex flex-col h-[500px]">
                {/* Header controls */}
                <div className="flex justify-between items-center pb-3 border-b border-slate-900">
                  <div className="flex items-center gap-2">
                    <span className="w-3 h-3 rounded-full bg-red-500"></span>
                    <span className="w-3 h-3 rounded-full bg-yellow-500"></span>
                    <span className="w-3 h-3 rounded-full bg-green-500"></span>
                    <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider ml-2">Console Stream</span>
                  </div>
                  <button 
                    onClick={() => setTelemetryLogs([])}
                    className="px-2 py-1 rounded bg-slate-900 border border-slate-800 hover:bg-slate-800 text-[10px] text-slate-400 font-semibold transition"
                  >
                    Clear Terminal
                  </button>
                </div>

                {/* Console Log Panel */}
                <div className="flex-1 overflow-y-auto p-4 space-y-2.5 font-mono text-xs text-slate-300">
                  {telemetryLogs.length === 0 ? (
                    <div className="text-slate-600 italic py-8 text-center">No logs generated. Run a diagnostics check to populate events.</div>
                  ) : (
                    telemetryLogs.map((log, idx) => (
                      <div key={idx} className="flex items-start gap-4">
                        <span className="text-slate-600 select-none">{log.timestamp}</span>
                        <span className={`font-bold select-none px-1.5 py-0.5 rounded text-[9px] ${log.level === 'ERROR' ? 'bg-red-950 text-red-400 border border-red-900/50' : log.level === 'WARN' ? 'bg-yellow-950 text-yellow-400 border border-yellow-900/50' : 'bg-slate-900 text-slate-400 border border-slate-800'}`}>
                          {log.level}
                        </span>
                        <span className="flex-1 break-all">{log.message}</span>
                      </div>
                    ))
                  )}
                </div>
              </div>
            </div>
          )}

        </div>
      </main>

    </div>
  );
}
