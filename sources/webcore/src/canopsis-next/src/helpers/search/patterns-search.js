function getRule(text = '') {
  return {
    $regex: text,
    $options: 'i',
  };
}

function getRulesByKeys(keys = [], text = '') {
  return keys.map(key => ({ [key]: getRule(text) }));
}

export function getContextWidgetSearchByText(text) {
  return {
    $or: getRulesByKeys(['name', 'type'], text),
  };
}

export function getUsersSearchByText(text) {
  return {
    $or: getRulesByKeys([
      '_id',
      'role',
      'username',
      'firstname',
      'lastname',
      'mail',
      'enable',
      'external',
    ], text),
  };
}

export function getRolesSearchByText(text) {
  return {
    $or: getRulesByKeys([
      '_id',
      'role',
      'enable',
      'description',
    ], text),
  };
}
