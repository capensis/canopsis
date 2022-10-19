<template lang="pug">
  widget-settings-item(:title="$t('settings.defaultNumberOfElementsPerPage')", :sub-title="subTitle")
    v-select.select(
      v-field="value",
      :label="$t('settings.elementsPerPage')",
      :items="$config.PAGINATION_PER_PAGE_VALUES",
      hide-details,
      single-line,
      dense
    )
</template>

<script>
import { PAGINATION_LIMIT } from '@/config';

import WidgetSettingsItem from '@/components/sidebars/settings/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  props: {
    value: {
      type: Number,
      default: PAGINATION_LIMIT,
    },
    subTitle: {
      type: String,
      required: false,
    },
  },
};
</script>
