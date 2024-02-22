<template>
  <widget-settings-item :title="$t('settings.defaultSampling')">
    <c-sampling-field
      v-field="value"
      :name="name"
    />
  </widget-settings-item>
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  props: {
    value: {
      type: String,
      required: false,
    },
    name: {
      type: String,
      required: false,
    },
  },
};
</script>
