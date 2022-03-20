import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import vueJsx from "@vitejs/plugin-vue-jsx";
import VitePluginLinaria from 'vite-plugin-linaria';

// https://vitejs.dev/config/
export default defineConfig({
  base: process.env.STATIC || "/",
  plugins: [
    vue(),
    vueJsx(),
    VitePluginLinaria(),
  ],
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          common: [
            "axios",
            "dayjs",
            "localforage",
            "pako",
          ],
          ui: [
            "vue",
            "vue-router",
          ],
          naive: [
            "naive-ui",
          ]
        },
      },
    },
  },
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:7001",
        rewrite: (path) => path.replace(/^\/api/, ""),
      },
    },
  },
});
