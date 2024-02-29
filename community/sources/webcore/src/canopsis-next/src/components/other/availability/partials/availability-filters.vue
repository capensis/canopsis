<template>
  <v-layout class="availability-filters gap-6" wrap>
    <c-quick-date-interval-field
      :interval="interval"
      :min="minIntervalDate"
      :quick-ranges="quickRanges"
      class="availability-filters__interval"
      short
      @input="$emit('update:interval', $event)"
    />
    <availability-show-type-field
      :value="showType"
      :label="$t('common.show')"
      class="availability-filters__show-type"
      @input="$emit('update:showType', $event)"
    />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { AVAILABILITY_QUICK_RANGES } from '@/constants';

import AvailabilityShowTypeField from '@/components/other/availability/form/fields/availability-show-type-field.vue';

export default {
  components: { AvailabilityShowTypeField },
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
.availability-filters {
  &__show-type, &__interval {
    max-width: 200px;
    flex: 1;
  }
}
</style>
