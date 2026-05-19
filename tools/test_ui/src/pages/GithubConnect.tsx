import { useState, useEffect, useCallback } from 'react'
import { useSearchParams, useNavigate } from 'react-router-dom'
import { request } from '../api/client'
import { useAuth } from '../hooks/useAuth'

interface ConnectionState {
  connected: boolean
  installationId?: string
  githubOwner?: string
  githubOwnerType?: string
  status?: string
}

export default function GithubConnect() {
  const { token } = useAuth()
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()

  const [state, setState] = useState<ConnectionState | null>(null)
  const [loading, setLoading] = useState(true)
  const [linking, setLinking] = useState(false)
  const [error, setError] = useState('')

  const checkConnection = async () => {
    try {
      setLoading(true)
      const d = await request<ConnectionState>('GET', '/github/connection')
      setState(d)
      setError('')
    } catch (e: any) {
      setError('Failed to check connection: ' + e.message)
    } finally {
      setLoading(false)
    }
  }

  const handleLink = useCallback(async (installationId: string) => {
    setLoading(false)
    setLinking(true)
    setError('')
    try {
      await request('POST', '/github/installations/link', { installationId })
      await checkConnection()
    } catch (e: any) {
      setError(e.message || 'Failed to link installation')
      await checkConnection()
    } finally {
      setLinking(false)
    }
  }, [])

  // On mount: check connection + handle callback
  useEffect(() => {
    const installationId = searchParams.get('installation_id')
    if (installationId) {
      handleLink(installationId)
    } else {
      checkConnection()
    }
  }, [searchParams, handleLink])

  if (!token) {
    return <div className="card"><h3>GitHub Connect</h3><p>Please log in first.</p></div>
  }

  if (loading) {
    return <div className="card"><h3>GitHub Connect</h3><p>Loading...</p></div>
  }

  if (linking) {
    return <div className="card"><h3>GitHub Connect</h3><p>Linking installation...</p></div>
  }

  return (
    <div>
      <h2>GitHub Connect</h2>

      {error && <div className="card" style={{ border: '1px solid #e53e3e', background: '#fff5f5' }}>
        <p className="error">{error}</p>
        <button onClick={checkConnection} style={{ marginTop: '0.5rem' }}>Retry</button>
      </div>}

      {state?.connected ? (
        <div className="card">
          <h3>✅ Connected</h3>
          <div style={{ marginTop: '0.5rem' }}>
            <p><strong>Owner:</strong> {state.githubOwner}</p>
            <p><strong>Type:</strong> {state.githubOwnerType}</p>
            <p><strong>Installation ID:</strong> {state.installationId}</p>
            <p><strong>Status:</strong> <span className="success">{state.status}</span></p>
          </div>
          <button onClick={checkConnection} style={{ marginTop: '1rem' }}>Refresh Status</button>
        </div>
      ) : (
        <div className="card">
          <h3>Not Connected</h3>
          <p style={{ marginBottom: '1rem' }}>Install the RinnoGen GitHub App to enable repository commits.</p>
          <button
            onClick={() => {
              const url = (import.meta as any).env?.VITE_GITHUB_APP_INSTALL_URL || 'https://github.com/apps/rinnogen/installations/new'
              window.location.href = url
            }}
          >
            Connect GitHub App
          </button>
        </div>
      )}

      <div className="info" style={{ marginTop: '1rem' }}>
        <strong>Note:</strong> GitHub App connection is stored in the backend database.
        If you just installed the app and it still shows not connected, GitHub may need a moment to deliver the webhook.
        Click <b>Refresh Status</b> or revisit this page.
      </div>

      <div className="info" style={{ marginTop: '1rem' }}>
        <strong>What happens after connecting?</strong>
        <p style={{ marginTop: '0.25rem' }}>
          When you submit a solution that passes all tests (<b>Accepted</b>), 
          your code is automatically committed to a GitHub repository under your 
          connected account. Each subject has its own repository named{' '}
          <code>&lt;subject&gt;-RinnoGen</code>. You can view commits on GitHub 
          and track your learning progress.
        </p>
      </div>
    </div>
  )
}
