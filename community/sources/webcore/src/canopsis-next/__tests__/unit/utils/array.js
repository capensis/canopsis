/**
 * Get random array item
 *
 * @param {Array} items
 */
export const randomArrayItem = items => items[Math.round(Math.random() * (items.length - 1))];
