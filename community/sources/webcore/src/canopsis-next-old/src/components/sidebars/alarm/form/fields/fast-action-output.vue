<template>
  <widget-settings-item
    :title="label"
    optional
  >
    <c-enabled-field v-field="value.enabled" />
    <v-text-field
      v-field="value.value"
      :label="$t('common.output')"
      :disabled="!value.enabled"
    />
  </widget-settings-item>
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

/**
 * Component to set fast-ack output
 *
 * @prop {Object} [value] - Value
 */
export default {
  components: { WidgetSettingsItem },
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    label: {
      type: String,
      default: '',
    },
  },
};
</script>
