<template>
  <widget-settings-item :title="$t('settings.density.title')">
    <c-density-btn-toggle
      :value="modelValue"
      :name="name"
      column
      @change="$emit('update:modelValue', $event)"
    />
  </widget-settings-item>
</template>

<script>
import { DENSE_TYPES } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },

  props: {
    modelValue: {
      type: Number,
      default: DENSE_TYPES.large,
    },
    name: {
      type: String,
      default: 'dense',
    },
  },

  emits: ['update:modelValue'],
};
</script>
