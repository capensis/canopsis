<template>
  <widget-settings-item :title="title">
    <v-slider
      :model-value="modelValue"
      :max="max"
      :min="min"
      color="primary"
      show-ticks="always"
      thumb-label
      @update:model-value="$emit('update:modelValue', $event)"
    />
  </widget-settings-item>
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },

  props: {
    title: {
      type: String,
      required: true,
    },
    modelValue: {
      type: Number,
      required: true,
    },
    min: {
      type: Number,
      default: 0,
    },
    max: {
      type: Number,
      default: 5,
    },
  },

  emits: ['update:modelValue'],
};
</script>
