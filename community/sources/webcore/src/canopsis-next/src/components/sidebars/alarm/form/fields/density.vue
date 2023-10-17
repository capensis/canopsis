<template>
  <widget-settings-item :title="$t('settings.density.title')">
    <c-density-btn-toggle
      v-field="value"
      :name="name"
      column
    />
  </widget-settings-item>
</template>

<script>
import { ALARM_DENSE_TYPES } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  props: {
    value: {
      type: Number,
      default: ALARM_DENSE_TYPES.large,
    },
    name: {
      type: String,
      default: 'dense',
    },
  },
};
</script>

<style lang="scss" scoped>
.density {
  &-icon {
    padding: 2px;
    border-radius: 5px;
    border: 1px solid #707070;
  }

  &-radio-group ::v-deep .v-input__control {
    width: 100%;
  }
}
</style>
