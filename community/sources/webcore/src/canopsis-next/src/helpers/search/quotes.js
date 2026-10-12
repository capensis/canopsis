import { escapeRegExp } from 'lodash';

export function replaceTextNotInQuotes(str, text, value) {
  const escapedText = escapeRegExp(text);
  const regexp = new RegExp(`(?!\\B"[^"]*)${escapedText}(?![^"]*"\\B)|(?!\\B'[^']*)${escapedText}(?![^']*'\\B)`, 'gu');

  return str.replace(regexp, value);
}
