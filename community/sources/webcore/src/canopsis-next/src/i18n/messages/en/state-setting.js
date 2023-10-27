import {
  STATE_SETTING_COMPUTE_METHODS,
  STATE_SETTING_CALCULATION_METHODS,
  STATE_SETTING_CONDITIONS,
} from '@/constants';

export default {
  worstLabel: 'The worst of:',
  worstHelpText: 'Canopsis counts the state for each criterion defined. The final state of JUnit test suite is taken as a worst of resulting states.',
  criterion: 'Criterion',
  serviceState: 'Service state',
  targetEntityState: 'Target entity state',
  computeMethod: 'State compute method',
  computeMethods: {
    [STATE_SETTING_COMPUTE_METHODS.inherit]: {
      label: 'State is inherited from impacting entities',
      tooltip: 'State is defined by the one or several impacting resources.\n'
          + 'When several resources are defined, the worst state of them is taken.',
    },
    [STATE_SETTING_COMPUTE_METHODS.shareOfDependencies]: {
      label: 'State is defined by a share or number of impacting entities of a specific state',
      tooltip: 'Entity states can be overridden with custom rule defined by number or share of impacting entities of specific states.',
    },
  },
  calculationMethods: {
    [STATE_SETTING_CALCULATION_METHODS.share]: 'Share',
    [STATE_SETTING_CALCULATION_METHODS.number]: 'Number',
  },
  conditions: {
    [STATE_SETTING_CONDITIONS.greater]: 'Greater than',
    [STATE_SETTING_CONDITIONS.less]: 'Less than',
  },
  entitiesStates: 'Impacting entities states',
  states: {
    ok: 'ok',
    minor: 'Minor',
    major: 'Major',
    critical: 'Critical',
  },
};
