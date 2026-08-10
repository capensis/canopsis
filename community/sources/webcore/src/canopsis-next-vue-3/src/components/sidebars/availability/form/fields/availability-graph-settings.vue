<template>
  <widget-settings-item :title="$t('settings.availability.graphSettings')">
    <c-enabled-field
      :model-value="modelValue.enabled"
      @update:model-value="updateField('enabled', $event)"
    />
    <v-expand-transition>
      <div v-if="modelValue.enabled">
        <c-quick-date-interval-type-field
          :value="modelValue.default_time_range"
          :name="defaultTimeRangeFieldName"
          :ranges="intervalRanges"
          :label="$t('settings.defaultTimeRange')"
          @update:value="updateField('default_time_range', $event)"
        />
        <availability-show-type-radio-field
          :model-value="modelValue.default_show_type"
          :label="$t('settings.availability.defaultAvailabilityDisplay')"
          :name="defaultShowTypeFieldName"
          @update:model-value="updateField('default_show_type', $event)"
        />
      </div>
    </v-expand-transition>
  </widget-settings-item>
</template>

<script>
import AvailabilityShowTypeRadioField from '@/components/other/availability/form/fields/availability-show-type-radio-field.vue';
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';
import { AVAILABILITY_QUICK_RANGES } from '@/constants';

export default {
  components: { WidgetSettingsItem, AvailabilityShowTypeRadioField },

  props: {
    modelValue: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      default: 'availability',
    },
  },

  emits: ['update:modelValue'],

  computed: {
    defaultTimeRangeFieldName() {
      return `${this.name}.default_time_range`;
    },

    defaultShowTypeFieldName() {
      return `${this.name}.show_type`;
    },

    intervalRanges() {
      return Object.values(AVAILABILITY_QUICK_RANGES);
    },
  },

  methods: {
    updateField(field, value) {
      this.$emit('update:modelValue', {
        ...this.modelValue,
        [field]: value,
      });
    },
  },
};
</script>
