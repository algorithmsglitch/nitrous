import type { Event } from '@/types'
import styles from './EventCard.module.css'

interface EventCardProps {
  event: Event
}

export default function EventCard({ event }: EventCardProps) {
  return (
    <div className={styles.card} data-testid="event-card">
      {event.isLive && (
        <div className={styles.liveBadge} data-testid="live-badge">
          <span className={styles.liveDot} />
          LIVE
        </div>
      )}

      <div className={styles.body}>
        <div className={styles.category}>{event.category}</div>
        <h3 className={styles.title}>{event.title}</h3>
        <p className={styles.location}>{event.location}</p>

        <div className={styles.footer}>
          <span className={styles.date}>{event.date}</span>
          {event.time && <span className={styles.time}>{event.time}</span>}
        </div>
      </div>
    </div>
  )
}
