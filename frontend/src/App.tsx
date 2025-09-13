import { useState } from 'react'
import { apiService, type ApiResponse } from './api'

const API_ENDPOINT = 'https://ujxtgteuma.execute-api.ap-northeast-1.amazonaws.com/prod'

function App() {
  const [response, setResponse] = useState<ApiResponse | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [lastMethod, setLastMethod] = useState<string | null>(null)

  const handleGetRequest = async () => {
    setLoading(true)
    setError(null)
    setLastMethod('GET')
    try {
      const data = await apiService.get()
      setResponse(data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'GET request failed')
    } finally {
      setLoading(false)
    }
  }

  const handlePostRequest = async () => {
    setLoading(true)
    setError(null)
    setLastMethod('POST')
    try {
      const data = await apiService.post({ name: '太郎', message: 'こんにちは' })
      setResponse(data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'POST request failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-700 via-gray-800 to-slate-900 flex flex-col justify-center items-center p-4">
      <div className="w-full max-w-4xl text-center">
        <h1 className="text-5xl font-bold text-white mb-8 drop-shadow-lg">
          Echo API Test
        </h1>

        <div className="bg-white/90 backdrop-blur-lg rounded-2xl p-6 mb-8 shadow-xl hover:shadow-2xl transition-all duration-300 hover:-translate-y-1">
          <h3 className="text-lg font-semibold text-gray-800 mb-4">API Endpoint:</h3>
          <code className="bg-gray-100 border border-gray-300 rounded-xl px-4 py-3 font-mono text-sm text-gray-700 inline-block break-all shadow-inner">
            {API_ENDPOINT}
          </code>
        </div>

        <div className="flex flex-col sm:flex-row gap-6 justify-center mb-8">
          <button
            onClick={handleGetRequest}
            disabled={loading}
            className="relative px-8 py-4 bg-gradient-to-r from-slate-600 to-slate-700 text-white font-semibold rounded-xl shadow-lg hover:shadow-xl transform hover:-translate-y-1 transition-all duration-300 disabled:opacity-60 disabled:cursor-not-allowed disabled:transform-none overflow-hidden group min-w-[160px]"
          >
            <div className="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent -translate-x-full group-hover:translate-x-full transition-transform duration-500"></div>
            <span className="relative">{loading ? 'Loading...' : 'GET Request'}</span>
          </button>

          <button
            onClick={handlePostRequest}
            disabled={loading}
            className="relative px-8 py-4 bg-gradient-to-r from-gray-600 to-gray-700 text-white font-semibold rounded-xl shadow-lg hover:shadow-xl transform hover:-translate-y-1 transition-all duration-300 disabled:opacity-60 disabled:cursor-not-allowed disabled:transform-none overflow-hidden group min-w-[160px]"
          >
            <div className="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent -translate-x-full group-hover:translate-x-full transition-transform duration-500"></div>
            <span className="relative">{loading ? 'Loading...' : 'POST Request'}</span>
          </button>
        </div>

        {error && (
          <div className="bg-red-50/90 backdrop-blur-lg border border-red-200 rounded-2xl p-6 mb-8 text-red-700 shadow-xl animate-pulse">
            <h3 className="font-semibold mb-2">Error:</h3>
            <p>{error}</p>
          </div>
        )}

        {response && (
          <div className="bg-white/90 backdrop-blur-lg rounded-2xl p-6 text-left shadow-xl animate-fadeIn">
            <div className="flex items-center gap-3 mb-4">
              <h3 className="font-semibold text-gray-800">API Response:</h3>
              {lastMethod && (
                <span className={`px-3 py-1 rounded-lg text-xs font-bold tracking-wider ${
                  lastMethod === 'GET'
                    ? 'bg-green-100 text-green-800 border border-green-200'
                    : 'bg-blue-100 text-blue-800 border border-blue-200'
                }`}>
                  {lastMethod}
                </span>
              )}
            </div>
            <pre className="bg-gray-50 border border-gray-200 rounded-xl p-4 overflow-x-auto font-mono text-sm text-gray-700 leading-relaxed shadow-inner">
              {JSON.stringify(response, null, 2)}
            </pre>
          </div>
        )}
      </div>
    </div>
  )
}

export default App
