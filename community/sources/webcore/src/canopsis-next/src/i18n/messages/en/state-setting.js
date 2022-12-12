import { STATE_SETTING_METHODS } from '@/constants';

export default {
  worstLabel: 'The worst of:',
  worstHelpText: 'Canopsis counts the state for each criterion defined. The final state of JUnit test suite is taken as a worst of resulting states.',
  criterion: 'Criterion',
  serviceState: 'Service state',
  methods: {
    [STATE_SETTING_METHODS.worst]: 'Worst',
    [STATE_SETTING_METHODS.worstOfShare]: 'Worst of share',
  },
  states: {
    minor: 'Minor',
    major: 'Major',
    critical: 'Critical',
  },
};
