<template>
  <widget-settings-item
    :title="title"
    optional
  >
    <c-number-field v-field="value" />
  </widget-settings-item>
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

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
      required: false,
    },
    title: {
      type: String,
      default: '',
    },
  },
};
</script>
