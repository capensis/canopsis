<template lang="pug">
  div
    template(v-if="step._t !== 'statecounter'")
      div.header
        c-alarm-chip.chips.pr-2(
          v-if="!isStepTypeAction",
          :value="step.val",
          :type="stepType"
        )
        p {{ stepTitle }}
      div.content
        p(v-if="isHtmlEnabled", v-html="step.m")
        p(v-else) {{ step.m }}
    template(v-else)
      div.header
        p {{ $t('alarmList.timeLine.stateCounter.header') }}
      div.content
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
import { widgetExpandPanelAlarmTimelineCard } from '@/mixins/widget/expand-panel/alarm/timeline-card';

export default {
  mixins: [widgetExpandPanelAlarmTimelineCard],
  props: {
    step: {
      type: Object,
      required: true,
    },
    isHtmlEnabled: {
      type: Boolean,
      default: false,
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
    white-space: pre-line;
    overflow-wrap: break-word;
    text-overflow: ellipsis;
    width: 90%;
  }
</style>
