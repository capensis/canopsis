<template lang="pug">
  widget-settings-item(:title="title", optional)
    v-text-field(:value="value", type="number", @input="$emit('input', parseInt($event, 10))")
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/settings/partials/widget-settings-item.vue';
/**
 * Component to select a Number type setting
 *
 * @prop {Number} [value=undefined] - Value
 * @prop {String} [title=''] - Title
 */
export default {
  components: { WidgetSettingsItem },
  props: {
    value: {
      type: Number,
      default: undefined,
    },
    title: {
      type: String,
      default: '',
    },
  },
};
</script>
