<template>
  <div class="healthcheck-history-filters col-gap-6">
    <c-quick-date-interval-field
      :interval="interval"
      :quick-ranges="quickRanges"
      :min="deletedBefore"
      class="healthcheck-history-filters__interval"
      short
      @input="$emit('update:interval', $event)"
    />
    <healthcheck-event-types-field
      :value="eventTypes"
      class="healthcheck-history-filters__event-types"
      @input="$emit('update:event-types', $event)"
    />
    <healthcheck-connector-names-field
      :value="connectorNames"
      class="healthcheck-history-filters__connector-names"
      @input="$emit('update:connector-names', $event)"
    />
  </div>
</template>

<script>
import { computed } from 'vue';

import { HEALTHCHECK_QUICK_RANGES } from '@/constants';

import HealthcheckEventTypesField from '../form/fields/healthcheck-event-types-field.vue';
import HealthcheckConnectorNamesField from '../form/fields/healthcheck-connector-names-field.vue';

export default {
  components: {
    HealthcheckConnectorNamesField,
    HealthcheckEventTypesField,
  },
  props: {
    interval: {
      type: Object,
      required: true,
    },
    deletedBefore: {
      type: Number,
      required: false,
    },
    eventTypes: {
      type: Array,
      required: true,
    },
    connectorNames: {
      type: Array,
      required: true,
    },
  },
  setup() {
    const quickRanges = computed(() => Object.values(HEALTHCHECK_QUICK_RANGES));

    return {
      quickRanges,
    };
  },
};
</script>

<style lang="scss">
.healthcheck-history-filters {
  display: flex;
  flex-wrap: wrap;

  & > * {
    flex-grow: 0;
  }

  &__interval {
    width: 280px;
  }

  &__connector-names, &__event-types {
    width: 340px;
  }
}
</style>
