import { PBEHAVIOR_TYPE_TYPES } from '@/constants';

export default {
  isEnabled: 'Is enabled',
  begins: 'Begins',
  ends: 'Ends',
  lastAlarmDate: 'Last alarm date',
  alarmCount: 'Alarms count',
  massRemove: 'Remove pbehaviors',
  massEnable: 'Enable pbehaviors',
  massDisable: 'Disable pbehaviors',
  periodsCalendar: 'Calendar with periods',
  notEditable: 'Cannot be modified',
  pbehaviorInfo: 'Pbehavior info',
  pbehaviorType: 'Pbehavior type',
  pbehaviorReason: 'Pbehavior reason',
  pbehaviorName: 'Pbehavior name',
  pbehaviorCanonicalType: 'Pbehavior canonical type',
  rruleEnd: 'End of recurrence',
  buttons: {
    addFilter: 'Add filter',
    editFilter: 'Edit filter',
    addRRule: 'Add recurrence rule',
    editRrule: 'Edit recurrence rule',
  },
  tabs: {
    type: 'Type',
    reason: 'Reason',
    exceptions: 'Exception dates',
  },

  exceptions: {
    title: 'Exception dates',
    create: 'Add an exception date',
    choose: 'Choose list of exceptions',
    usingException: 'Cannot be deleted since it is in use',
    emptyExceptions: 'No exceptions added yet',
  },

  types: {
    usingType: 'Cannot be deleted since it is in use',
    defaultType: 'The type is default, you can edit only color field',
    hidden: 'Hide this type on pbehavior form ?',
    types: {
      [PBEHAVIOR_TYPE_TYPES.active]: 'Active',
      [PBEHAVIOR_TYPE_TYPES.inactive]: 'Inactive',
      [PBEHAVIOR_TYPE_TYPES.pause]: 'Pause',
      [PBEHAVIOR_TYPE_TYPES.maintenance]: 'Maintenance',
    },
  },

  reasons: {
    usingReason: 'Cannot be deleted since it is in use',
    hidden: 'Hide this reason on pbehavior form ?',
  },
};
