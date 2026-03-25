describe('Login page', () => {
  beforeEach(() => {
    cy.visit('/login')
  })

  it('loads the login page', () => {
    cy.contains('Sign In').should('be.visible')
  })

  it('renders email and password inputs', () => {
    cy.get('[data-testid="email-input"]').should('be.visible')
    cy.get('[data-testid="password-input"]').should('be.visible')
  })

  it('renders the submit button', () => {
    cy.get('[data-testid="submit-button"]').should('be.visible')
  })

  it('shows an error message on invalid credentials', () => {
    cy.intercept('POST', '**/api/auth/login', {
      statusCode: 401,
      body: { error: 'Invalid credentials' },
    }).as('loginFail')

    cy.get('[data-testid="email-input"]').type('wrong@test.com')
    cy.get('[data-testid="password-input"]').type('badpassword')
    cy.get('[data-testid="login-form"]').submit()

    cy.wait('@loginFail')
    cy.get('[data-testid="error-message"]').should('be.visible')
    cy.contains(/invalid credentials/i).should('be.visible')
  })

  it('redirects to homepage on successful login', () => {
    cy.intercept('POST', '**/api/auth/login', {
      statusCode: 200,
      body: {
        token: 'fake-jwt-token',
        user: { id: '1', name: 'Test User', email: 'test@test.com' },
      },
    }).as('loginSuccess')

    cy.get('[data-testid="email-input"]').type('test@test.com')
    cy.get('[data-testid="password-input"]').type('password123')
    cy.get('[data-testid="login-form"]').submit()

    cy.wait('@loginSuccess')
    cy.url().should('eq', Cypress.config().baseUrl + '/')
  })

  it('stores JWT token in localStorage after login', () => {
    cy.intercept('POST', '**/api/auth/login', {
      statusCode: 200,
      body: { token: 'my-test-token', user: { id: '1', name: 'Test User' } },
    }).as('loginSuccess')

    cy.get('[data-testid="email-input"]').type('test@test.com')
    cy.get('[data-testid="password-input"]').type('password123')
    cy.get('[data-testid="login-form"]').submit()

    cy.wait('@loginSuccess')
    cy.window().then((win) => {
      expect(win.localStorage.getItem('nitrous_token')).to.eq('my-test-token')
    })
  })

  it('NITROUS logo links back to homepage', () => {
    cy.contains('NITROUS').should('have.attr', 'href', '/')
  })

  it('Register link points to /register', () => {
    cy.contains('Register').should('have.attr', 'href', '/register')
  })

  it('email input accepts typed text', () => {
    cy.get('[data-testid="email-input"]').type('test@test.com')
    cy.get('[data-testid="email-input"]').should('have.value', 'test@test.com')
  })

  it('password input masks typed text', () => {
    cy.get('[data-testid="password-input"]').should('have.attr', 'type', 'password')
  })
})

describe('Nav bar on homepage', () => {
  beforeEach(() => {
    cy.visit('/')
  })

  it('renders the NITROUS logo', () => {
    cy.contains('NITROUS').should('be.visible')
  })

  it('renders all main nav links', () => {
    cy.contains('Events').should('be.visible')
    cy.contains('Teams').should('be.visible')
    cy.contains('Journeys').should('be.visible')
    cy.contains('Merch').should('be.visible')
  })

  it('Sign In button is present', () => {
    cy.contains('Sign In').should('be.visible')
  })

  it('clicking Sign In navigates to /login', () => {
    // Wire up the Sign In button — currently it has no handler,
    // so this test will fail until Nav is updated to link to /login.
    // Update Nav.tsx: change <button> to <Link href="/login">
    cy.contains('Sign In').click()
    cy.url().should('include', '/login')
  })
})
