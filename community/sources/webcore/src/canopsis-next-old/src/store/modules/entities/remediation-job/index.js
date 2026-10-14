import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.remediation.jobs,
  withFetchingParams: true,
  withWithoutStore: true,
}, {
  actions: {
    bulkRemove(context, { data = [] }) {
      return Promise.all(
        data.map(({ _id: id }) => request.delete(
          `${API_ROUTES.remediation.jobs}/${encodeURIComponent(id)}`,
        )),
      );
    },
  },
});
