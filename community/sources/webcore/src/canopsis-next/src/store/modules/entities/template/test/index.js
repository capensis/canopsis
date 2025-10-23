import { API_ROUTES } from '@/config';

import { createBasicCRUDModule } from '@/store/plugins/entities';

export default createBasicCRUDModule({
  route: API_ROUTES.templateTests,
  withWithoutStore: true,
});
