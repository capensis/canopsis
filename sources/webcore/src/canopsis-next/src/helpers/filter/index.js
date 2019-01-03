import get from 'lodash/get';

import { FILTER_DEFAULT_VALUES } from '@/constants';

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

export default {
  prepareMainFilterToQueryFilter,
};
