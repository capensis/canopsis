<template>
  <widget-settings-item :title="$t('settings.contextTypeOfEntities.title')">
    <v-checkbox
      v-for="entityType in entitiesTypes"
      :key="entityType.value"
      :input-value="value"
      :label="entityType.label"
      :value="entityType.value"
      color="primary"
      hide-details
      @change="$listeners.input"
    />
  </widget-settings-item>
</template>

<script>
import { ENTITY_TYPES } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

/**
 * Component to select entities type to filter on entities-list
 *
 * @prop {Array} [value] - Array of selected entities types values to filter on
 *
 * @event value#input
 */
export default {
  components: { WidgetSettingsItem },
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
