import js from '@eslint/js'
import pluginVue from 'eslint-plugin-vue'
import pluginPrettier from 'eslint-plugin-prettier'
import configPrettier from 'eslint-config-prettier'
import tsParser from '@typescript-eslint/parser'
import pluginTs from '@typescript-eslint/eslint-plugin'
import globals from 'globals'
import vueEslintParser from 'vue-eslint-parser'

export default [
  {
    name: 'app/files-to-lint',
    files: ['**/*.{js,mjs,jsx,vue,ts,tsx}'],
  },

  {
    name: 'app/files-to-ignore',
    ignores: ['**/dist/**', '**/dist-ssr/**', '**/coverage/**', '**/node_modules/**'],
  },
  js.configs.recommended,
  ...pluginVue.configs['flat/recommended'],
  configPrettier,
  {
    files: ['**/*.vue'],
    languageOptions: {
      parser: vueEslintParser,
      parserOptions: {
        parser: tsParser,
        ecmaVersion: 'latest',
        sourceType: 'module',
      },
      globals: {
        ...globals.browser,
      },
    },
  },
  {
    files: ['**/*.{ts,tsx}'],
    languageOptions: {
      parser: tsParser,
      parserOptions: {
        ecmaVersion: 'latest',
        sourceType: 'module',
      },
      globals: {
        ...globals.browser,
      },
    },
  },
  {
    files: ['**/*.{ts,tsx,vue}'],
    plugins: {
      '@typescript-eslint': pluginTs,
      prettier: pluginPrettier,
    },
    rules: {
      'no-unused-vars': 'off',
      'prettier/prettier': 'error',
      '@typescript-eslint/no-unused-vars': ['error', { argsIgnorePattern: '^_' }],
      '@typescript-eslint/no-explicit-any': 'off',
      '@typescript-eslint/explicit-module-boundary-types': 'off',
      'vue/attributes-order': 'off',
    },
  },
  {
    ignores: ['dist/', 'node_modules/', '*.config.js', '*.config.ts'],
  },
]
