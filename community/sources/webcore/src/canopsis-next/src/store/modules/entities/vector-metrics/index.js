import { API_ROUTES } from '@/config';

import { createWidgetModule } from '@/store/plugins/entities';

export default createWidgetModule({ route: API_ROUTES.metrics.alarm });
