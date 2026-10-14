<template>
  <c-layout class="availability-history-filters" gap="6" wrap>
    <c-sampling-field
      :value="sampling"
      :samplings="availableSamplings"
      class="availability-history-filters__sampling"
      @update:value="$emit('update:sampling', $event)"
    />
    <availability-show-type-field
      :model-value="showType"
      :label="$t('common.show')"
      class="availability-history-filters__show-type"
      @update:model-value="$emit('update:showType', $event)"
    />
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { SAMPLINGS } from '@/constants';

import AvailabilityShowTypeField from '@/components/other/availability/form/fields/availability-show-type-field.vue';

export default {
  components: { AvailabilityShowTypeField },

  props: {
    sampling: {
      type: String,
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

  emits: ['update:sampling', 'update:showType'],

  setup() {
    const availableSamplings = computed(() => [SAMPLINGS.hour, SAMPLINGS.day]);

    return {
      availableSamplings,
    };
  },
};
</script>

<style lang="scss">
.availability-history-filters {
  &__sampling,
  &__show-type {
    max-width: 200px;
    flex: 1;
  }
}
</style>
