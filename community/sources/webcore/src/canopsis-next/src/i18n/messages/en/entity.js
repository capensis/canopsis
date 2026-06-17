import { ENTITY_TYPES, TREE_OF_DEPENDENCIES_SHOW_TYPES } from '@/constants';

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
    lastAlarmUpdateDate: 'Last alarm update date',
    alarmLastComment: 'Alarm last comment',
    importSource: 'Import source',
    imported: 'Imported date',
  },
  treeOfDependenciesShowTypes: {
    [TREE_OF_DEPENDENCIES_SHOW_TYPES.allDependencies]: 'Show all dependencies',
    [TREE_OF_DEPENDENCIES_SHOW_TYPES.dependenciesDefiningTheState]: 'Show dependencies defining the state',
    [TREE_OF_DEPENDENCIES_SHOW_TYPES.custom]: 'Show selector',
  },
  comments: {
    emptyList: 'No comments are added yet',
  },
  infosLog: {
    eventFilterId: 'Event filter ID',
    eventFilterDescription: 'Event filter description',
    oldValue: 'Old value',
    newValue: 'New value',
  },
};
