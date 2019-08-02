export function getRule(text) {
  return {
    $regex: text,
    $options: 'i',
  };
}

export function getContextWidgetSearchByText(text) {
  return {
    $or: [
      { name: getRule(text) },
      { type: getRule(text) },
    ],
  };
}

export function getUsersSearchByText(text) {
  return {
    $or: [
      { _id: getRule(text) },
      { role: getRule(text) },
      { enable: getRule(text) },
      { external: getRule(text) },
    ],
  };
}

export function getRolesSearchByText(text) {
  return {
    $or: [
      { _id: getRule(text) },
      { role: getRule(text) },
      { enable: getRule(text) },
      { description: getRule(text) },
    ],
  };
}

export default {
  getContextWidgetSearchByText,
  getUsersSearchByText,
  getRolesSearchByText,
};
