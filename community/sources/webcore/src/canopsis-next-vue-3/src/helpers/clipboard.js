/**
 * Write text to clipboard
 *
 * @param {*} data
 * @returns {Promise<void>}
 */
export const writeTextToClipboard = data => navigator.clipboard.writeText(data);

/**
 * Read text from clipboard
 *
 * @returns {Promise<string>}
 */
export const readTextFromClipboard = () => navigator.clipboard.readText();
