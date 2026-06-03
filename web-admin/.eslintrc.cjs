 /* eslint-env node */
 module.exports = {
   root: true,
   env: { browser: true, es2022: true, node: true },
   parser: 'vue-eslint-parser',
   parserOptions: {
     parser: '@typescript-eslint/parser',
     ecmaVersion: 2022,
     sourceType: 'module',
     extraFileExtensions: ['.vue'],
   },
   extends: [
     'eslint:recommended',
     'plugin:vue/vue3-recommended',
     'plugin:@typescript-eslint/recommended',
     'prettier',
   ],
   plugins: ['@typescript-eslint', 'vue'],
   rules: {
     'vue/multi-word-component-names': 'off',
     'vue/require-default-prop': 'off',
     '@typescript-eslint/consistent-type-imports': 'warn',
     '@typescript-eslint/no-unused-vars': ['warn', { argsIgnorePattern: '^_' }],
     '@typescript-eslint/no-explicit-any': 'warn',
   },
   ignorePatterns: ['dist', 'node_modules', 'src/auto-imports.d.ts', 'src/components.d.ts'],
 }
