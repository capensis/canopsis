import { BACKGROUND_AND_ICONS_COLORS_DIFF } from '@/constants'

/**
 * @typedef {Object} ReadableColorOptions
 * @property {'AA' | 'AAA'} level
 * @property {'small' | 'large'} size
 */

const NAMED_COLORS = {
  black: '#000000',
  white: '#FFFFFF',
}

const clamp = (value, min = 0, max = 255) => Math.min(max, Math.max(min, value))

const expandHex = (value) =>
  value.length === 3
    ? value
        .split('')
        .map((char) => char + char)
        .join('')
    : value

const parseHexColor = (color = '') => {
  const value = expandHex(String(color).replace('#', '').slice(0, 6))
  const number = Number.parseInt(value, 16)

  if (Number.isNaN(number)) {
    return null
  }

  return {
    r: (number >> 16) & 255,
    g: (number >> 8) & 255,
    b: number & 255,
  }
}

const parseRgbColor = (color = '') => {
  const [, r, g, b, alpha] =
    String(color).match(/rgba?\(([^,]+),\s*([^,]+),\s*([^,)]+)(?:,\s*([^)]+))?\)/) ?? []

  if (r === undefined || g === undefined || b === undefined) {
    return null
  }

  return {
    r: Number.parseFloat(r),
    g: Number.parseFloat(g),
    b: Number.parseFloat(b),
    a: alpha === undefined ? undefined : Number.parseFloat(alpha),
  }
}

export const parseColor = (color) => {
  if (!color) {
    return null
  }

  if (typeof color === 'object') {
    return color
  }

  const preparedColor = NAMED_COLORS[String(color).toLowerCase()] ?? color

  return String(preparedColor).startsWith('#')
    ? parseHexColor(preparedColor)
    : parseRgbColor(preparedColor)
}

const componentToHex = (value) => clamp(Math.round(value)).toString(16).padStart(2, '0')

const colorToString = ({ r, g, b }) =>
  `#${componentToHex(r)}${componentToHex(g)}${componentToHex(b)}`

const rgbToHsl = ({ r, g, b }) => {
  const normalizedR = r / 255
  const normalizedG = g / 255
  const normalizedB = b / 255
  const max = Math.max(normalizedR, normalizedG, normalizedB)
  const min = Math.min(normalizedR, normalizedG, normalizedB)
  const lightness = (max + min) / 2

  if (max === min) {
    return { h: 0, s: 0, l: lightness }
  }

  const delta = max - min
  const saturation = lightness > 0.5 ? delta / (2 - max - min) : delta / (max + min)
  const hue =
    {
      [normalizedR]: (normalizedG - normalizedB) / delta + (normalizedG < normalizedB ? 6 : 0),
      [normalizedG]: (normalizedB - normalizedR) / delta + 2,
      [normalizedB]: (normalizedR - normalizedG) / delta + 4,
    }[max] / 6

  return { h: hue, s: saturation, l: lightness }
}

const hueToRgb = (p, q, t) => {
  let hue = t

  if (hue < 0) {
    hue += 1
  }

  if (hue > 1) {
    hue -= 1
  }

  if (hue < 1 / 6) {
    return p + (q - p) * 6 * hue
  }

  if (hue < 1 / 2) {
    return q
  }

  if (hue < 2 / 3) {
    return p + (q - p) * (2 / 3 - hue) * 6
  }

  return p
}

const hslToRgb = ({ h, s, l }) => {
  if (s === 0) {
    return {
      r: l * 255,
      g: l * 255,
      b: l * 255,
    }
  }

  const q = l < 0.5 ? l * (1 + s) : l + s - l * s
  const p = 2 * l - q

  return {
    r: hueToRgb(p, q, h + 1 / 3) * 255,
    g: hueToRgb(p, q, h) * 255,
    b: hueToRgb(p, q, h - 1 / 3) * 255,
  }
}

const adjustColor = (color, amount, direction) => {
  const rgb = parseColor(color)

  if (!rgb) {
    return color
  }

  const hsl = rgbToHsl(rgb)
  const ratio = direction === 'lighten' ? amount / 100 : -amount / 100

  return colorToString(
    hslToRgb({
      ...hsl,
      l: clamp(hsl.l + ratio, 0, 1),
    }),
  )
}

