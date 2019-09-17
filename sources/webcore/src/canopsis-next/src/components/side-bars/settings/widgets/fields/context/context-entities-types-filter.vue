<template lang="pug">
  v-list-group(data-test="contextTypeOfEntities")
    v-list-tile(slot="activator") {{ $t('settings.contextTypeOfEntities.title') }}
    v-container
      v-checkbox(
      v-for="(entitiesType, index) in entitiesTypes",
      v-bind="entitiesType",
      data-test="entitiesTypeCheckbox",
      :key="`entities-type-${index}`",
      :input-value="value",
      @change="$emit('input', $event)",
      color="primary",
      hide-details
      )
</template>

<script>

/**
 * Component to select entities type to filter on entities-list
 *
 * @prop {Array} [value] - Array of selected entities types values to filter on
 *
 * @event value#input
 */
export default {
  props: {
    value: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      entitiesTypes: [
        { label: this.$t('settings.contextTypeOfEntities.fields.component'), value: 'component' },
        { label: this.$t('settings.contextTypeOfEntities.fields.connector'), value: 'connector' },
        { label: this.$t('settings.contextTypeOfEntities.fields.resource'), value: 'resource' },
        { label: this.$t('settings.contextTypeOfEntities.fields.watcher'), value: 'watcher' },
      ],
    };
  },
};
</script>
