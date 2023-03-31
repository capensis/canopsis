<template lang="pug">
  widget-settings-item(:title="title")
    c-storage-field(
      v-field="storage",
      v-on="$listeners",
      :required="required",
      :disabled="disabled"
    )
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/settings/partials/widget-settings-item.vue';

export default {
  inject: ['$validator'],
  components: { WidgetSettingsItem },
  model: {
    prop: 'storage',
    event: 'input',
  },
  props: {
    storage: {
      type: String,
      default: '',
    },
    title: {
      type: String,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
