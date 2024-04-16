<template>
  <v-layout class="healthcheck-last-hour-filters gap-6" wrap>
    <healthcheck-event-types-field
      :value="eventTypes"
      class="healthcheck-last-hour-filters__event-types"
      @input="$emit('update:event-types', $event)"
    />
    <healthcheck-connector-names-field
      :value="connectorNames"
      class="healthcheck-last-hour-filters__connector-names"
      @input="$emit('update:connector-names', $event)"
    />
  </v-layout>
</template>

<script>
import HealthcheckEventTypesField from '@/components/other/healthcheck/form/fields/healthcheck-event-types-field.vue';

import HealthcheckConnectorNamesField from '../form/fields/healthcheck-connector-names-field.vue';

export default {
  components: { HealthcheckEventTypesField, HealthcheckConnectorNamesField },
  props: {
    eventTypes: {
      type: Array,
      required: true,
    },
    connectorNames: {
      type: Array,
      required: true,
    },
  },
};
</script>

<style lang="scss">
.healthcheck-last-hour-filters {
  &__connector-names, &__event-types {
    max-width: 350px;
  }
}
</style>