/**
 * Check color is readable
 *
 * @param {string} firstColor
 * @param {string} secondColor
 * @param {ReadableColorOptions} [options = {}]
 * @return {boolean}
 */
export const isReadableColor = (firstColor, secondColor) => {
  const first = parseColor(firstColor)
  const second = parseColor(secondColor)

  if (!first || !second) {
    return false
  }

  const contrast = (getRelativeLuminance(first) + 0.05) / (getRelativeLuminance(second) + 0.05)
  const normalizedContrast = contrast < 1 ? 1 / contrast : contrast

  return normalizedContrast >= 4.5
}

/**
 * Get most readable text color ('white' or 'black')
 *
 * @param {string} color
 * @param {ReadableColorOptions} [options = {}]
 */
export const getMostReadableTextColor = (color) =>
  isReadableColor(color, 'white') ? 'white' : 'black'

/**
 * Convert color to rgb
 *
 * @param {string|Object} color
 * @return {string}
 */
export const colorToRgb = (color) => {
  const rgb = parseColor(color)

  return rgb ? `rgb(${rgb.r}, ${rgb.g}, ${rgb.b})` : color
}

/**
 * Convert color to rgba with alpha
 *
 * @param {string|Object} color
 * @param {number} alpha
 * @return {string}
 */
export const colorToRgba = (color, alpha = 1.0) => {
  const rgb = parseColor(color)

  return rgb ? `rgba(${rgb.r}, ${rgb.g}, ${rgb.b}, ${alpha})` : color
}

/**
 * Convert color to hex
 *
 * @param {string|Object} color
 * @return {string}
 */
export const colorToHex = (color) => {
  const rgb = parseColor(color)

  return rgb ? colorToString(rgb) : color
}

/**
 * Convert color to hex8
 *
 * @param {string|Object} color
 * @return {string}
 */
export const colorToHex8 = (color, alpha = 1) =>
  `${colorToHex(color)}${componentToHex(alpha * 255)}`

/**
 * Check color is valid
 *
 * @param {string|Object} color
 * @return {boolean}
 */
export const isValidColor = (color) => Boolean(parseColor(color))

const getRelativeLuminance = ({ r, g, b }) => {
  const channels = [r, g, b].map((value) => {
    const normalized = value / 255

    return normalized <= 0.03928 ? normalized / 12.92 : ((normalized + 0.055) / 1.055) ** 2.4
  })

  return channels[0] * 0.2126 + channels[1] * 0.7152 + channels[2] * 0.0722
}

/**
 * Check color is dark
 *
 * @param {string|Object} color
 * @return {boolean}
 */
export const isDarkColor = (color) => {
  const rgb = parseColor(color)

  return rgb ? getRelativeLuminance(rgb) < 0.18 : false
}

/**
 * Get darken color
 *
 * @param {string} color
 * @param {number} amount
 */
export const getDarkenColor = (color, amount) => adjustColor(color, amount, 'darken')

/**
 * Get lighten color
 *
 * @param {string} color
 * @param {number} amount
 */
export const getLightenColor = (color, amount) => adjustColor(color, amount, 'lighten')

/**
 * Adjusts the color by either lightening or darkening it based on its current brightness.
 *
 * @param {string|Object} color - The color to be adjusted. It can be a string or an object compatible with tinycolor.
 * @param {number} [amount=BACKGROUND_AND_ICONS_COLORS_DIFF] - The amount by which to lighten or darken the color.
 * @returns {string} - The adjusted color as a string.
 */
export const getLightenOrDarkenColor = (color, amount = BACKGROUND_AND_ICONS_COLORS_DIFF) =>
  isDarkColor(color) ? getLightenColor(color, amount) : getDarkenColor(color, amount)

/**
 * Check property is css variable
 *
 * @param {string} property
 * @returns {boolean}
 */
export const isCSSVariable = (property) => /^var\(.+\)$/.test(property)

/**
 * Get css variable name
 *
 * @param {string} property
 * @returns {string}
 */
export const getCSSVariableName = (property) => property.match(/^var\((.+)\)$/)[1]

/**
 * Get color from css variable
 *
 * @param {Element} element
 * @param {string} property
 */
export const getCSSVariableColor = (element, property) =>
  getComputedStyle(element).getPropertyValue(property)
