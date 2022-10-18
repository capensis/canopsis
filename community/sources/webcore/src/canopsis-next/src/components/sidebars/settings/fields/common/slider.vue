<template lang="pug">
  widget-settings-item(:title="title")
    v-slider(
      :value="value",
      :max="max",
      :min="min",
      ticks="always",
      @input="$emit('input', $event)",
      always-dirty,
      thumb-label
    )
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/settings/partials/widget-settings-item.vue';

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
