import {
  ALARM_TEMPLATE_FIELDS,
  ALARM_FIELDS_TO_LABELS_KEYS,
  ENTITY_TEMPLATE_FIELDS,
  PBEHAVIOR_INFO_FIELDS,
  ALARM_EXPORT_PDF_FIELDS,
  ALARM_EXPORT_PDF_FIELDS_TO_ORIGINAL_FIELDS,
} from '@/constants';

import { variablesMixin } from './common';

export const alarmVariablesMixin = {
  mixins: [variablesMixin],
  computed: {
    ticketVariables() {
      return [this.stepValueVariable, this.stepAuthorVariable];
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

    exportPdfAlarmVariables() {
      return [
        {
          value: ALARM_EXPORT_PDF_FIELDS.currentDate,
          text: this.$t('common.currentDate'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.displayName,
          text: this.$t('alarm.fields.displayName'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.state,
          text: this.$t('common.state'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.status,
          text: this.$t('common.status'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.connector,
          text: this.$t('common.connector'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.connectorName,
          text: this.$t('common.connectorName'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.component,
          text: this.$t('common.component'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.resource,
          text: this.$t('common.resource'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.creationDate,
          text: this.$t('common.created'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.duration,
          text: this.$t('common.duration'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.infos,
          text: this.$t('common.infos'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.pbehaviorInfo,
          text: this.$t('pbehavior.pbehaviorInfo'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.initialOutput,
          text: this.$t('alarm.fields.initialOutput'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.output,
          text: this.$t('common.output'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.lastEventDate,
          text: this.$t('common.lastEventDate'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.lastUpdateDate,
          text: this.$t('common.updated'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.acknowledgeDate,
          text: this.$t('alarm.fields.ackAt'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.resolved,
          text: this.$t('alarm.fields.resolved'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.eventsCount,
          text: this.$t('alarm.fields.eventsCount'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.ticket,
          text: this.$tc('common.ticket'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.comments,
          text: this.$tc('common.comment', 2),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.activationDate,
          text: this.$t('alarm.fields.activationDate'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.tags,
          text: this.$tc('common.tag', 2),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.links,
          text: this.$tc('common.link', 2),
        },
      ].map(variable => ({ ...variable, value: ALARM_EXPORT_PDF_FIELDS_TO_ORIGINAL_FIELDS[variable.value] }));
    },
  },
};
