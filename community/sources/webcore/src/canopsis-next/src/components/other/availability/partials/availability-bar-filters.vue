<template>
  <v-layout class="availability-bar-filters gap-6" wrap>
    <c-quick-date-interval-field
      :interval="interval"
      :min="minIntervalDate"
      :quick-ranges="quickRanges"
      class="availability-bar-filters__interval"
      short
      @input="$emit('update:interval', $event)"
    />
    <c-availability-show-type-field
      :value="showType"
      :label="$t('common.show')"
      class="availability-bar-filters__show-type"
      @input="$emit('update:showType', $event)"
    />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { AVAILABILITY_QUICK_RANGES } from '@/constants';

export default {
  props: {
    interval: {
      type: Object,
      required: false,
    },
    showType: {
      type: Number,
      required: false,
    },
    minIntervalDate: {
      type: Number,
      required: false,
    },
  },
  setup() {
    const quickRanges = computed(() => Object.values(AVAILABILITY_QUICK_RANGES));

    return {
      quickRanges,
    };
  },
};
</script>

<style lang="scss">
.availability-bar-filters {
  &__show-type, &__interval {
    max-width: 200px;
  }
}
</style>
