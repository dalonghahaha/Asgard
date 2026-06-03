 /// <reference types="vite/client" />

 declare module '*.vue' {
   import type { DefineComponent } from 'vue'
   // eslint-disable-next-line @typescript-eslint/no-explicit-any
   const component: DefineComponent<{}, {}, any>
   export default component
 }

 interface ImportMetaEnv {
   readonly VITE_API_BASE: string
   readonly VITE_BACKEND_TARGET: string
   readonly VITE_TITLE: string
 }

 interface ImportMeta {
   readonly env: ImportMetaEnv
 }
