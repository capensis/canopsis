import { STATE_SETTING_METHODS } from '@/constants';

export default {
  worstLabel: 'Le pire critère :',
  worstHelpText: 'Canopsis compte l\'état pour chaque critère défini. L\'état final de la suite de tests JUnit est considéré comme le pire des états résultants.',
  criterion: 'Critère',
  serviceState: 'État du service',
  methods: {
    [STATE_SETTING_METHODS.worst]: 'Pire',
    [STATE_SETTING_METHODS.worstOfShare]: 'Pire des états',
  },
  states: {
    minor: 'Mineur',
    major: 'Majeur',
    critical: 'Critique',
  },
};
