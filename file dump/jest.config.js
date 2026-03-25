const nextJest = require('next/jest')

const createJestConfig = nextJest({
  // Path to Next.js app root
  dir: './',
})

const customJestConfig = {
  setupFilesAfterFramework: ['<rootDir>/jest.setup.ts'],
  testEnvironment: 'jest-environment-jsdom',
  moduleNameMapper: {
    '^@/(.*)$': '<rootDir>/$1',
    // CSS modules — return empty object so imports don't crash
    '\\.module\\.(css|sass|scss)$': '<rootDir>/__mocks__/styleMock.js',
  },
  testPathPattern: '\\.(test|spec)\\.(ts|tsx)$',
  collectCoverageFrom: [
    'lib/**/*.{ts,tsx}',
    'components/**/*.{ts,tsx}',
    'app/**/*.{ts,tsx}',
    '!**/*.test.{ts,tsx}',
    '!**/*.d.ts',
  ],
}

module.exports = createJestConfig(customJestConfig)
