<template>
  <v-tooltip top>
    <template #activator="{ on }">
      <v-layout
        :class="{ 'event-statistics__tooltip--inactive': hasInactivePbehavior }"
        class="event-statistics__tooltip"
        justify-center
        v-on="on"
      >
        <span class="mr-1 success--text font-weight-bold">{{ entity.ok_events }}</span>
        <span>/</span>
        <span class="ml-1 error--text font-weight-bold">{{ entity.ko_events }}</span>
      </v-layout>
    </template>
    <span class="pre-wrap">{{ statisticsMessage }}</span>
  </v-tooltip>
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
