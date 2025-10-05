import { withMermaid } from 'vitepress-plugin-mermaid'

export default withMermaid({
  mermaid: {
    theme: 'default',
    securityLevel: 'loose',
    flowchart: {
      useMaxWidth: true,
      htmlLabels: true
    }
  },
  lang: 'en-EN',
  title: 'Doc-Builder example doc 1',
  description: 'Doc-Builder example doc 1',
  base: '/',
  cleanUrls: true,
  ignoreDeadLinks: true,
  themeConfig: {
    nav: [
      { text: 'Home page', link: '/' }
    ],
    sidebar: [],
    lastUpdated: true,
    search: {
      provider: 'local'
    }
  }
})
