/**
 * Get base step class by tour name and step number
 *
 * @param {string} tour - Tour name
 * @param {number} step - Step number
 * @returns {string}
 */
export function getStepClass(tour, step) {
  return `v-tour-${tour}-step-${step}`;
}

/**
 * Get base step query selector by tour name and step number
 *
 * @param {string} tour - Tour name
 * @param {number} step - Step number
 * @returns {string}
 */
export function getStepTarget(tour, step) {
  return `.${getStepClass(tour, step)}`;
}
