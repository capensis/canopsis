<template>
  <div class="healthcheck-last-hour-filters col-gap-6">
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
  </div>
</template>

<script>
import HealthcheckConnectorNamesField from '../form/fields/healthcheck-connector-names-field.vue'
import HealthcheckEventTypesField from '../form/fields/healthcheck-event-types-field.vue'

export default {
  components: {
    HealthcheckConnectorNamesField,
    HealthcheckEventTypesField,
  },
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

  emits: ['update:event-types', 'update:connector-names'],
}
</script>

<style lang="scss">
.healthcheck-last-hour-filters {
  display: flex;
  flex-wrap: wrap;

  & > * {
    flex-grow: 0;
  }

  &__connector-names,
  &__event-types {
    width: 340px;
  }
}
</style>
