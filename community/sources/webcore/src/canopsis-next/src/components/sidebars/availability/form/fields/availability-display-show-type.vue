<template>
  <widget-settings-item :title="$t('settings.availability.defaultAvailabilityDisplay')">
    <availability-show-type-radio-field v-field="value" :name="name" class="mt-0" />
  </widget-settings-item>
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';
import AvailabilityShowTypeRadioField from '@/components/other/availability/form/fields/availability-show-type-radio-field.vue';

export default {
  components: { WidgetSettingsItem, AvailabilityShowTypeRadioField },
  props: {
    value: {
      type: Number,
      required: false,
    },
    name: {
      type: String,
      default: 'default_show_type',
    },
  },
};
</script>
