import escapeStringRegexp from 'escape-string-regexp';

export function replaceTextNotInQuotes(str, text, value) {
  const escapedText = escapeStringRegexp(text);
  const regexp = new RegExp(`(?!\\B"[^"]*)${escapedText}(?![^"]*"\\B)|(?!\\\\B'[^']*)${escapedText}(?![^']*'\\\\B)`, 'gu');

  return str.replace(regexp, value);
}
