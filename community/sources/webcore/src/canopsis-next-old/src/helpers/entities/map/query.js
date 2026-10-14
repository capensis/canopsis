/**
 * This function converts widget with type 'Map' to query Object
 *
 * @param {Object} widget
 * @returns {{ lockedFilter: string }}
 */
export const convertMapWidgetToQuery = (widget) => {
  const { mainFilter } = widget.parameters;

  return {
    lockedFilter: mainFilter,
  };
};

/**
 * This function converts userPreference with widget type 'Map' to query Object
 *
 * @param {Object} userPreference
 * @returns {{ category: string }}
 */
export const convertMapUserPreferenceToQuery = ({ content: { category, mainFilter } }) => ({
  category,
  filter: mainFilter,
});
