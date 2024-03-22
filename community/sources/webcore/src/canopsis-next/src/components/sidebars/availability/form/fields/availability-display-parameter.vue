<template>
  <widget-settings-item :title="$t('settings.availability.defaultParameterToDisplay')">
    <availability-display-parameter-radio-field v-field="value" :name="name" class="mt-0" />
  </widget-settings-item>
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';
import AvailabilityDisplayParameterRadioField from '@/components/other/availability/form/fields/availability-display-parameter-radio-field.vue';

export default {
  components: { WidgetSettingsItem, AvailabilityDisplayParameterRadioField },
  props: {
    value: {
      type: Number,
      required: false,
    },
    name: {
      type: String,
      default: 'default_display_parameter',
    },
  },
};
</script>
