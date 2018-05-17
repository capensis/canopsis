import { normalize } from 'normalizr';

import request from '@/services/request';

import types from './types';

export default {
  actions: {
    async fetch(
      { commit },
      {
        route,
        schema,
        params,
        dataPreparer,
        mutationType,
      },
    ) {
      const [data] = await request.get(route, { params });
      const normalizedData = normalize(dataPreparer(data), schema);

      commit(`entities/${mutationType || types.ENTITIES_UPDATE}`, normalizedData.entities, { root: true });

      return { data, normalizedData };
    },
  },
};
