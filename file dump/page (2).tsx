'use client'
import { useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { login } from '@/lib/api'
import styles from './login.module.css'

export default function LoginPage() {
  const router = useRouter()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      const { token, user } = await login(email, password)
      localStorage.setItem('nitrous_token', token)
      localStorage.setItem('nitrous_user', JSON.stringify(user))
      router.push('/')
    } catch (err: any) {
      setError(err.message || 'Invalid credentials')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className={styles.page}>
      {/* Tron grid inherited from globals.css via body */}

      <Link href="/" className={styles.logo}>
        NITROUS<span>.</span>
      </Link>

      <div className={styles.card}>
        <div className={styles.cardHeader}>
          <div className={styles.hudRow}>
            <span className={styles.hudLine} />
            <span className={styles.hudDot} />
            <span className={styles.hudTxt}>IDENTITY VERIFICATION</span>
          </div>
          <h1 className={styles.title}>Sign In</h1>
          <p className={styles.sub}>Access your garage, streams &amp; passes</p>
        </div>

        <form
          onSubmit={handleSubmit}
          className={styles.form}
          data-testid="login-form"
        >
          <div className={styles.field}>
            <label htmlFor="email" className={styles.label}>
              Email
            </label>
            <input
              id="email"
              name="email"
              type="email"
              autoComplete="email"
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className={styles.input}
              placeholder="pilot@nitrous.io"
              data-testid="email-input"
            />
          </div>

          <div className={styles.field}>
            <label htmlFor="password" className={styles.label}>
              Password
            </label>
            <input
              id="password"
              name="password"
              type="password"
              autoComplete="current-password"
              required
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className={styles.input}
              placeholder="••••••••"
              data-testid="password-input"
            />
          </div>

          {error && (
            <div className={styles.errorMsg} data-testid="error-message" role="alert">
              <span className={styles.errorDot} />
              {error}
            </div>
          )}

          <button
            type="submit"
            className={styles.btnSubmit}
            disabled={loading}
            data-testid="submit-button"
          >
            {loading ? 'AUTHENTICATING...' : '▶  IGNITE ACCESS'}
          </button>
        </form>

        <div className={styles.footer}>
          <span className={styles.footerTxt}>No account?</span>
          <Link href="/register" className={styles.footerLink}>
            Register
          </Link>
        </div>
      </div>
    </div>
  )
}
