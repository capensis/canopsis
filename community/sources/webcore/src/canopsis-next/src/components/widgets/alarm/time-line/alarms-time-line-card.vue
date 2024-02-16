<template>
  <div class="time-line-card">
    <div class="time-line-card__header text--secondary">
      <template v-if="isNotStateCounter">
        <c-alarm-chip
          class="chips pr-2"
          v-if="!isStepTypeAction"
          :value="step.val"
          :type="stepType"
        />
        <p>{{ stepTitle }}</p>
      </template>
      <p v-else>
        {{ $t('alarm.timeLine.stateCounter.header') }}
      </p>
    </div>
    <div class="time-line-card__content text--disabled">
      <template v-if="isNotStateCounter">
        <div
          v-if="isHtmlEnabled"
          v-html="sanitizedStepMessage"
        />
        <p v-else>
          {{ step.m }}
        </p>
      </template>
      <table v-else>
        <tr>
          <td>{{ $t('alarm.timeLine.stateCounter.stateIncreased') }} :</td>
          <td>{{ step.val.stateinc }}</td>
        </tr>
        <tr>
          <td>{{ $t('alarm.timeLine.stateCounter.stateDecreased') }} :</td>
          <td>{{ step.val.statedec }}</td>
        </tr>
        <tr
          v-for="state in states"
          :key="state.value"
        >
          <td>{{ $t('common.state') }} {{ state.text }} :</td>
          <td>{{ state.value }}</td>
        </tr>
      </table>
    </div>
  </div>
</template>

<script>
import { sanitizeHtml, linkifyHtml } from '@/helpers/html';

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
  computed: {
    isNotStateCounter() {
      return this.step._t !== 'statecounter';
    },

    sanitizedStepMessage() {
      return sanitizeHtml(linkifyHtml(String(this.step?.m ?? '')));
    },
  },
};
</script>

<style lang="scss" scoped>
$border_line: #DDDDE0;

.time-line-card {
  &__content {
    padding-left: 20px;
    padding-top: 20px;
    overflow-wrap: break-word;
    word-break: break-all;
    width: 90%;
    max-height: 600px;
    overflow-y: auto;
  }

  &__header {
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
}
</style>
