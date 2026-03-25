import { render, screen } from '@testing-library/react'
import Nav from '@/components/Nav'

// Next/link renders as <a> in test environment
jest.mock('next/link', () => {
  return function MockLink({ href, children, className }: any) {
    return <a href={href} className={className}>{children}</a>
  }
})

describe('Nav', () => {
  beforeEach(() => {
    render(<Nav />)
  })

  it('renders the NITROUS logo', () => {
    expect(screen.getByText(/NITROUS/i)).toBeInTheDocument()
  })

  it('logo links to homepage', () => {
    const logo = screen.getByRole('link', { name: /NITROUS/i })
    expect(logo).toHaveAttribute('href', '/')
  })

  it('renders Live nav link', () => {
    expect(screen.getByRole('link', { name: /^live$/i })).toBeInTheDocument()
  })

  it('renders Events nav link', () => {
    expect(screen.getByRole('link', { name: /^events$/i })).toBeInTheDocument()
  })

  it('renders Teams nav link', () => {
    expect(screen.getByRole('link', { name: /^teams$/i })).toBeInTheDocument()
  })

  it('renders Journeys nav link', () => {
    expect(screen.getByRole('link', { name: /^journeys$/i })).toBeInTheDocument()
  })

  it('renders Merch nav link', () => {
    expect(screen.getByRole('link', { name: /^merch$/i })).toBeInTheDocument()
  })

  it('Events link points to /events', () => {
    const link = screen.getByRole('link', { name: /^events$/i })
    expect(link).toHaveAttribute('href', '/events')
  })

  it('Teams link points to /teams', () => {
    const link = screen.getByRole('link', { name: /^teams$/i })
    expect(link).toHaveAttribute('href', '/teams')
  })

  it('renders Sign In button', () => {
    expect(screen.getByRole('button', { name: /sign in/i })).toBeInTheDocument()
  })

  it('renders live status indicator text', () => {
    expect(screen.getByText(/events live/i)).toBeInTheDocument()
  })
})
