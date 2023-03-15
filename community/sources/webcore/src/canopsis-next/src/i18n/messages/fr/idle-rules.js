import { IDLE_RULE_ALARM_CONDITIONS, IDLE_RULE_TYPES } from '@/constants';

export default {
  timeAwaiting: 'Temps d\'attente',
  timeRangeAwaiting: 'Plage de temps en attente',
  types: {
    [IDLE_RULE_TYPES.alarm]: 'Règle d\'inactivité d\'alarme',
    [IDLE_RULE_TYPES.entity]: 'Règle d\'inactivité d\'entité',
  },
  alarmConditions: {
    [IDLE_RULE_ALARM_CONDITIONS.lastEvent]: 'Aucun événement reçu',
    [IDLE_RULE_ALARM_CONDITIONS.lastUpdate]: 'Aucun changement d\'état',
  },
};
