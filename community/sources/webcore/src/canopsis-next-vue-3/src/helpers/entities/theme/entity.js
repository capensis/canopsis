import { isString, kebabCase } from 'lodash'

const themeObjectColorsToCSSVariables = (colors, prefix = '') => Object.entries(colors).reduce((acc, [key, value]) => {
  if (!value) {
    return acc
  }

  if (isString(value) || typeof value === 'number') {
    acc[`${prefix}${kebabCase(key)}`] = value

    return acc
  }

  return {
    ...acc,
    ...themeObjectColorsToCSSVariables(value, `${key}-`),
  }
}, {})

export const themePropertiesToCSSVariables = (colors = {}) => themeObjectColorsToCSSVariables(colors)

export const getThemeVariableKeyFromCssVariable = (cssVariable = '') => (
  cssVariable.replace('var(--v-', '').replace('-base)', '')
)
