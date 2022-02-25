<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.contextTypeOfEntities.title') }}
    v-container
      v-checkbox(
        v-for="entityType in entitiesTypes",
        :input-value="value",
        :label="entityType.label",
        :value="entityType.value",
        :key="entityType.value",
        color="primary",
        hide-details,
        @change="$listeners.input"
      )
</template>

<script>
import { ENTITY_TYPES } from '@/constants';

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
  computed: {
    entitiesTypes() {
      return [
        ENTITY_TYPES.component,
        ENTITY_TYPES.connector,
        ENTITY_TYPES.resource,
        ENTITY_TYPES.service,
      ].map(type => ({
        label: this.$t(`settings.contextTypeOfEntities.fields.${type}`),
        value: type,
      }));
    },
  },
};
</script>
