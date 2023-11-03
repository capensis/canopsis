import {
  ENTITIES_STATES_KEYS,
  STATE_SETTING_METHODS,
  STATE_SETTING_THRESHOLDS_METHODS,
  STATE_SETTING_THRESHOLDS_CONDITIONS,
  JUNIT_STATE_SETTING_METHODS,
} from '@/constants';

export default {
  title: 'State setting',
  targetEntityState: 'Target entity state',
  entitiesStates: 'Impacting entities states',
  computeMethod: 'State compute method',
  addImpactingEntityPattern: 'Impacting entity pattern',
  conditionsError: 'Please add at least one condition',
  steps: {
    basics: 'Basics',
    rulePatterns: 'Define target entities',
    conditions: 'Add conditions',
  },
  methods: {
    [STATE_SETTING_METHODS.inherited]: {
      label: 'State is inherited from impacting entities',
      tooltip: 'State is defined by the one or several impacting resources.\n'
          + 'When several resources are defined, the worst state of them is taken.',
    },
    [STATE_SETTING_METHODS.dependencies]: {
      label: 'State is defined by a share or number of impacting entities of a specific state',
      tooltip: 'Entity states can be overridden with custom rule defined by number or share of impacting entities of specific states.',
    },
  },
  calculationMethods: {
    [STATE_SETTING_THRESHOLDS_METHODS.share]: 'Share',
    [STATE_SETTING_THRESHOLDS_METHODS.number]: 'Number',
  },
  thresholdConditions: {
    [STATE_SETTING_THRESHOLDS_CONDITIONS.greater]: 'Greater than',
    [STATE_SETTING_THRESHOLDS_CONDITIONS.less]: 'Less than',
  },
  states: {
    [ENTITIES_STATES_KEYS.ok]: 'OK',
    [ENTITIES_STATES_KEYS.minor]: 'Minor',
    [ENTITIES_STATES_KEYS.major]: 'Major',
    [ENTITIES_STATES_KEYS.critical]: 'Critical',
  },
  junit: {
    worstLabel: 'The worst of:',
    worstHelpText: 'Canopsis counts the state for each criterion defined. The final state of JUnit test suite is taken as a worst of resulting states.',
    criterion: 'Criterion',
    serviceState: 'Service state',
    methods: {
      [JUNIT_STATE_SETTING_METHODS.worst]: 'Worst',
      [JUNIT_STATE_SETTING_METHODS.worstOfShare]: 'Worst of share',
    },
  },
};
