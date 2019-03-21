import { get } from 'lodash';

import { FILTER_DEFAULT_VALUES } from '@/constants';

import parseGroupToFilter from './editor/parse-group-to-filter';
import parseFilterToRequest from './editor/parse-filter-to-request';

export function prepareMainFilterToQueryFilter(
  filterObject = {},
  condition = FILTER_DEFAULT_VALUES.condition,
) {
  let filter;

  if (Array.isArray(filterObject) && filterObject.length) {
    filter = { [condition]: filterObject.map(item => JSON.parse(item.filter)) };
  } else if (!Array.isArray(filterObject)) {
    filter = get(filterObject, 'filter');
  }

  return filter;
}

export function addPrefixIntoGroupFields(filterGroup, prefix) {
  const groups = Object.keys(filterGroup.groups).reduce((acc, key) => {
    const group = filterGroup.groups[key];

    if (group) {
      acc[key] = addPrefixIntoGroupFields(group, prefix);
    }

    return acc;
  }, {});

  const rules = Object.keys(filterGroup.rules).reduce((acc, key) => {
    const rule = filterGroup.rules[key];

    acc[key] = {
      ...rule,

      field: prefix + rule.field,
    };

    return acc;
  }, {});

  return {
    rules,
    groups,
    condition: filterGroup.condition,
  };
}

export function prepareFilterWithFieldsPrefix(filterObject) {
  const formattedFilter = parseGroupToFilter(filterObject);
  const formattedFilterWithPrefixes = addPrefixIntoGroupFields(formattedFilter, 'entity.');

  return parseFilterToRequest(formattedFilterWithPrefixes);
}

export default {
  prepareMainFilterToQueryFilter,
};
