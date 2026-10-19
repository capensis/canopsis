<template>
  <c-layout class="kpi-sli-toolbar ml-4 my-4" gap="3" wrap>
    <c-quick-date-interval-field
      :interval="query.interval"
      :min="minDate"
      :quick-ranges="quickRanges"
      @update:interval="updateQueryField('interval', $event)"
    />
    <c-sampling-field
      :value="query.sampling"
      class="kpi-sli-toolbar__sampling"
      @update:value="updateQueryField('sampling', $event)"
    />
    <kpi-sli-show-mode-field
      :value="query.type"
      class="kpi-sli-toolbar__show-mode"
      @update:value="updateQueryField('type', $event)"
    />
  </c-layout>
</template>

<script>
import { METRICS_QUICK_RANGES } from '@/constants';

import KpiSliShowModeField from '../form/fields/kpi-sli-show-mode-field.vue';

export default {
  components: {
    KpiSliShowModeField,
  },

  props: {
    query: {
      type: Object,
      required: true,
    },
    minDate: {
      type: Number,
      required: false,
    },
  },

  emits: ['update:query'],

  computed: {
    quickRanges() {
      return Object.values(METRICS_QUICK_RANGES);
    },
  },

  methods: {
    updateQueryField(field, value) {
      this.$emit('update:query', {
        ...this.query,
        [field]: value,
      });
    },
  },
};
</script>

<style lang="scss">
.kpi-sli-toolbar {
  &__sampling {
    max-width: 200px;
  }

  &__show-mode {
    max-width: 150px;
  }
}
</style>
