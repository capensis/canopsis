<template>
  <c-layout :class="wrapperClass" class="timeline-step" align="center">
    <div class="timeline-step__time text-medium-emphasis">
      {{ time }}
    </div>
    <div class="timeline-step__main-content-wrapper">
      <c-layout class="timeline-step__main-content pb-2" column>
        <c-layout gap="2">
          <div class="timeline-step__icon">
            <alarm-timeline-step-icon :step="step" />
          </div>
          <c-expand-btn
            v-if="hasSteps"
            v-model="expanded"
            size="x-small"
            icon-size="18"
            class="timeline-step__expand-button"
          />
          <div class="timeline-step__title">
            <alarm-timeline-step-title :step="step" :deep="deep" />
          </div>
          <div v-if="resultIcon" class="timeline-step__result-icon">
            <c-enabled :value="resultIcon.value" />
          </div>
        </c-layout>
        <div v-if="step.m" class="text-medium-emphasis mt-2 pre-wrap">
          <span v-if="isHtmlEnabled" v-html="sanitizedStepMessage" />
          <span v-else>
            {{ step.m }}
          </span>
        </div>
      </c-layout>
    </div>
  </c-layout>
</template>

<script>
import { computed, ref, watch } from 'vue';

import { ALARM_LIST_STEPS, DATETIME_FORMATS } from '@/constants';

import { convertDateToString } from '@/helpers/date/date';
import { linkifyHtml, sanitizeHtml } from '@/helpers/html';

import AlarmTimelineStepIcon from './alarm-timeline-step-icon.vue';
import AlarmTimelineStepTitle from './alarm-timeline-step-title.vue';

export default {
  components: { AlarmTimelineStepIcon, AlarmTimelineStepTitle },

  props: {
    step: {
      type: Object,
      default: () => ({}),
    },
    deep: {
      type: Boolean,
      default: false,
    },
    isHtmlEnabled: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['expand'],

  setup(props, { emit }) {
    const expanded = ref(false);

    const time = computed(() => convertDateToString(props.step.t, DATETIME_FORMATS.time));
    const hasSteps = computed(() => !!props.step.steps?.length);
    const wrapperClass = computed(() => ({
      'timeline-step--pbehavior': props.step.in_pbh,
      'timeline-step--pbehavior-enter': props.step._t === ALARM_LIST_STEPS.pbhenter,
      'timeline-step--pbehavior-leave': props.step._t === ALARM_LIST_STEPS.pbhleave,
      'timeline-step--deep': props.deep,
    }));

    const resultIcon = computed(() => {
      switch (props.step._t) {
        case ALARM_LIST_STEPS.declareTicketRuleComplete:
        case ALARM_LIST_STEPS.webhookComplete:
          return { value: true };

        case ALARM_LIST_STEPS.declareTicketRuleFail:
        case ALARM_LIST_STEPS.webhookFail:
          return { value: false };

        default:
          return null;
      }
    });

    const sanitizedStepMessage = computed(() => sanitizeHtml(linkifyHtml(String(props.step?.m ?? ''))),
    );

    watch(expanded, value => emit('expand', value));

    return {
      expanded,
      hasSteps,
      wrapperClass,
      time,
      resultIcon,
      sanitizedStepMessage,
    };
  },
};
</script>

<style lang="scss" scoped>
.timeline-step {
  --alarm-timeline-step-margins: 35px;
  --alarm-timeline-step-time-width: 56px;

  &__time {
    width: var(--alarm-timeline-step-time-width);
  }

  &__main-content {
    border-bottom: 1px solid rgba(var(--v-theme-on-surface), 0.12);

    &-wrapper {
      position: relative;
      flex: 1 1 auto;
      margin-left: var(--alarm-timeline-step-margins);
      padding: 0 var(--alarm-timeline-step-margins);

      .timeline-step--pbehavior-leave &:before,
      .timeline-step--pbehavior &:before {
        content: '';
        position: absolute;
        width: 1px;
        left: 0;
        bottom: 0;
        background: rgba(var(--v-theme-on-surface), 0.12);
        height: 100%;
      }

      .timeline-step--pbehavior-leave &:before {
        height: calc(100% - 30px);
      }

      .timeline-step--deep & {
        padding-left: calc(var(--alarm-timeline-step-margins) * 2);
      }
    }
  }

  &__icon {
    .timeline-step--pbehavior-leave &,
    .timeline-step--pbehavior-enter & {
      position: absolute;
      left: 0;
      transform: translateX(-50%);
    }
  }

  &__title {
    display: inline-flex;
    align-items: center;
  }
}
</style>
