import { LINK_RULE_TYPES } from '@/constants';

export default {
  simpleMode: 'Simple mode',
  advancedMode: 'Advanced mode',
  addLink: 'Add link',
  linksEmpty: 'No link added yet',
  linksEmptyError: 'You should add at least 1 link in simple mode or edit source code in advanced mode',
  sourceCodeAlert: 'Please, change this script only if you are completely aware of what you are doing',
  type: 'Link type',
  single: 'Apply this link only to single alarm ?',
  types: {
    [LINK_RULE_TYPES.alarm]: 'Alarm',
    [LINK_RULE_TYPES.entity]: 'Entity',
  },
};
