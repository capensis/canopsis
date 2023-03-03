import { ENTITY_TYPES } from '@/constants';

export default {
  manageInfos: 'Manage Infos',
  form: 'Form',
  impact: 'Impact',
  depends: 'Depends',
  addInformation: 'Add Information',
  emptyInfos: 'No information',
  availabilityState: 'Hi availability state',
  types: {
    [ENTITY_TYPES.component]: 'Component',
    [ENTITY_TYPES.connector]: 'Connector',
    [ENTITY_TYPES.resource]: 'Resource',
    [ENTITY_TYPES.service]: 'Service',
  },
  fields: {
    categoryName: 'Category name',
    koEvents: 'KO events',
    okEvents: 'OK events',
    statsKo: 'Stats KO',
    statsOk: 'Stats OK',
    idleSince: 'Idle since',
    componentInfos: 'Component infos',
    alarmDisplayName: 'Alarm display name',
    alarmCreationDate: 'Alarm creation date',
  },
};
