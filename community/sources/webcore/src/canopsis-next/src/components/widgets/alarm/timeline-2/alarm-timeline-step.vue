<template>
  <v-layout :class="wrapperClass" class="timeline-step" align-center>
    <div class="timeline-step__time grey--text">
      {{ time }}
    </div>
    <v-flex class="timeline-step__main-content-wrapper">
      <v-layout class="timeline-step__main-content pb-2" column>
        <v-layout class="gap-2">
          <div class="timeline-step__icon">
            <v-icon size="20">
              $vuetify.icons.webhook
            </v-icon>
          </div>
          <c-expand-btn
            v-if="hasChildren"
            v-model="expanded"
            class="timeline-step__expand-button"
          />
          <div class="timeline-step__title text-subtitle-2">
            title
          </div>
          <div v-if="resultIcon" class="timeline-step__result-icon">
            <v-icon :color="resultIcon.color">
              {{ resultIcon.icon }}
            </v-icon>
          </div>
        </v-layout>
        <div class="grey--text mt-2 pre-wrap">
          message
        </div>
      </v-layout>
    </v-flex>
  </v-layout>
</template>

<script>
import { ALARM_LIST_STEPS, DATETIME_FORMATS } from '@/constants';

import { convertDateToString } from '@/helpers/date/date';

export default {
  props: {
    step: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      expanded: false,
    };
  },
  computed: {
    wrapperClass() {
      return {
        'timeline-step--pbehavior': this.step.in_pbh,
        'timeline-step--pbehavior-start': this.step._t === ALARM_LIST_STEPS.pbhenter,
        'timeline-step--pbehavior-end': this.step._t === ALARM_LIST_STEPS.pbhleave,
      };
    },

    time() {
      return convertDateToString(this.step.t, DATETIME_FORMATS.time);
    },

    hasChildren() {
      return !!this.step.steps?.length;
    },

    resultIcon() {
      return {
        icon: 'check_circle',
        color: 'success',
      };
      /* switch (this.step._t) {
        case ALARM_LIST_STEPS.ticketDeclarationRuleCompleted:
        case ALARM_LIST_STEPS.webhookComplete:
          return {
            icon: 'check_circle',
            color: 'success',
          };

        case ALARM_LIST_STEPS.ticketDeclarationRuleFailed:
        case ALARM_LIST_STEPS.webhookFail:
          return {
            icon: 'cancel',
            color: 'error',
          };

        default:
          return null;
      } */
    },
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

      .timeline-step--pbehavior-start &:before, .timeline-step--pbehavior &:before {
        content: '';
        position: absolute;
        width: 1px;
        left: 0;
        bottom: 0;
        background: $borderColor;
        height: 100%;
      }

      .timeline-step--pbehavior-start &:before {
        height: calc(100% - 30px);
      }
    }
  }

  &__icon {
    .timeline-step--pbehavior-start &, .timeline-step--pbehavior-end & {
      position: absolute;
      left: 0;
      transform: translateX(-50%);
    }
  }
}
</style>
