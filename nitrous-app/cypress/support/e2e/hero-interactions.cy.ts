describe('Hero Section Interactions', () => {
  beforeEach(() => {
    cy.visit('/')
  })

  it('displays hero action buttons and verifies they are clickable', () => {
    const igniteButton = cy.contains('Ignite Stream')
    igniteButton.should('be.visible')
    igniteButton.should('not.be.disabled')
  })

  it('verifies all navigation cards are visible in hero section', () => {
    // Access Garage card
    cy.contains(/ACCESS.*GARAGE/).should('be.visible')
    // Access Event Passes card
    cy.contains(/ACCESS.*EVENT PASSES/).should('be.visible')
    // Access Live Streams card
    cy.contains(/ACCESS.*LIVE STREAMS/).should('be.visible')
    // Teams card
    cy.contains(/TEAMS/).should('be.visible')
    // Journeys card
    cy.contains(/JOURNEYS/).should('be.visible')
    // Merch card
    cy.contains(/MERCH/).should('be.visible')
  })

  it('can click the Ignite Stream button', () => {
    cy.contains('Ignite Stream').click()
    // Verify the action completes (page doesn't error)
    cy.contains('NITROUS').should('be.visible')
  })

  it('can click the Explore Events button', () => {
    cy.contains('Explore Events').click()
    // Verify the action completes (page doesn't error)
    cy.contains('NITROUS').should('be.visible')
  })

  it('can navigate through hero nav cards', () => {
    // Click on the garage link from hero nav
    cy.contains(/ACCESS.*GARAGE/).click()
    cy.url().should('include', '/garage')
  })

  it('can navigate to live streams from hero', () => {
    cy.contains(/ACCESS.*LIVE STREAMS/).click()
    cy.url().should('include', '/live')
  })

  it('verifies hero section styling elements exist', () => {
    // Check for hero section existence
    cy.get('h1').should('be.visible')
    // Verify buttons are rendered
    cy.get('button').should('have.length.greaterThan', 0)
  })
})
