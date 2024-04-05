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
  searchHelp: '<span>Help on the advanced research :</span>\n'
    + '<p>- [ NOT ] &lt;ColumnName&gt; &lt;Operator&gt; &lt;Value&gt;</p> [ AND|OR [ NOT ] &lt;ColumnName&gt; &lt;Operator&gt; &lt;Value&gt; ]\n'
    + '<p>The "-" before the research is required</p>\n'
    + '<p>Operators : <=, <,=, !=,>=, >, LIKE (For MongoDB regular expression)</p>\n'
    + '<p>For querying patterns, use "pattern" keyword as the &lt;ColumnName&gt; alias</p>\n'
    + '<p>Value\'s type : String between quote, Boolean ("TRUE", "FALSE"), Integer, Float, "NULL"</p>\n'
    + '<dl>'
    + '  <dt>Examples :</dt>'
    + '  <dt>- name = "name_1"</dt>\n'
    + '  <dd>Pbehavior name are "name_1"</dd>\n'
    + '  <dt>- rrule = "rrule_1"</dt>\n'
    + '  <dd>Pbehavior rrule are "rrule_1"</dd>\n'
    + '  <dt>- filter = "filter_1"</dt>\n'
    + '  <dd>Pbehavior filter are "filter_1"</dd>\n'
    + '  <dt>- type.name = "type_name_1"</dt>\n'
    + '  <dd>Pbehavior type name are "type_name_1"</dd>\n'
    + '  <dt>- reason.name = "reason_name_1"</dt>\n'
    + '  <dd>Pbehavior reason name are "reason_name_1"</dd>'
    + '</dl>',
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
