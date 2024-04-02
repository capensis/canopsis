<template>
  <v-layout :class="wrapperClass" class="timeline-step" align-center>
    <div class="timeline-step__time grey--text">
      {{ time }}
    </div>
    <v-flex class="timeline-step__main-content-wrapper">
      <v-layout class="timeline-step__main-content pb-2" column>
        <v-layout class="gap-2">
          <div class="timeline-step__icon">
            <alarm-timeline-step-icon :step="step" />
          </div>
          <c-expand-btn
            v-if="hasSteps"
            :expanded="expanded"
            class="timeline-step__expand-button"
            @expand="expand"
          />
          <div class="timeline-step__title">
            <alarm-timeline-step-title :step="step" :deep="deep" />
          </div>
          <div v-if="resultIcon" class="timeline-step__result-icon">
            <v-icon :color="resultIcon.color">
              {{ resultIcon.icon }}
            </v-icon>
          </div>
        </v-layout>
        <div v-if="step.m" class="grey--text mt-2 pre-wrap">
          {{ step.m }}
        </div>
      </v-layout>
    </v-flex>
  </v-layout>
</template>

<script>
import { ref, computed } from 'vue';

import { ALARM_LIST_STEPS, DATETIME_FORMATS } from '@/constants';

import { convertDateToString } from '@/helpers/date/date';

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
  },
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
          return {
            icon: 'check_circle',
            color: 'success',
          };

        case ALARM_LIST_STEPS.declareTicketRuleFailed:
        case ALARM_LIST_STEPS.webhookFail:
          return {
            icon: 'cancel',
            color: 'error',
          };

        default:
          return null;
      }
    });

    const expand = () => {
      expanded.value = !expanded.value;

      emit('expand', expanded.value);
    };

    return {
      expanded,

      hasSteps,
      wrapperClass,
      time,
      resultIcon,

      expand,
    };
  },
};
</script>

<style lang="scss" scoped>
$margins: 35px;
$borderColor: #cecece;

.timeline-step {
  &__main-content {
    border-bottom: 1px solid $borderColor;

    &-wrapper {
      position: relative;
      margin-left: $margins;
      padding: 0 $margins;

      .timeline-step--pbehavior-leave &:before, .timeline-step--pbehavior &:before {
        content: '';
        position: absolute;
        width: 1px;
        left: 0;
        bottom: 0;
        background: $borderColor;
        height: 100%;
      }

      .timeline-step--pbehavior-leave &:before {
        height: calc(100% - 30px);
      }

      .timeline-step--deep & {
        padding-left: $margins * 2;
      }
    }
  }

  &__icon {
    .timeline-step--pbehavior-leave &, .timeline-step--pbehavior-enter & {
      position: absolute;
      left: 0;
      transform: translateX(-50%);
    }
  }
}
</style>
