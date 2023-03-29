<template lang="pug">
  widget-settings-item(:title="$t('settings.defaultSampling')")
    c-sampling-field(v-field="value", :name="name")
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/settings/partials/widget-settings-item.vue';

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
