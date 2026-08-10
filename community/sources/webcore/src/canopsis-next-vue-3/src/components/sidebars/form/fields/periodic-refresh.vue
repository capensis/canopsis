<template>
  <widget-settings-item :title="$t('view.periodicRefresh')">
    <periodic-refresh-field
      :periodic-refresh="modelValue"
      name="parameters.periodic_refresh"
      @update:periodic-refresh="$emit('update:modelValue', $event)"
    />
  </widget-settings-item>
</template>

<script>
import PeriodicRefreshField from '@/components/forms/fields/periodic-refresh-field.vue';
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: {
    PeriodicRefreshField,
    WidgetSettingsItem,
  },

  props: {
    modelValue: {
      type: Object,
      default: () => ({}),
    },
  },

  emits: ['update:modelValue'],
};
</script>
