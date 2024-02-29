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
  entitiesStates: 'Dependencies state',
  computeMethod: 'State compute method',
  dependenciesEntityPattern: 'Pattern for dependencies',
  conditionsError: 'Please add at least one condition',
  targetEntityThresholdSummary: 'A targeted entity state is {state} when the {method} of dependencies of the {dependenciesEntitiesState} state is {condition} {value}.',
  entityThresholdSummary: '{name} state is {state} when the {method} of dependencies of the {dependenciesEntitiesState} state is {condition} {value}.',
  appliedFor: 'Applied for',
  appliedForEntityType: 'Applied for entity type',
  stateIsInheritFrom: '{name} state is inherit from',
  seeFilterPattern: 'See filter pattern',
  dependsCount: 'Total number of dependencies',
  stateDependsCount: 'Number of dependencies of the {state} state',
  steps: {
    basics: 'Basics',
    rulePatterns: 'Define target entities',
    conditions: 'Add conditions',
  },
  methods: {
    [STATE_SETTING_METHODS.inherited]: {
      label: 'State is inherited from dependencies',
      tooltip: 'State is defined by the one or several dependencies.\n'
          + 'When several resources are defined, the worst state of them is taken.',
      stepTitle: 'The target entity state is inherited from one or several dependencies. When several dependencies fit the pattern, the worst state is taken.',
    },
    [STATE_SETTING_METHODS.dependencies]: {
      label: 'State is defined by a share or number of dependencies of a specific state',
      tooltip: 'Entity states can be overridden with custom rule defined by number or share of dependencies of specific states. ',
      stepTitle: 'The target entity states can be overridden by conditions based on a number or share of dependencies of a specific states.',
    },
  },
  thresholdMethods: {
    [STATE_SETTING_THRESHOLDS_METHODS.share]: 'Share',
    [STATE_SETTING_THRESHOLDS_METHODS.number]: 'Number',
  },
  thresholdConditions: {
    [STATE_SETTING_THRESHOLDS_CONDITIONS.greater]: 'Greater than',
    [STATE_SETTING_THRESHOLDS_CONDITIONS.less]: 'Less than',
  },
  states: { // TODO: remove
    [ENTITIES_STATES_KEYS.ok]: 'OK',
    [ENTITIES_STATES_KEYS.minor]: 'Minor',
    [ENTITIES_STATES_KEYS.major]: 'Major',
    [ENTITIES_STATES_KEYS.critical]: 'Critical',
  },
  junit: {
    worstLabel: 'The worst of:',
    worstHelpText: 'Canopsis counts the state for each criterion defined. The final state of JUnit test suite is taken as a worst of resulting states.',
    criterion: 'Criterion',
    methods: {
      [JUNIT_STATE_SETTING_METHODS.worst]: 'Worst',
      [JUNIT_STATE_SETTING_METHODS.worstOfShare]: 'Worst of share',
    },
  },
};
