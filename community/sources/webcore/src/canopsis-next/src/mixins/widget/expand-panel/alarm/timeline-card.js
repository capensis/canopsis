import { ALARMS_LIST_TIME_LINE_SYSTEM_AUTHOR, ENTITY_INFOS_TYPE } from '@/constants';

export const widgetExpandPanelAlarmTimelineCard = {
  props: {
    step: {
      type: Object,
      required: true,
    },
  },
  computed: {
    stepType() {
      const { _t: type } = this.step;

      if (type.startsWith('status')) {
        return ENTITY_INFOS_TYPE.status;
      } if (type.startsWith('state')) {
        return ENTITY_INFOS_TYPE.state;
      }

      return ENTITY_INFOS_TYPE.action;
    },

    isStepTypeAction() {
      return this.stepType === ENTITY_INFOS_TYPE.action;
    },

    isStepTypeState() {
      return this.stepType === ENTITY_INFOS_TYPE.state;
    },

    states() {
      const { val: states } = this.step;
      const prefix = 'state:';

      return Object.keys(states).reduce((acc, key) => {
        if (key.startsWith(prefix)) {
          const stateValue = parseInt(key.replace(prefix, ''), 10);

          acc.push({
            text: stateValue,
            value: states[key],
          });
        }

        return acc;
      }, []);
    },

    stepTitle() {
      const { _t: type, a: author, role } = this.step;
      const typeMessageKey = `alarm.timeLine.types.${type}`;

      let formattedStepTitle = '';

      if (this.$te(typeMessageKey)) {
        formattedStepTitle = this.$t(typeMessageKey);
      } else if (this.isStepTypeAction) {
        formattedStepTitle = type.replace(/(declare)|(ack)/g, '$& ');
      } else {
        formattedStepTitle = type.replace(/(status)|(state)/g, '$& ').replace(/(inc)|(dec)/g, '$&reased');
      }

      formattedStepTitle += ` ${this.$t('alarm.timeLine.titlePaths.by')} `;

      if (author === ALARMS_LIST_TIME_LINE_SYSTEM_AUTHOR) {
        formattedStepTitle += 'system';
      } else {
        formattedStepTitle += author;
      }

      if (role) {
        formattedStepTitle += ` (${role})`;
      }

      return formattedStepTitle;
    },
  },
};
