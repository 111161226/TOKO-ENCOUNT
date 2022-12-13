/* eslint-env node */
require('@rushstack/eslint-patch/modern-module-resolution')

module.exports = {
  root: true,
  env: {
    browser: true,
    es6: true,
    node: true
  },
  extends: [
    'plugin:vue/vue3-essential',
    'eslint:recommended',
    '@vue/eslint-config-typescript/recommended',
    '@vue/eslint-config-prettier'
  ],
  parserOptions: {
    parser: '@typescript-eslint/parser'
  },
  rules: {
    'no-unused-vars': 'off',
    '@typescript-eslint/no-unused-vars': [
      'error',
      {
        args: 'none',
        varsIgnorePattern: '^_' // 変数名が "_" で始まる場合は許可
      }
    ],
    '@typescript-eslint/no-inferrable-types': 'off',
    '@typescript-eslint/no-non-null-assertion': 'off',
    '@typescript-eslint/no-explicit-any': 'off',
    'vue/array-bracket-spacing': 'error',
    'vue/arrow-spacing': 'error',
    'vue/block-lang': [
      'error',
      {
        script: { lang: 'ts' },
        style: { lang: 'scss' }
      }
    ],
    'vue/block-spacing': 'error',
    'vue/brace-style': 'error',
    'vue/camelcase': 'error',
    'vue/comma-dangle': 'error',
    'vue/component-api-style': ['error', ['script-setup']],
    'vue/component-name-in-template-casing': ['error', 'kebab-case'],
    'vue/eqeqeq': 'error',
    'vue/key-spacing': 'error',
    'vue/match-component-file-name': 'error',
    'vue/multi-word-component-names': 'off',
    'vue/v-on-event-hyphenation': [
      'error',
      'always',
      {
        autofix: true,
        ignore: []
      }
    ]
  }
}
