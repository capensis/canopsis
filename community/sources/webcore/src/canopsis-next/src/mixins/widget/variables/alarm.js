import {
  ALARM_TEMPLATE_FIELDS,
  ALARM_FIELDS_TO_LABELS_KEYS,
  ENTITY_TEMPLATE_FIELDS,
  PBEHAVIOR_INFO_FIELDS,
} from '@/constants';

import { variablesMixin } from './common';

export const alarmVariablesMixin = {
  mixins: [variablesMixin],
  computed: {
    ticketVariables() {
      return [this.alarmStepValueVariable, this.alarmStepAuthorVariable];
    },

    pbehaviorInfoVariables() {
      return [
        {
          text: this.$t('pbehavior.pbehaviorType'),
          value: PBEHAVIOR_INFO_FIELDS.typeName,
        },
        {
          text: this.$tc('pbehavior.pbehaviorReason'),
          value: PBEHAVIOR_INFO_FIELDS.reason,
        },
        {
          text: this.$t('pbehavior.pbehaviorName'),
          value: PBEHAVIOR_INFO_FIELDS.name,
        },
        {
          text: this.$t('pbehavior.pbehaviorCanonicalType'),
          value: PBEHAVIOR_INFO_FIELDS.canonicalType,
        },
      ];
    },

    alarmVariables() {
      return [
        { value: ENTITY_TEMPLATE_FIELDS.id },
        {
          value: ALARM_TEMPLATE_FIELDS.ack,
          variables: this.ackVariables,
        },
        {
          value: ALARM_TEMPLATE_FIELDS.state,
          variables: this.stateVariables,
        },
        {
          value: ALARM_TEMPLATE_FIELDS.status,
          variables: this.statusVariables,
        },
        {
          value: ALARM_TEMPLATE_FIELDS.ticket,
          variables: this.ticketVariables,
        },
        { value: ALARM_TEMPLATE_FIELDS.component },
        { value: ALARM_TEMPLATE_FIELDS.connector },
        { value: ALARM_TEMPLATE_FIELDS.connectorName },
        { value: ALARM_TEMPLATE_FIELDS.resource },
        { value: ALARM_TEMPLATE_FIELDS.creationDate },
        { value: ALARM_TEMPLATE_FIELDS.displayName },
        { value: ALARM_TEMPLATE_FIELDS.output },
        { value: ALARM_TEMPLATE_FIELDS.lastUpdateDate },
        { value: ALARM_TEMPLATE_FIELDS.lastEventDate },
        {
          value: ALARM_TEMPLATE_FIELDS.pbehaviorInfo,
          variables: this.pbehaviorInfoVariables,
        },
        { value: ALARM_TEMPLATE_FIELDS.duration },
        { value: ALARM_TEMPLATE_FIELDS.eventsCount },
      ].map(variable => ({
        ...variable,

        text: this.$tc(ALARM_FIELDS_TO_LABELS_KEYS[variable.value.replace('alarm.', '')], 2),
      }));
    },
  },
};
