import {
  getEvents,
  getLiveEvents,
  getEventById,
  getCategories,
  getCategoryBySlug,
  getJourneys,
  getJourneyById,
  bookJourney,
  getMerchItems,
  getMerchItemById,
  login,
  register,
  getCurrentUser,
} from '@/lib/api'

// ─── Mock fetch globally ───────────────────────────────────────────────────

const mockFetch = jest.fn()
global.fetch = mockFetch

function mockOk(body: unknown) {
  mockFetch.mockResolvedValueOnce({
    ok: true,
    json: async () => body,
  })
}

function mockError(status: number, message: string) {
  mockFetch.mockResolvedValueOnce({
    ok: false,
    status,
    json: async () => ({ error: message }),
  })
}

beforeEach(() => {
  jest.clearAllMocks()
})

// ─── Events ───────────────────────────────────────────────────────────────

describe('getEvents', () => {
  it('returns events array on success', async () => {
    const fakeEvents = [{ id: '1', title: 'Daytona 500', isLive: true }]
    mockOk({ events: fakeEvents, count: 1 })

    const result = await getEvents()
    expect(result).toEqual(fakeEvents)
  })

  it('passes category query param when provided', async () => {
    mockOk({ events: [], count: 0 })
    await getEvents('motorsport')
    expect(mockFetch).toHaveBeenCalledWith(
      expect.stringContaining('?category=motorsport'),
      expect.any(Object)
    )
  })

  it('throws on API error', async () => {
    mockError(500, 'server error')
    await expect(getEvents()).rejects.toThrow('server error')
  })
})

describe('getLiveEvents', () => {
  it('returns only live events', async () => {
    const liveEvents = [{ id: '1', title: 'Live Race', isLive: true }]
    mockOk({ events: liveEvents, count: 1 })

    const result = await getLiveEvents()
    expect(result).toEqual(liveEvents)
    expect(mockFetch).toHaveBeenCalledWith(
      expect.stringContaining('/events/live'),
      expect.any(Object)
    )
  })

  it('throws on API error', async () => {
    mockError(500, 'server error')
    await expect(getLiveEvents()).rejects.toThrow()
  })
})

describe('getEventById', () => {
  it('returns a single event on success', async () => {
    const fakeEvent = { id: '42', title: 'Dakar Rally', isLive: false }
    mockOk(fakeEvent)

    const result = await getEventById('42')
    expect(result).toEqual(fakeEvent)
    expect(mockFetch).toHaveBeenCalledWith(
      expect.stringContaining('/events/42'),
      expect.any(Object)
    )
  })

  it('throws 404 on unknown id', async () => {
    mockError(404, 'event not found')
    await expect(getEventById('999')).rejects.toThrow('event not found')
  })
})

// ─── Categories ───────────────────────────────────────────────────────────

describe('getCategories', () => {
  it('returns categories array on success', async () => {
    const fakeCats = [{ id: '1', name: 'MOTORSPORT', slug: 'motorsport' }]
    mockOk({ categories: fakeCats, count: 1 })

    const result = await getCategories()
    expect(result).toEqual(fakeCats)
  })

  it('throws on API error', async () => {
    mockError(500, 'server error')
    await expect(getCategories()).rejects.toThrow()
  })
})

describe('getCategoryBySlug', () => {
  it('returns a category on success', async () => {
    const fakeCat = { id: '1', name: 'MOTORSPORT', slug: 'motorsport' }
    mockOk(fakeCat)

    const result = await getCategoryBySlug('motorsport')
    expect(result).toEqual(fakeCat)
    expect(mockFetch).toHaveBeenCalledWith(
      expect.stringContaining('/categories/motorsport'),
      expect.any(Object)
    )
  })

  it('throws 404 on unknown slug', async () => {
    mockError(404, 'category not found')
    await expect(getCategoryBySlug('unknown')).rejects.toThrow('category not found')
  })
})

// ─── Journeys ─────────────────────────────────────────────────────────────

describe('getJourneys', () => {
  it('returns journeys array on success', async () => {
    const fakeJourneys = [{ id: '1', title: 'Daytona Pit Experience', slotsLeft: 12 }]
    mockOk({ journeys: fakeJourneys, count: 1 })

    const result = await getJourneys()
    expect(result).toEqual(fakeJourneys)
  })

  it('throws on API error', async () => {
    mockError(500, 'server error')
    await expect(getJourneys()).rejects.toThrow()
  })
})

describe('getJourneyById', () => {
  it('returns a single journey on success', async () => {
    const fakeJourney = { id: '1', title: 'Daytona Pit Experience', slotsLeft: 12 }
    mockOk(fakeJourney)

    const result = await getJourneyById('1')
    expect(result).toEqual(fakeJourney)
  })

  it('throws 404 on unknown id', async () => {
    mockError(404, 'journey not found')
    await expect(getJourneyById('999')).rejects.toThrow('journey not found')
  })
})

