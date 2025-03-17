import {defineConfig} from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Kick SDK",
  description: "Powerful Golang toolkit for Kick APIs",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      {text: 'Home', link: '/'},
      {text: 'Overview', link: '/overview/about-kick-sdk'}
    ],

    sidebar: [
      {
        text: 'Introduction',
        items: [
          {text: 'About the Kick SDK', link: '/overview/about-kick-sdk'},
        ]
      },
      {
        text: 'Authorization',
        items: [
          {text: 'About', link: '/overview/about-authorization'},
          {text: 'Authorization Flow', link: '/overview/authorization-flow'},
        ]
      },
      {
        text: 'API',
        items: [
          {text: 'About', link: '/overview/about-api'},
          {
            text: 'Resources',
            items: [
              {text: 'Categories', link: '/overview/api-categories'},
              {text: 'Users', link: '/overview/api-users'},
              {text: 'Channels', link: '/overview/api-channels'},
              {text: 'Chat', link: '/overview/api-chat'},
              {text: 'Public Key', link: '/overview/api-public-key'},
              {text: 'OAuth', link: '/overview/api-oauth'},
              {text: 'Events', link: '/overview/api-events'},
            ],
          },
        ]
      },
      {
        text: 'Events',
        items: [
          {text: 'About', link: '/overview/about-events'},
        ]
      },
    ],

    socialLinks: [
      {icon: 'github', link: 'https://github.com/glichtv/kick-sdk'}
    ],

    search: {
      provider: 'local',
    },
  }
})
