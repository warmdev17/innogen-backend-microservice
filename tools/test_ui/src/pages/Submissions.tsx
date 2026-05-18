import { useState } from 'react'
import { request } from '../api/client'
import { usePolling } from '../hooks/usePolling'
import CodeEditor from '../components/CodeEditor'

export default function Submissions() {
  const [problemId, setProblemId] = useState('1')
  const [languageId, setLanguageId] = useState('1')
  const [code, setCode] = useState("const fs = require('fs');\nconst input = fs.readFileSync(0, 'utf8').trim().split(' ').map(Number);\nconsole.log(input[0] + input[1]);")
  const [submission, setSubmission] = useState<any>(null)
  const [error, setError] = useState('')
  const [pollUrl, setPollUrl] = useState<string | null>(null)
  const [list, setList] = useState<any[]>([])

  const submit = async () => {
    setError('')
    try {
      const d = await request<any>('POST', '/submit', { problemId: +problemId, languageId: +languageId, code })
      setSubmission(d.submission || d)
      setPollUrl(`/submissions/${d.submission.id}`)
    } catch (e: any) { setError(e.message) }
  }

  const pollResult = usePolling<any>(pollUrl, 2000, !!pollUrl)
  const finalStatus = pollResult.data?.submission?.status || pollResult.data?.status

  if (finalStatus && !['Pending', 'Running'].includes(finalStatus)) {
    if (pollUrl) setTimeout(() => setPollUrl(null), 0)
  }

  const loadList = async () => {
    try { const d = await request<any>('GET', '/me/submissions'); setList(d.submissions || []) }
    catch (e: any) { setError(e.message) }
  }

  return (
    <div>
      <h2>Submissions</h2>
      <div className="form-group"><label>Problem ID</label><input value={problemId} onChange={e => setProblemId(e.target.value)} /></div>
      <div className="form-group"><label>Language ID</label><input value={languageId} onChange={e => setLanguageId(e.target.value)} /></div>
      <CodeEditor value={code} onChange={setCode} label="Code" />
      <button onClick={submit}>Submit</button>
      <button onClick={loadList} style={{ marginLeft: '0.5rem' }}>List My Submissions</button>
      {error && <p className="error">{error}</p>}
      {submission && <pre>{JSON.stringify(submission, null, 2)}</pre>}
      {pollResult.data && (
        <div className="card">
          <h4>Poll Result</h4>
          <span className={`status-badge status-${finalStatus?.replace(/ /g, '')}`}>{finalStatus}</span>
          <pre>{JSON.stringify(pollResult.data, null, 2)}</pre>
        </div>
      )}
      {list.length > 0 && (
        <div className="card">
          <h4>My Submissions ({list.length})</h4>
          <pre>{JSON.stringify(list.slice(0, 10), null, 2)}</pre>
        </div>
      )}
    </div>
  )
}
