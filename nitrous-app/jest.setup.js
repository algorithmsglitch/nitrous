// Learn more: https://github.com/testing-library/jest-dom
import '@testing-library/jest-dom'

// Suppress console warnings that don't affect test results
const originalError = console.error
const originalWarn = console.warn

beforeAll(() => {
  console.error = (...args) => {
    // Suppress React/Next.js warnings that are non-critical
    const firstArg = args[0]?.toString?.() || ''
    
    // Suppress Next.js Image component attribute warnings
    if (firstArg.includes('non-boolean attribute')) {
      return
    }
    if (firstArg.includes('fill') || firstArg.includes('priority')) {
      return
    }
    
    return originalError.call(console, ...args)
  }
  
  console.warn = (...args) => {
    // Suppress similar warnings
    const firstArg = args[0]?.toString?.() || ''
    if (firstArg.includes('non-boolean attribute')) {
      return
    }
    return originalWarn.call(console, ...args)
  }
})

afterAll(() => {
  console.error = originalError
  console.warn = originalWarn
})
