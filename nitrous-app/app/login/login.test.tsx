import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import LoginPage from '@/app/login/page'

// Mock next/navigation
const mockPush = jest.fn()
jest.mock('next/navigation', () => ({
  useRouter: () => ({ push: mockPush }),
}))

// Mock next/link
jest.mock('next/link', () => {
  return function MockLink({ href, children, className }: any) {
    return <a href={href} className={className}>{children}</a>
  }
})

// Mock the api module
jest.mock('@/lib/api', () => ({
  login: jest.fn(),
}))

import { login } from '@/lib/api'
const mockLogin = login as jest.Mock

// Mock localStorage
const localStorageMock = (() => {
  let store: Record<string, string> = {}
  return {
    getItem: (key: string) => store[key] || null,
    setItem: (key: string, value: string) => { store[key] = value },
    clear: () => { store = {} },
  }
})()
Object.defineProperty(window, 'localStorage', { value: localStorageMock })

describe('LoginPage', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    localStorageMock.clear()
  })

  it('renders the Sign In heading', () => {
    render(<LoginPage />)
    expect(screen.getByRole('heading', { name: /sign in/i })).toBeInTheDocument()
  })

  it('renders email and password inputs', () => {
    render(<LoginPage />)
    expect(screen.getByTestId('email-input')).toBeInTheDocument()
    expect(screen.getByTestId('password-input')).toBeInTheDocument()
  })

  it('renders the submit button', () => {
    render(<LoginPage />)
    expect(screen.getByTestId('submit-button')).toBeInTheDocument()
  })

  it('renders a link back to homepage (NITROUS logo)', () => {
    render(<LoginPage />)
    const logo = screen.getByRole('link', { name: /nitrous/i })
    expect(logo).toHaveAttribute('href', '/')
  })

  it('renders a Register link', () => {
    render(<LoginPage />)
    const registerLink = screen.getByRole('link', { name: /register/i })
    expect(registerLink).toHaveAttribute('href', '/register')
  })

  it('updates email field on input', async () => {
    render(<LoginPage />)
    const input = screen.getByTestId('email-input')
    await userEvent.type(input, 'test@test.com')
    expect(input).toHaveValue('test@test.com')
  })

  it('updates password field on input', async () => {
    render(<LoginPage />)
    const input = screen.getByTestId('password-input')
    await userEvent.type(input, 'password123')
    expect(input).toHaveValue('password123')
  })

  it('shows error message on failed login', async () => {
    mockLogin.mockRejectedValueOnce(new Error('Invalid credentials'))
    render(<LoginPage />)

    await userEvent.type(screen.getByTestId('email-input'), 'bad@test.com')
    await userEvent.type(screen.getByTestId('password-input'), 'wrongpass')
    fireEvent.submit(screen.getByTestId('login-form'))

    await waitFor(() => {
      expect(screen.getByTestId('error-message')).toBeInTheDocument()
      expect(screen.getByText(/invalid credentials/i)).toBeInTheDocument()
    })
  })

  it('does not show error message before submit', () => {
    render(<LoginPage />)
    expect(screen.queryByTestId('error-message')).not.toBeInTheDocument()
  })

  it('redirects to homepage on successful login', async () => {
    mockLogin.mockResolvedValueOnce({
      token: 'fake-jwt',
      user: { id: '1', name: 'Alice', email: 'alice@test.com' },
    })
    render(<LoginPage />)

    await userEvent.type(screen.getByTestId('email-input'), 'alice@test.com')
    await userEvent.type(screen.getByTestId('password-input'), 'password123')
    fireEvent.submit(screen.getByTestId('login-form'))

    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledWith('/')
    })
  })

  it('saves token to localStorage on successful login', async () => {
    mockLogin.mockResolvedValueOnce({
      token: 'fake-jwt',
      user: { id: '1', name: 'Alice' },
    })
    render(<LoginPage />)

    await userEvent.type(screen.getByTestId('email-input'), 'alice@test.com')
    await userEvent.type(screen.getByTestId('password-input'), 'password123')
    fireEvent.submit(screen.getByTestId('login-form'))

    await waitFor(() => {
      expect(localStorageMock.getItem('nitrous_token')).toBe('fake-jwt')
    })
  })

  it('disables submit button while loading', async () => {
    // login never resolves during this test
    mockLogin.mockImplementationOnce(() => new Promise(() => {}))
    render(<LoginPage />)

    await userEvent.type(screen.getByTestId('email-input'), 'alice@test.com')
    await userEvent.type(screen.getByTestId('password-input'), 'password123')
    fireEvent.submit(screen.getByTestId('login-form'))

    await waitFor(() => {
      expect(screen.getByTestId('submit-button')).toBeDisabled()
    })
  })
})
