import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'
import { ApiError } from '../api/client'

export default function Login() {
  const { login, token, user } = useAuth()
  const [email, setEmail] = useState('admin@example.com')
  const [password, setPassword] = useState('password')
  const [error, setError] = useState('')
  const navigate = useNavigate()

  const handleLogin = async () => {
    setError('')
    try { await login(email, password); navigate('/') }
    catch (e) { setError(e instanceof ApiError ? `Error ${e.status}: ${e.message}` : String(e)) }
  }

  if (token) return <div className="card"><h3>Logged In</h3><pre>{JSON.stringify(user, null, 2)}</pre></div>

  return (
    <div className="card" style={{ maxWidth: 400 }}>
      <h3>Login</h3>
      <div className="form-group"><label>Email</label><input value={email} onChange={e => setEmail(e.target.value)} /></div>
      <div className="form-group"><label>Password</label><input type="password" value={password} onChange={e => setPassword(e.target.value)} /></div>
      {error && <p className="error">{error}</p>}
      <button onClick={handleLogin}>Login</button>
    </div>
  )
}
