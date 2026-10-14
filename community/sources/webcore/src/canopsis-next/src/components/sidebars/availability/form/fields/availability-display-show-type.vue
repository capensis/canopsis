<template>
  <widget-settings-item :title="$t('settings.availability.defaultAvailabilityDisplay')">
    <availability-show-type-radio-field
      :model-value="modelValue"
      :name="name"
      class="mt-0"
      @update:model-value="$emit('update:modelValue', $event)"
    />
  </widget-settings-item>
</template>

<script>
import AvailabilityShowTypeRadioField from '@/components/other/availability/form/fields/availability-show-type-radio-field.vue';
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem, AvailabilityShowTypeRadioField },

  props: {
    modelValue: {
      type: Number,
      required: false,
    },
    name: {
      type: String,
      default: 'default_show_type',
    },
  },

  emits: ['update:modelValue'],
};
</script>
