<template>
  <c-tooltip location="top">
    <template #activator="{ props: activatorProps }">
      <c-layout
        v-bind="activatorProps"
        :class="{ 'event-statistics__tooltip--inactive': hasInactivePbehavior }"
        class="event-statistics__tooltip"
        justify="center"
      >
        <span class="mr-1 text-success font-weight-bold">
          {{ entity.ok_events }}
        </span>
        <span>/</span>
        <span class="ml-1 text-error font-weight-bold">
          {{ entity.ko_events }}
        </span>
      </c-layout>
    </template>
    <span class="pre-wrap">
      {{ statisticsMessage }}
    </span>
  </c-tooltip>
</template>

<script>
import { PBEHAVIOR_TYPE_TYPES } from '@/constants';

export default {
  props: {
    entity: {
      type: Object,
      required: true,
    },
  },

  computed: {
    statisticsMessage() {
      return this.$t('context.eventStatisticsMessage', {
        ok: this.entity.ok_events,
        ko: this.entity.ko_events,
      });
    },

    hasInactivePbehavior() {
      return this.entity?.pbehavior_info?.canonical_type === PBEHAVIOR_TYPE_TYPES.inactive;
    },
  },
};
</script>

<style lang="scss">
.event-statistics {
  &__tooltip {
    &--inactive {
      opacity: 0.5;
    }
  }
}
</style>
