import pluginVitest from '@vitest/eslint-plugin';
import skipFormatting from 'eslint-config-prettier/flat';
import pluginOxlint from 'eslint-plugin-oxlint';
import perfectionist from 'eslint-plugin-perfectionist';
import pluginPlaywright from 'eslint-plugin-playwright';
import pluginVue from 'eslint-plugin-vue';
import { globalIgnores } from 'eslint/config';
import tseslint from 'typescript-eslint';

const vueTsRecommended = tseslint.configs.recommended.map((config) => {
  if (!config.files?.includes('**/*.ts')) {
    return config;
  }

  return {
    ...config,
    files: [...config.files, '**/*.vue'],
  };
});

const vueParser = pluginVue.configs['flat/base'][1].languageOptions.parser;

export default tseslint.config(
  {
    name: 'app/files-to-lint',
    files: ['**/*.{vue,js,cjs,mjs,jsx,ts,mts,tsx}'],
  },

  globalIgnores([
    '**/build/**',
    '**/dist/**',
    '**/dist-ssr/**',
    '**/coverage/**',
    '**/public/scripts/**',
  ]),

  ...pluginVue.configs['flat/essential'],
  ...vueTsRecommended,

  ...pluginVue.configs['flat/base'],
  {
    name: 'app/typescript-in-vue',
    files: ['*.vue', '**/*.vue'],
    languageOptions: {
      parser: vueParser,
      parserOptions: {
        parser: {
          js: 'espree',
          jsx: 'espree',
          ts: tseslint.parser,
          tsx: tseslint.parser,
        },
        ecmaVersion: 2024,
        ecmaFeatures: {
          jsx: true,
        },
        extraFileExtensions: ['.vue'],
      },
    },
    rules: {
      'vue/block-lang': [
        'error',
        {
          script: {
            lang: ['js', 'jsx', 'ts', 'tsx'],
            allowNoLang: true,
          },
        },
      ],
    },
  },

  {
    ...pluginPlaywright.configs['flat/recommended'],
    files: ['e2e/**/*.{test,spec}.{js,ts,jsx,tsx}'],
  },

  {
    ...pluginVitest.configs.recommended,
    files: ['__tests__/**/*.{js,jsx,ts,tsx}'],
  },

  ...pluginOxlint.buildFromOxlintConfigFile('.oxlintrc.json'),

  skipFormatting,

  {
    name: 'app/legacy-canopsis-rules',
    plugins: {
      perfectionist,
    },
    languageOptions: {
      globals: {
        sel: 'readonly',
      },
    },
    rules: {
      'object-curly-newline': 'off',
      'arrow-parens': ['error', 'as-needed', { requireForBlockBody: true }],
      'arrow-body-style': ['error', 'as-needed'],
      'implicit-arrow-linebreak': 'off',
      'brace-style': ['error', '1tbs', { allowSingleLine: false }],
      curly: ['error', 'all'],
      eqeqeq: ['error', 'always'],
      semi: ['error', 'always'],
      quotes: ['error', 'single', { avoidEscape: false }],
      'default-case': 'off',
      'max-len': [
        'warn',
        {
          code: 120,
          ignoreStrings: true,
          ignoreTemplateLiterals: true,
          ignoreRegExpLiterals: true,
        },
      ],
      'no-param-reassign': [
        'error',
        {
          props: true,
          ignorePropertyModificationsFor: [
            'state',
            'config',
            'acc',
            'Vue',
            'browser',
            'wrapper',
            'app',
            'instance',
          ],
        },
      ],
      'no-console': ['error', { allow: ['warn', 'error'] }],
      'no-return-assign': 'off',
      'no-underscore-dangle': [
        'error',
        {
          allow: [
            '_id',
            '_embedded',
            '_t',
            '_filter',
            'type_',
            '_original',
            '_modulesNamespaceMap',
            '_handlebarsRequestHelper',
            '__vueParentComponent',
            '__TEST_STORE__',
            '__TEST_FORM__',
            '__TEST_MODALS__',
            '__TEST_POPUPS__',
            '__TEST_SIDEBAR__',
            '__TEST_ALARM_TAG_HOOK__',
            '__TEST_PATTERN_HOOK__',
            '__TEST_CORPORATE_PATTERN_HOOK__',
            '__TEST_SELECT_ADVANCED_SEARCH_FIELD__',
          ],
        },
      ],
      'no-bitwise': ['error', { allow: ['&', '>>'] }],
      'no-continue': 'off',
      'no-cond-assign': 'off',
      'no-promise-executor-return': 'off',
      'default-param-last': 'off',
      'no-restricted-exports': 'off',
      'prefer-regex-literals': ['error', { disallowRedundantWrapping: true }],
      'no-unused-expressions': 'off',
      'no-restricted-syntax': 'off',
      'no-restricted-imports': [
        'error',
        {
          patterns: [
            {
              group: ['**/form/partials/**', '**/form/partial/**'],
              message: 'All fields should be moved inside `/form/fields/`.',
            },
            {
              group: ['**/field/**'],
              message: '`field` folder deprecated. Please, rename to `fields`.',
            },
            {
              group: ['**/partial/**'],
              message: '`partial` folder deprecated. Please, rename to `partials`.',
            },
          ],
        },
      ],
      'perfectionist/sort-imports': [
        'error',
        {
          type: 'natural',
          order: 'asc',
          internalPattern: ['^@/.*'],
          customGroups: [
            {
              groupName: 'internal-config-constants',
              elementNamePattern: '^@/(config|constants)$',
            },
            {
              groupName: 'internal-plugins',
              elementNamePattern: '^@/plugins(/.*)?$',
            },
            {
              groupName: 'internal-services',
              elementNamePattern: '^@/services(/.*)?$',
            },
            {
              groupName: 'internal-core',
              elementNamePattern: '^@/(store|stores|i18n|filters|router)$',
            },
            {
              groupName: 'internal-store',
              elementNamePattern: '^@/(store|stores)/.*$',
            },
            {
              groupName: 'internal-helpers',
              elementNamePattern: '^@/helpers/.*$',
            },
            {
              groupName: 'internal-mixins-hooks',
              elementNamePattern: '^@/(mixins|hooks)/.*$',
            },
            {
              groupName: 'internal-component-hooks',
              elementNamePattern: '^@/components/.*/hooks/.*$',
            },
            {
              groupName: 'internal-components',
              elementNamePattern: '^@/components/.*$',
            },
          ],
          groups: [
            'builtin',
            { newlinesBetween: 1 },
            'external',
            { newlinesBetween: 1 },
            'internal-config-constants',
            { newlinesBetween: 1 },
            'internal-plugins',
            { newlinesBetween: 1 },
            'internal-services',
            { newlinesBetween: 1 },
            'internal-core',
            { newlinesBetween: 1 },
            'internal-store',
            { newlinesBetween: 1 },
            'internal-helpers',
            { newlinesBetween: 1 },
            'internal-mixins-hooks',
            { newlinesBetween: 1 },
            'internal-component-hooks',
            { newlinesBetween: 1 },
            'internal-components',
            { newlinesBetween: 1 },
            'internal',
            { newlinesBetween: 1 },
            'parent',
            { newlinesBetween: 1 },
            'sibling',
            { newlinesBetween: 1 },
            'index',
          ],
          newlinesBetween: 1,
          sortSideEffects: false,
        },
      ],
      'vue/valid-next-tick': 'off',
      'vue/no-deprecated-model-definition': 'off',
      'vue/multi-word-component-names': 'off',
      'vue/no-use-computed-property-like-method': 'off',
      'vue/component-definition-name-casing': 'off',
      'vue/require-default-prop': 'off',
      'vue/require-prop-types': 'error',
      'vue/no-unused-components': [
        'error',
        {
          ignoreWhenBindingPresent: true,
        },
      ],
      'vue/order-in-components': [
        'warn',
        {
          order: [
            'el',
            'name',
            'parent',
            'functional',
            '$_veeValidate',
            ['provide', 'inject'],
            ['delimiters', 'comments'],
            ['components', 'directives', 'filters'],
            'extends',
            'mixins',
            'inheritAttrs',
            'model',
            ['props', 'propsData'],
            'data',
            'computed',
            'watch',
            'LIFECYCLE_HOOKS',
            'methods',
            ['template', 'render'],
            'renderError',
          ],
        },
      ],
      'vue/block-order': ['error', { order: ['template', 'script', 'style'] }],
      'vue/no-v-text-v-html-on-component': 'off',
      'vue/require-component-is': 'off',
      'vue/max-attributes-per-line': [
        'warn',
        {
          singleline: 3,
          multiline: 1,
        },
      ],
      'vue/first-attribute-linebreak': [
        'warn',
        {
          singleline: 'beside',
        },
      ],
      'vue/singleline-html-element-content-newline': 'off',
      'vue/component-name-in-template-casing': [
        'warn',
        'kebab-case',
        {
          registeredComponentsOnly: true,
        },
      ],
      'vue/eqeqeq': ['error', 'always'],
      'vue/attributes-order': [
        'warn',
        {
          order: [
            'GLOBAL',
            'CONDITIONALS',
            'LIST_RENDERING',
            'DEFINITION',
            'TWO_WAY_BINDING',
            'CONTENT',
            'OTHER_DIRECTIVES',
            'RENDER_MODIFIERS',
            'UNIQUE',
            'SLOT',
            'ATTR_DYNAMIC',
            'ATTR_STATIC',
            'ATTR_SHORTHAND_BOOL',
            'EVENTS',
          ],
          alphabetical: false,
        },
      ],
      'vue/valid-v-slot': ['error', { allowModifiers: true }],
    },
  },

  {
    name: 'app/migration-rules',
    rules: {
      '@typescript-eslint/no-unused-vars': ['warn', { ignoreRestSiblings: true }],
      'vue/multi-word-component-names': 'off',
      'no-unused-vars': 'off',
    },
  },

  {
    name: 'app/tools-rules',
    files: ['tools/**/*.js'],
    rules: {
      'no-console': 'off',
    },
  },

  {
    name: 'app/required-semicolons',
    rules: {
      semi: ['error', 'always'],
    },
  },
);
