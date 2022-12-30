import {
  ALARM_STEP_FIELDS,
  PBEHAVIOR_INFO_FIELDS,
} from '@/constants';

export const variablesMixin = {
  computed: {
    alarmStepValueVariable() {
      return {
        text: this.$t('common.value'),
        value: ALARM_STEP_FIELDS.value,
      };
    },

    alarmStepTimestampVariable() {
      return {
        text: this.$t('common.timestamp'),
        value: ALARM_STEP_FIELDS.timestamp,
      };
    },

    alarmStepMessageVariable() {
      return {
        text: this.$t('common.message'),
        value: ALARM_STEP_FIELDS.message,
      };
    },

    alarmStepAuthorVariable() {
      return {
        text: this.$t('common.author'),
        value: ALARM_STEP_FIELDS.author,
      };
    },

    stateVariables() {
      return [this.alarmStepValueVariable];
    },

    statusVariables() {
      return [this.alarmStepValueVariable];
    },

    ackVariables() {
      return [
        this.alarmStepTimestampVariable,
        this.alarmStepAuthorVariable,
        this.alarmStepMessageVariable,
      ];
    },

    ticketVariables() {
      return [this.alarmStepValueVariable];
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
      ];
    },
  },
};
