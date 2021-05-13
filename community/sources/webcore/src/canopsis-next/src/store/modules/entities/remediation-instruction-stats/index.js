import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import { createEntityModule } from '@/store/plugins/entities';

export default createEntityModule({
  route: API_ROUTES.remediation.instructionStats,
  entityType: ENTITIES_TYPES.remediationInstructionStats,
  dataPreparer: d => d.data,
  withMeta: true,
});
