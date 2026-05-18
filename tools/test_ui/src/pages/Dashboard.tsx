import { useState } from 'react'
import { request } from '../api/client'
import { useConfig } from '../hooks/useConfig'
import { useAuth } from '../hooks/useAuth'

export default function Dashboard() {
  const { baseURL, setURL } = useConfig()
  const { token, user } = useAuth()
  const [status, setStatus] = useState<Record<string, string>>({})
  const [directStatus, setDirectStatus] = useState<Record<string, string>>({})

  const checkGateway = async () => {
    try {
      const d = await request<any>('GET', '/health')
      setStatus(s => ({ ...s, gateway: `OK (${d.service || 'api-gateway'})` }))
    } catch { setStatus(s => ({ ...s, gateway: 'FAIL' })) }
  }

  const checkDirect = async (name: string, url: string) => {
    try {
      const res = await fetch(url)
      const d = await res.json()
      setDirectStatus(s => ({ ...s, [name]: `OK (${d.service || name})` }))
    } catch { setDirectStatus(s => ({ ...s, [name]: 'FAIL' })) }
  }

  return (
    <div>
      <h2>Dashboard</h2>
      <div className="card">
        <h3>API Base URL</h3>
        <div className="form-group">
          <input value={baseURL} onChange={e => setURL(e.target.value)} />
        </div>
        <button onClick={checkGateway}>Check Gateway Health</button>
        {status.gateway && <p className={status.gateway.startsWith('OK') ? 'success' : 'error'}>{status.gateway}</p>}
      </div>

      <div className="card">
        <h3>Direct Service Health</h3>
        <div style={{ display: 'flex', gap: '0.5rem', flexWrap: 'wrap' }}>
          {[
            ['auth', 'http://localhost:8081/health'],
            ['run', 'http://localhost:8082/health'],
            ['submission', 'http://localhost:8083/health'],
            ['repo', 'http://localhost:8084/health'],
          ].map(([name, url]) => (
            <button key={name} onClick={() => checkDirect(name, url)}>{name}</button>
          ))}
        </div>
        {Object.entries(directStatus).map(([k, v]) => (
          <p key={k} className={v.startsWith('OK') ? 'success' : 'error'}>{k}: {v}</p>
        ))}
      </div>

      {user && (
        <div className="card">
          <h3>Logged In</h3>
          <pre>{JSON.stringify(user, null, 2)}</pre>
        </div>
      )}
    </div>
  )
}
