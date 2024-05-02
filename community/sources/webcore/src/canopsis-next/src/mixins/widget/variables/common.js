import { VARIABLES_STEP_FIELDS, PBEHAVIOR_INFO_FIELDS } from '@/constants';

export const variablesMixin = {
  computed: {
    stepValueVariable() {
      return {
        text: this.$t('common.value'),
        value: VARIABLES_STEP_FIELDS.value,
      };
    },

    stepTimestampVariable() {
      return {
        text: this.$t('common.timestamp'),
        value: VARIABLES_STEP_FIELDS.timestamp,
      };
    },

    stepMessageVariable() {
      return {
        text: this.$t('common.message'),
        value: VARIABLES_STEP_FIELDS.message,
      };
    },

    stepAuthorVariable() {
      return {
        text: this.$t('common.author'),
        value: VARIABLES_STEP_FIELDS.author,
      };
    },

    stateVariables() {
      return [this.stepValueVariable];
    },

    statusVariables() {
      return [this.stepValueVariable];
    },

    ackVariables() {
      return [
        this.stepTimestampVariable,
        this.stepAuthorVariable,
        this.stepMessageVariable,
      ];
    },

    ticketVariables() {
      return [this.stepValueVariable];
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
