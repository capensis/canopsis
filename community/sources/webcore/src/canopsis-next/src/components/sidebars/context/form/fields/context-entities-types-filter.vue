<template>
  <widget-settings-item :title="$t('settings.contextTypeOfEntities.title')">
    <c-checkbox
      v-for="entityType in entitiesTypes"
      :key="entityType.value"
      :model-value="isSelected(entityType.value)"
      :label="entityType.label"
      hide-details
      @update:model-value="updateType(entityType.value, $event)"
    />
  </widget-settings-item>
</template>

<script>
import { ENTITY_TYPES } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

/**
 * Component to select entities type to filter on entities-list
 *
 * @prop {Array} [modelValue] - Array of selected entities types values to filter on
 *
 * @event modelValue#update:modelValue
 */
export default {
  components: { WidgetSettingsItem },

  props: {
    modelValue: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['update:modelValue'],

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

  methods: {
    isSelected(type) {
      return this.modelValue.includes(type);
    },

    updateType(type, selected) {
      const nextValue = selected
        ? [...this.modelValue, type]
        : this.modelValue.filter(selectedType => selectedType !== type);

      this.$emit('update:modelValue', nextValue);
    },
  },
};
</script>
