import {
  getEvents,
  getCategories,
  getJourneys,
  getMerchItems,
  getEventById,
  login,
  register,
  getCurrentUser,
} from '@/lib/api'

// Mock the fetch function
globalThis.fetch = jest.fn()

// Suppress console.error for expected API errors in tests
const consoleErrorSpy = jest.spyOn(console, 'error').mockImplementation(() => {})

describe('API Utility Functions', () => {
  afterAll(() => {
    consoleErrorSpy.mockRestore()
  })

  beforeEach(() => {
    jest.clearAllMocks()
  })

  describe('getEvents', () => {
    it('fetches events successfully', async () => {
      const mockEvents = [
        { id: '1', name: 'Event 1', category: 'racing' },
        { id: '2', name: 'Event 2', category: 'motorsports' },
      ]
      ;(globalThis.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ events: mockEvents, count: 2 }),
      })

      const result = await getEvents()
      expect(result).toEqual(mockEvents)
      expect(globalThis.fetch).toHaveBeenCalled()
    })

    it('returns empty array on error', async () => {
      ;(globalThis.fetch as jest.Mock).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Not found' }),
      })

      await expect(getEvents()).rejects.toThrow()
    })
  })

  describe('getCategories', () => {
    it('fetches categories successfully', async () => {
      const mockCategories = [
        { id: '1', name: 'Racing', slug: 'racing' },
        { id: '2', name: 'Motorsports', slug: 'motorsports' },
      ]
      ;(globalThis.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ categories: mockCategories, count: 2 }),
      })

      const result = await getCategories()
      expect(result).toEqual(mockCategories)
    })
  })

  describe('getJourneys', () => {
    it('fetches journeys successfully', async () => {
      const mockJourneys = [
        { id: '1', name: 'Journey 1', description: 'A great journey' },
        { id: '2', name: 'Journey 2', description: 'Another journey' },
      ]
      ;(globalThis.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ journeys: mockJourneys, count: 2 }),
      })

      const result = await getJourneys()
      expect(result).toEqual(mockJourneys)
    })
  })

  describe('getMerchItems', () => {
    it('fetches merch items successfully', async () => {
      const mockMerch = [
        { id: '1', name: 'T-Shirt', price: 29.99 },
        { id: '2', name: 'Hat', price: 19.99 },
      ]
      ;(globalThis.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({ items: mockMerch, count: 2 }),
      })

      const result = await getMerchItems()
      expect(result).toEqual(mockMerch)
    })
  })

  describe('getEventById', () => {
    it('fetches event by ID successfully', async () => {
      const mockEvent = { id: '1', name: 'Event 1', category: 'racing' }
      ;(globalThis.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockEvent,
      })

      const result = await getEventById('1')
      expect(result).toEqual(mockEvent)
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/events/1'),
        expect.any(Object)
      )
    })
  })

  describe('Authentication Functions', () => {
    it('registers user successfully', async () => {
      const mockResponse = {
        user: { id: '1', email: 'test@example.com', name: 'Test User' },
        token: 'jwt_token_here',
      }
      ;(globalThis.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse,
      })

      const result = await register('test@example.com', 'password123', 'Test User')
      expect(result).toEqual(mockResponse)
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/auth/register'),
        expect.objectContaining({
          method: 'POST',
        })
      )
    })

    it('logs in user successfully', async () => {
      const mockResponse = {
        user: { id: '1', email: 'test@example.com' },
        token: 'jwt_token_here',
      }
      ;(globalThis.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse,
      })

      const result = await login('test@example.com', 'password123')
      expect(result).toEqual(mockResponse)
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/auth/login'),
        expect.objectContaining({
          method: 'POST',
        })
      )
    })

    it('gets current user with token', async () => {
      const mockUser = { id: '1', email: 'test@example.com', name: 'Test User' }
      ;(globalThis.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockUser,
      })

      const result = await getCurrentUser('test_token')
      expect(result).toEqual(mockUser)
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/auth/me'),
        expect.objectContaining({
          headers: expect.objectContaining({
            Authorization: 'Bearer test_token',
          }),
        })
      )
    })
  })

  describe('Error Handling', () => {
    it('handles API errors properly', async () => {
      ;(globalThis.fetch as jest.Mock).mockResolvedValueOnce({
        ok: false,
        json: async () => ({ error: 'Invalid request' }),
      })

      await expect(register('test@example.com', 'password', 'Test')).rejects.toThrow()
    })

    it('handles network errors', async () => {
      ;(globalThis.fetch as jest.Mock).mockRejectedValueOnce(
        new Error('Network error')
      )

      await expect(getEvents()).rejects.toThrow()
    })
  })
})
