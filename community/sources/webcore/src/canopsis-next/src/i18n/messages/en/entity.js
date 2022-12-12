import { ENTITY_TYPES } from '@/constants';

export default {
  manageInfos: 'Manage Infos',
  form: 'Form',
  impact: 'Impact',
  depends: 'Depends',
  addInformation: 'Add Information',
  emptyInfos: 'No information',
  availabilityState: 'Hi availability state',
  okEvents: 'OK events',
  koEvents: 'KO events',
  types: {
    [ENTITY_TYPES.component]: 'Component',
    [ENTITY_TYPES.connector]: 'Connector',
    [ENTITY_TYPES.resource]: 'Resource',
    [ENTITY_TYPES.service]: 'Service',
  },
};
