import Nav from '@/components/Nav'

describe('Nav Component', () => {
  beforeEach(() => {
    cy.mount(<Nav />)
  })

  describe('Logo Section', () => {
    it('renders the navigation element', () => {
      cy.get('nav').should('exist')
    })

    it('displays the NITROUS logo text', () => {
      cy.get('nav').within(() => {
        cy.contains('NITROUS').should('be.visible')
      })
    })

    it('logo has link to home page', () => {
      cy.get('nav a[href="/"]').should('exist')
      cy.get('nav a[href="/"]').within(() => {
        cy.contains('NITROUS').should('be.visible')
      })
    })

    it('displays logo dot span element', () => {
      cy.get('nav a[href="/"]').within(() => {
        cy.get('span').should('contain', '.')
      })
    })

    it('logo is clickable', () => {
      cy.get('nav a[href="/"]').should('not.be.disabled')
    })
  })

  describe('Navigation Links', () => {
    it('renders Live navigation link', () => {
      cy.get('nav').within(() => {
        cy.contains('Live').should('be.visible')
        cy.get('a').contains('Live').should('have.attr', 'href', '/')
      })
    })

    it('renders Events navigation link', () => {
      cy.get('nav').within(() => {
        cy.contains('Events').should('be.visible')
        cy.get('a').contains('Events').should('have.attr', 'href', '/events')
      })
    })

    it('renders Teams navigation link', () => {
      cy.get('nav').within(() => {
        cy.contains('Teams').should('be.visible')
        cy.get('a').contains('Teams').should('have.attr', 'href', '/teams')
      })
    })

    it('renders Journeys navigation link', () => {
      cy.get('nav').within(() => {
        cy.contains('Journeys').should('be.visible')
        cy.get('a').contains('Journeys').should('have.attr', 'href', '/journeys')
      })
    })

    it('renders Merch navigation link', () => {
      cy.get('nav').within(() => {
        cy.contains('Merch').should('be.visible')
        cy.get('a').contains('Merch').should('have.attr', 'href', '/merch')
      })
    })

    it('all nav links have correct href attributes', () => {
      cy.get('nav a[href="/"]').should('exist')
      cy.get('nav a[href="/events"]').should('exist')
      cy.get('nav a[href="/teams"]').should('exist')
      cy.get('nav a[href="/journeys"]').should('exist')
      cy.get('nav a[href="/merch"]').should('exist')
    })
  })

  describe('Right Section - Status and Button', () => {
    it('displays live status indicator', () => {
      cy.get('[class*="dotLive"]').should('be.visible')
    })

    it('displays live events count text', () => {
      cy.contains('4 Events Live').should('be.visible')
    })

    it('displays status and count together', () => {
      cy.get('[class*="navStatus"]').within(() => {
        cy.get('[class*="dotLive"]').should('exist')
        cy.contains('4 Events Live').should('be.visible')
      })
    })

    it('renders Sign In button', () => {
      cy.get('nav').within(() => {
        cy.get('button').contains('Sign In').should('be.visible')
      })
    })

    it('Sign In button is correct element type', () => {
      cy.get('nav').within(() => {
        cy.get('button').contains('Sign In').should('be.a', 'button')
      })
    })

    it('Sign In button is not disabled', () => {
      cy.get('button').contains('Sign In').should('not.be.disabled')
    })

    it('Sign In button is clickable', () => {
      cy.get('button').contains('Sign In').should('not.have.attr', 'disabled')
    })
  })

  describe('Navigation Structure', () => {
    it('renders navigation with three sections', () => {
      cy.get('nav').within(() => {
        // Logo section
        cy.get('a').contains('NITROUS').should('exist')
        
        // Center nav section
        cy.get('[class*="navCenter"]').should('exist')
        
        // Right section
        cy.get('[class*="navRight"]').should('exist')
      })
    })

    it('renders center nav section with links', () => {
      cy.get('[class*="navCenter"]').within(() => {
        cy.get('a').should('have.length', 5) // Live, Events, Teams, Journeys, Merch
      })
    })

    it('renders right section with status and button', () => {
      cy.get('[class*="navRight"]').within(() => {
        cy.get('[class*="navStatus"]').should('exist')
        cy.get('button').should('exist')
      })
    })

    it('has correct number of navigation links in nav', () => {
      cy.get('nav a').should('have.length', 6) // Logo + 5 nav links
    })
  })

  describe('Styling and Classes', () => {
    it('nav element has nav styling class', () => {
      cy.get('nav[class*="nav"]').should('exist')
    })

    it('logo has logo styling class', () => {
      cy.get('[class*="logo"]').should('exist')
    })

    it('nav links have navLink styling class', () => {
      cy.get('[class*="navLink"]').should('have.length', 5)
    })

    it('sign in button has btnNav styling class', () => {
      cy.get('[class*="btnNav"]').should('exist')
    })

    it('nav status has navStatus styling class', () => {
      cy.get('[class*="navStatus"]').should('exist')
    })
  })

  describe('Accessibility', () => {
    it('all links are navigable', () => {
      cy.get('nav a').each(($link) => {
        cy.wrap($link).should('be.visible')
      })
    })

    it('all interactive elements are keyboard accessible', () => {
      cy.get('nav a').should('have.length.greaterThan', 0)
      cy.get('nav button').should('have.length.greaterThan', 0)
    })

    it('button is not disabled for interaction', () => {
      cy.get('nav button').each(($button) => {
        expect($button.prop('disabled')).to.be.false
      })
    })
  })
})
