import { defineConfig } from 'vitepress'

export default defineConfig({
  lang: 'en-US',
  title: 'Doc Builder Example',
  description: 'Sample site assembled with doc-builder',
  ignoreDeadLinks: true,
  themeConfig: {
    nav: [
      { text: 'Home', link: '/' }
    ],
    sidebar: [
      {
        text: '🏠 Home',
        link: '/'
      },
      // SIDEBAR_ITEMS - will be replaced by build script
    ]
  }
})