describe('bookJourney', () => {
  it('sends POST with auth header and returns message + journey', async () => {
    const fakeResponse = { message: 'booked', journey: { id: '1', slotsLeft: 11 } }
    mockOk(fakeResponse)

    const result = await bookJourney('1', 'test-token')
    expect(result.message).toBe('booked')
    expect(mockFetch).toHaveBeenCalledWith(
      expect.stringContaining('/journeys/1/book'),
      expect.objectContaining({
        method: 'POST',
        headers: expect.objectContaining({ Authorization: 'Bearer test-token' }),
      })
    )
  })

  it('throws 401 when no valid token', async () => {
    mockError(401, 'unauthorized')
    await expect(bookJourney('1', 'bad-token')).rejects.toThrow('unauthorized')
  })
})

// ─── Merch ────────────────────────────────────────────────────────────────

describe('getMerchItems', () => {
  it('returns merch items array on success', async () => {
    const fakeItems = [{ id: '1', name: 'Team Hoodie', price: 89 }]
    mockOk({ items: fakeItems, count: 1 })

    const result = await getMerchItems()
    expect(result).toEqual(fakeItems)
  })

  it('throws on API error', async () => {
    mockError(500, 'server error')
    await expect(getMerchItems()).rejects.toThrow()
  })
})

describe('getMerchItemById', () => {
  it('returns a single merch item on success', async () => {
    const fakeItem = { id: '1', name: 'Team Hoodie', price: 89 }
    mockOk(fakeItem)

    const result = await getMerchItemById('1')
    expect(result).toEqual(fakeItem)
  })

  it('throws 404 on unknown id', async () => {
    mockError(404, 'item not found')
    await expect(getMerchItemById('999')).rejects.toThrow('item not found')
  })
})

// ─── Auth ─────────────────────────────────────────────────────────────────

describe('login', () => {
  it('returns token and user on valid credentials', async () => {
    const fakeResponse = { token: 'jwt-abc123', user: { id: '1', name: 'Alice', email: 'alice@test.com' } }
    mockOk(fakeResponse)

    const result = await login('alice@test.com', 'password123')
    expect(result.token).toBe('jwt-abc123')
    expect(result.user.name).toBe('Alice')
  })

  it('sends email and password in POST body', async () => {
    mockOk({ token: 'tok', user: {} })
    await login('alice@test.com', 'password123')

    const callArgs = mockFetch.mock.calls[0]
    const body = JSON.parse(callArgs[1].body)
    expect(body.email).toBe('alice@test.com')
    expect(body.password).toBe('password123')
  })

  it('throws on wrong password', async () => {
    mockError(401, 'invalid credentials')
    await expect(login('alice@test.com', 'wrongpass')).rejects.toThrow('invalid credentials')
  })

  it('throws on unknown email', async () => {
    mockError(401, 'invalid credentials')
    await expect(login('nobody@test.com', 'pass')).rejects.toThrow('invalid credentials')
  })
})

describe('register', () => {
  it('returns token and user on successful registration', async () => {
    const fakeResponse = { token: 'jwt-new', user: { id: '2', name: 'Bob', email: 'bob@test.com' } }
    mockOk(fakeResponse)

    const result = await register('bob@test.com', 'securepass', 'Bob')
    expect(result.token).toBe('jwt-new')
    expect(result.user.email).toBe('bob@test.com')
  })

  it('sends email, password, name in POST body', async () => {
    mockOk({ token: 'tok', user: {} })
    await register('bob@test.com', 'securepass', 'Bob')

    const callArgs = mockFetch.mock.calls[0]
    const body = JSON.parse(callArgs[1].body)
    expect(body.email).toBe('bob@test.com')
    expect(body.password).toBe('securepass')
    expect(body.name).toBe('Bob')
  })

  it('throws on duplicate email', async () => {
    mockError(409, 'email already exists')
    await expect(register('alice@test.com', 'pass', 'Alice2')).rejects.toThrow('email already exists')
  })
})

describe('getCurrentUser', () => {
  it('returns user data with valid token', async () => {
    const fakeUser = { id: '1', name: 'Alice', email: 'alice@test.com' }
    mockOk(fakeUser)

    const result = await getCurrentUser('valid-token')
    expect(result).toEqual(fakeUser)
    expect(mockFetch).toHaveBeenCalledWith(
      expect.stringContaining('/auth/me'),
      expect.objectContaining({
        headers: expect.objectContaining({ Authorization: 'Bearer valid-token' }),
      })
    )
  })

  it('throws 401 on invalid token', async () => {
    mockError(401, 'unauthorized')
    await expect(getCurrentUser('bad-token')).rejects.toThrow('unauthorized')
  })
})
