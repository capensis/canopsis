<template lang="pug">
  widget-settings-item(:title="$t('common.title')", optional)
    v-text-field(
      v-field="value",
      :placeholder="$t('settings.widgetTitle')"
    )
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/settings/partials/widget-settings-item.vue';

/**
 * Component to set the title of a widget
 *
 * @prop {String} [value] - Value of the title
 */
export default {
  components: { WidgetSettingsItem },
  props: {
    value: {
      type: String,
      required: false,
    },
  },
};
</script>
