/**
 * This function converts userPreference with widget type 'Map' to query Object
 *
 * @param {Object} userPreference
 * @returns {{ category: string }}
 */
export const convertMapUserPreferenceToQuery = ({ content: { category } }) => ({ category });
