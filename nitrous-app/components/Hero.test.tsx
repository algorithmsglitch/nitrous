import { render, screen } from '@testing-library/react'
import Hero from '@/components/Hero'

// Mock next/image — jsdom can't handle Next.js image optimization
jest.mock('next/image', () => {
  return function MockImage({ src, alt }: { src: string; alt: string }) {
    return <img src={src} alt={alt} />
  }
})

// Mock next/link
jest.mock('next/link', () => {
  return function MockLink({ href, children, className }: any) {
    return <a href={href} className={className}>{children}</a>
  }
})

describe('Hero', () => {
  beforeEach(() => {
    render(<Hero />)
  })

  it('renders the hero section', () => {
    expect(document.querySelector('section')).toBeInTheDocument()
  })

  it('renders the main NITROUS heading', () => {
    expect(screen.getByRole('heading', { level: 1 })).toBeInTheDocument()
    expect(screen.getByText('NITROUS')).toBeInTheDocument()
  })

  it('renders tagline text', () => {
    expect(screen.getByText(/FUEL/i)).toBeInTheDocument()
  })

  it('renders hero subtitle copy', () => {
    expect(screen.getByText(/Stream every race/i)).toBeInTheDocument()
  })

  it('renders Ignite Stream CTA button', () => {
    expect(screen.getByRole('button', { name: /ignite stream/i })).toBeInTheDocument()
  })

  it('renders Explore Events button', () => {
    expect(screen.getByRole('button', { name: /explore events/i })).toBeInTheDocument()
  })

  it('renders the hero background image', () => {
    const img = screen.getByAltText(/nitrous wireframe car/i)
    expect(img).toBeInTheDocument()
  })

  it('renders all 6 nav rail cards', () => {
    const navLinks = screen.getAllByRole('link')
    // Filter to nav rail links (they have hrefs like /garage, /passes etc)
    const railLinks = navLinks.filter((l) =>
      ['/garage', '/passes', '/live', '/teams', '/journeys', '/merch'].includes(
        l.getAttribute('href') || ''
      )
    )
    expect(railLinks).toHaveLength(6)
  })

  it('renders JOURNEYS nav card', () => {
    expect(screen.getByRole('link', { name: /journeys/i })).toBeInTheDocument()
  })

  it('renders MERCH nav card', () => {
    expect(screen.getByRole('link', { name: /merch/i })).toBeInTheDocument()
  })

  it('renders LIVE STREAMS nav card', () => {
    expect(screen.getByText(/live streams/i)).toBeInTheDocument()
  })

  it('renders the HUD online label', () => {
    expect(screen.getByText(/system online/i)).toBeInTheDocument()
  })
})
