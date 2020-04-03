<template lang="pug">
  div
    template(v-if="step._t !== 'statecounter'")
      .header
        alarm-chips.chips.pr-2(
          v-if="!isStepTypeAction",
          :value="step.val",
          :type="stepType"
        )
        p {{ stepTitle }}
      .content
        p(v-if="isHTMLEnabled", v-html="step.m")
        p(v-else) {{ step.m }}
    template(v-else)
      .header
        p {{ $t('alarmList.timeLine.stateCounter.header') }}
      .content
        table
          tr
            td {{ $t('alarmList.timeLine.stateCounter.stateIncreased') }} :
            td {{ step.val.stateinc }}
          tr
            td {{ $t('alarmList.timeLine.stateCounter.stateDecreased') }} :
            td {{ step.val.statedec }}
          tr(v-for="state in states")
            td {{ $t('common.state') }} {{ state.text }} :
            td {{ state.value }}
</template>

<script>
import { ENTITIES_STATES_STYLES, ENTITY_INFOS_TYPE, ALARMS_LIST_TIME_LINE_SYSTEM_AUTHOR } from '@/constants';

import AlarmChips from '@/components/other/alarm/alarm-chips.vue';

export default {
  components: { AlarmChips },
  props: {
    step: {
      type: Object,
      required: true,
    },
    isHTMLEnabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    stepType() {
      const { _t: type } = this.step;

      if (type.startsWith('status')) {
        return ENTITY_INFOS_TYPE.status;
      } else if (type.startsWith('state')) {
        return ENTITY_INFOS_TYPE.state;
      }

      return ENTITY_INFOS_TYPE.action;
    },

    isStepTypeAction() {
      return this.stepType === ENTITY_INFOS_TYPE.action;
    },

    states() {
      const { val: states } = this.step;
      const prefix = 'state:';

      return Object.keys(states).reduce((acc, key) => {
        if (key.startsWith(prefix)) {
          const stateValue = parseInt(key.replace(prefix, ''), 10);

          acc.push({
            text: ENTITIES_STATES_STYLES[stateValue] && ENTITIES_STATES_STYLES[stateValue].text,
            value: states[key],
          });
        }

        return acc;
      }, []);
    },

    stepTitle() {
      const { _t: type, a: author, role } = this.step;
      const typeMessageKey = `alarmList.timeLine.types.${type}`;

      let formattedStepTitle = '';

      if (this.$te(typeMessageKey)) {
        formattedStepTitle = this.$t(typeMessageKey);
      } else if (this.isStepTypeAction) {
        formattedStepTitle = type.replace(/(declare)|(ack)/g, '$& ');
      } else {
        formattedStepTitle = type.replace(/(status)|(state)/g, '$& ').replace(/(inc)|(dec)/g, '$&reased');
      }

      formattedStepTitle += ` ${this.$t('alarmList.timeLine.titlePaths.by')} `;

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
</script>

<style lang="scss" scoped>
  $border_line: #DDDDE0;

  .content {
    padding-left: 20px;
    padding-top: 20px;
    overflow-wrap: break-word;
    word-break: break-all;
    width: 90%;
    max-height: 600px;
    overflow-y: auto;
  }

  .header {
    color: #686868;
    display: flex;
    align-items: baseline;
    font-weight: bold;
    border-bottom: solid 1px $border_line;
    padding-left: 5px;
    padding-top: 5px;


    .chips {
      font-size: 15px;
      height: 25px;
    }

    p {
      font-size: 15px;

      &:first-letter {
        text-transform: uppercase;
      }
    }
  }

  p {
    overflow-wrap: break-word;
    text-overflow: ellipsis;
    width: 90%;
  }
</style>
