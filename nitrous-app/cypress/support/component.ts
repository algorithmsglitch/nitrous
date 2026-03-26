import React from 'react'
import './commands'
import { mount } from '@cypress/react'

// Fix Next.js CSS injection — must exist before any module loads
if (!document.getElementById('__next_css__DO_NOT_USE__')) {
  const anchor = document.createElement('div')
  anchor.id = '__next_css__DO_NOT_USE__'
  document.head.appendChild(anchor)
}

// Stub next/image — prevents useContext dual-React crash
const NextImage = (props: any) => React.createElement('img', { ...props, src: props.src })
require.cache[require.resolve('next/image')] = {
  id: require.resolve('next/image'),
  filename: require.resolve('next/image'),
  loaded: true,
  exports: NextImage,
} as any

// Stub next/link — prevents useContext dual-React crash
const NextLink = ({ href, children, className }: any) =>
  React.createElement('a', { href, className }, children)
require.cache[require.resolve('next/link')] = {
  id: require.resolve('next/link'),
  filename: require.resolve('next/link'),
  loaded: true,
  exports: NextLink,
} as any

Cypress.Commands.add('mount', mount)

declare global {
  namespace Cypress {
    interface Chainable {
      mount: typeof mount
    }
  }
}