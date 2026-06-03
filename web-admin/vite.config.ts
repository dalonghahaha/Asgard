import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import { fileURLToPath, URL } from 'node:url'

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const backendTarget = env.VITE_BACKEND_TARGET || 'http://localhost:12345'
  return {
    plugins: [
      vue(),
      AutoImport({
        resolvers: [ElementPlusResolver()],
        imports: ['vue', 'vue-router', 'pinia'],
        dts: 'src/auto-imports.d.ts',
        eslintrc: { enabled: true },
      }),
      // T-202:关闭 components.d.ts 自动生成,组件类型由 src/types/element-plus.d.ts
      // 统一维护。原因:自动生成的 `typeof import('element-plus/es')['ElTableColumn']`
      // 暴露运行时 buildProps 的 `PropType<EpPropMergeType<...>>` 复杂对象,vue-tsc 会
      // 拿这个去校验模板字面量 `prop="id"`,触发 TS2322;slot 里的 `row` 也被推为
      // 内部 DefaultRow。维护成本:新增本项目自研 .vue 组件时需要在 element-plus.d.ts
      // 追加一行。
      Components({
        resolvers: [ElementPlusResolver()],
        dts: false,
      }),
    ],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      },
    },
    server: {
      port: 5173,
      host: '0.0.0.0',
      // T-203:开发期把 /api/* 转发到后端 Asgard web
      proxy: {
        '/api': {
          target: backendTarget,
          changeOrigin: true,
        },
      },
    },
    build: {
      outDir: 'dist',
      sourcemap: false,
      chunkSizeWarningLimit: 1500,
    },
  }
})
