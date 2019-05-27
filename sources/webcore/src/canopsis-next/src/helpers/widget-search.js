/**
 * Get context filter by text
 *
 * @param {string} text
 * @param {Array} [additionalFields=[]] - if we want we can put additional fields (ex: ['_id'])
 * @returns {{$or: *[]}}
 */
export function getContextSearchByText(text, additionalFields = []) {
  const orConditions = [
    {
      name: { $regex: text, $options: 'i' },
    },
    {
      type: { $regex: text, $options: 'i' },
    },
  ];

  if (additionalFields.length) {
    const additionalOrConditions = additionalFields.map(field => ({
      [field]: { $regex: text, $options: 'i' },
    }));

    orConditions.push(...additionalOrConditions);
  }

  return {
    $or: orConditions,
  };
}

export default {
  getContextSearchByText,
};
