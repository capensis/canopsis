<template lang="pug">
  widget-settings-item(:title="$t('settings.counters')")
    v-layout(align-center)
      v-switch(
        v-field="value.enabled",
        color="primary",
        hide-details
      )
      c-pbehavior-type-field(
        v-field="value.types",
        :required="!value.enabled",
        :disabled="!value.enabled",
        :is-item-disabled="isItemDisabled",
        with-icon,
        chips,
        multiple
      )
</template>

<script>
import { COUNTERS_LIMIT } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/settings/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  props: {
    value: {
      type: Object,
      default: () => ({ enabled: false, types: [] }),
    },
  },
  methods: {
    isItemDisabled(item) {
      const { types } = this.value;

      return types.length === COUNTERS_LIMIT && !types.includes(item._id);
    },
  },
};
</script>
