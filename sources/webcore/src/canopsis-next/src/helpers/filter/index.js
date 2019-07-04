import { get, isString } from 'lodash';

import { FILTER_DEFAULT_VALUES } from '@/constants';

export default function prepareMainFilterToQueryFilter(
  filterObject = {},
  condition = FILTER_DEFAULT_VALUES.condition,
) {
  let filter;

  if (Array.isArray(filterObject) && filterObject.length) {
    filter = {
      [condition]: filterObject.map((item) => {
        if (isString(item.filter)) {
          return JSON.parse(item.filter);
        }

        return item.filter || {};
      }),
    };
  } else if (!Array.isArray(filterObject)) {
    filter = get(filterObject, 'filter');
  }

  return filter;
}

