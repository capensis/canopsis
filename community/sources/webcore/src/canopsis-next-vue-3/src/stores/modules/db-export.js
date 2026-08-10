import { upperFirst } from 'lodash';

import { API_ROUTES } from '@/config';
import request from '@/services/request';

const createExportActions = () =>
  Object.entries(API_ROUTES.dbExport).reduce((acc, [key, value]) => {
    acc[`export${upperFirst(key)}Db`] = ({ ids = [] } = {}) =>
      request.post(value, { ids }, { responseType: 'blob' });

    return acc;
  }, {});

export const dbExportActions = createExportActions();
