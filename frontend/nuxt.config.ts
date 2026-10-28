// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:5665/api'
    }
  },
  modules: ['@nuxtjs/tailwindcss'],
  css: ['~/assets/css/print.css'],
  postcss: {
    plugins: {
      'tailwindcss/nesting': false
    }
  },
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  nitro: {
    externals: {
      inline: [/[\\/]node_modules[\\/]nuxt[\\/]dist[\\/]/]
    }
  }
})
