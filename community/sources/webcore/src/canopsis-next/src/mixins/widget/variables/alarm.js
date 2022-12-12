import { ALARM_TEMPLATE_FIELDS, ENTITY_TEMPLATE_FIELDS, PBEHAVIOR_INFO_FIELDS } from '@/constants';

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
        {
          text: this.$t('common.id'),
          value: ENTITY_TEMPLATE_FIELDS.id,
        },
        {
          text: this.$t('common.ack'),
          value: ALARM_TEMPLATE_FIELDS.ack,
          variables: this.ackVariables,
        },
        {
          text: this.$t('common.state'),
          value: ALARM_TEMPLATE_FIELDS.state,
          variables: this.stateVariables,
        },
        {
          text: this.$t('common.status'),
          value: ALARM_TEMPLATE_FIELDS.status,
          variables: this.statusVariables,
        },
        {
          text: this.$t('common.ticket'),
          value: ALARM_TEMPLATE_FIELDS.ticket,
          variables: this.ticketVariables,
        },
        {
          text: this.$t('common.component'),
          value: ALARM_TEMPLATE_FIELDS.component,
        },
        {
          text: this.$t('common.connector'),
          value: ALARM_TEMPLATE_FIELDS.connector,
        },
        {
          text: this.$t('common.connectorName'),
          value: ALARM_TEMPLATE_FIELDS.connectorName,
        },
        {
          text: this.$t('common.resource'),
          value: ALARM_TEMPLATE_FIELDS.resource,
        },
        {
          text: this.$t('common.created'),
          value: ALARM_TEMPLATE_FIELDS.creationDate,
        },
        {
          text: this.$t('common.displayName'),
          value: ALARM_TEMPLATE_FIELDS.displayName,
        },
        {
          text: this.$t('common.output'),
          value: ALARM_TEMPLATE_FIELDS.output,
        },
        {
          text: this.$t('common.updated'),
          value: ALARM_TEMPLATE_FIELDS.lastUpdateDate,
        },
        {
          text: this.$t('common.lastEventDate'),
          value: ALARM_TEMPLATE_FIELDS.lastEventDate,
        },
        {
          text: this.$t('pbehavior.pbehaviorInfo'),
          value: ALARM_TEMPLATE_FIELDS.pbehaviorInfo,
          variables: this.pbehaviorInfoVariables,
        },
        {
          text: this.$t('common.duration'),
          value: ALARM_TEMPLATE_FIELDS.duration,
        },
        {
          text: this.$t('alarm.eventsCount'),
          value: ALARM_TEMPLATE_FIELDS.eventsCount,
        },
      ];
    },
  },
};
