<template>
  <widget-settings-item :title="title">
    <v-slider
      v-field="value"
      :max="max"
      :min="min"
      ticks="always"
      always-dirty
      thumb-label
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
    value: {
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
};
</script>
