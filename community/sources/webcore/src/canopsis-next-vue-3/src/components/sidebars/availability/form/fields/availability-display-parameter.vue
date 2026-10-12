<template>
  <widget-settings-item :title="$t('settings.availability.defaultParameterToDisplay')">
    <availability-display-parameter-radio-field
      :model-value="modelValue"
      :name="name"
      class="mt-0"
      @update:model-value="$emit('update:modelValue', $event)"
    />
  </widget-settings-item>
</template>

<script>
import AvailabilityDisplayParameterRadioField from '@/components/other/availability/form/fields/availability-display-parameter-radio-field.vue';
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem, AvailabilityDisplayParameterRadioField },

  props: {
    modelValue: {
      type: Number,
      required: false,
    },
    name: {
      type: String,
      default: 'default_display_parameter',
    },
  },

  emits: ['update:modelValue'],
};
</script>
