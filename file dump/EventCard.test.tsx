import { render, screen } from '@testing-library/react'
import EventCard from '@/components/EventCard'
import type { Event } from '@/types'

const baseEvent: Event = {
  id: '1',
  title: 'NASCAR Daytona 500',
  location: 'Daytona International Speedway · Florida',
  date: 'Feb 16, 2026',
  isLive: false,
  category: 'motorsport',
}

describe('EventCard', () => {
  it('renders the event title', () => {
    render(<EventCard event={baseEvent} />)
    expect(screen.getByText('NASCAR Daytona 500')).toBeInTheDocument()
  })

  it('renders the event location', () => {
    render(<EventCard event={baseEvent} />)
    expect(screen.getByText('Daytona International Speedway · Florida')).toBeInTheDocument()
  })

  it('renders the event date', () => {
    render(<EventCard event={baseEvent} />)
    expect(screen.getByText('Feb 16, 2026')).toBeInTheDocument()
  })

  it('renders the event category', () => {
    render(<EventCard event={baseEvent} />)
    expect(screen.getByText('motorsport')).toBeInTheDocument()
  })

  it('does NOT show LIVE badge when isLive is false', () => {
    render(<EventCard event={baseEvent} />)
    expect(screen.queryByTestId('live-badge')).not.toBeInTheDocument()
  })

  it('shows LIVE badge when isLive is true', () => {
    render(<EventCard event={{ ...baseEvent, isLive: true }} />)
    expect(screen.getByTestId('live-badge')).toBeInTheDocument()
  })

  it('LIVE badge contains the text LIVE', () => {
    render(<EventCard event={{ ...baseEvent, isLive: true }} />)
    expect(screen.getByText(/^LIVE$/)).toBeInTheDocument()
  })

  it('renders time when provided', () => {
    render(<EventCard event={{ ...baseEvent, time: '14:00 UTC' }} />)
    expect(screen.getByText('14:00 UTC')).toBeInTheDocument()
  })

  it('does not render time when not provided', () => {
    render(<EventCard event={baseEvent} />)
    // no time field on baseEvent, no time element should appear
    expect(screen.queryByText(/UTC/)).not.toBeInTheDocument()
  })

  it('renders the card container', () => {
    render(<EventCard event={baseEvent} />)
    expect(screen.getByTestId('event-card')).toBeInTheDocument()
  })
})
