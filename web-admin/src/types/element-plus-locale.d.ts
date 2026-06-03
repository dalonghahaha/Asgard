// Element Plus locale 子路径缺类型声明（element-plus 没为 dist/locale/*.mjs 出 .d.ts）
// 文件保持纯脚本形式（没有 top-level import/export），让 `declare module` 当作新模块声明
// 而不是 augmentation，这样 `import zhCn from 'element-plus/dist/locale/zh-cn.mjs'` 才会被识别。
declare module 'element-plus/dist/locale/zh-cn.mjs' {
  const locale: import('element-plus/es/locale').Language
  export default locale
}

declare module 'element-plus/dist/locale/en.mjs' {
  const locale: import('element-plus/es/locale').Language
  export default locale
}
