/**
 * @returns {boolean}
 */
export const isMacOS = () => /Mac|iPhone|iPod|iPad/.test(navigator.platform || '');

/**
 * Modifier used for alarm-list drag selection: Command on macOS, Control elsewhere.
 * macOS treats Control + left click as a right click.
 *
 * @param {KeyboardEvent|MouseEvent} event
 * @returns {boolean}
 */
export const hasAlarmSelectingModifier = event => (isMacOS() ? event.metaKey : event.ctrlKey);

/**
 * @param {string} key
 * @returns {boolean}
 */
export const isAlarmSelectingModifierKey = key => key === (isMacOS() ? 'Meta' : 'Control');
