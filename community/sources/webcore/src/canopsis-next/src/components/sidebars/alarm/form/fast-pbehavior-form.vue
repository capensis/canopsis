<template>
  <widget-settings-item :title="$t('settings.fastPbehavior')">
    <v-text-field
      v-field="form.fastPbehaviorNamePrefix"
      :label="$t('common.namePrefix')"
    />
    <c-pbehavior-type-field
      v-field="form.fastPbehaviorType"
      :label="$t('common.type')"
      :types="types"
      independent
      clearable
    />
    <c-pbehavior-reason-field
      v-field="form.fastPbehaviorReason"
      :label="$t('common.reason')"
      clearable
    />
  </widget-settings-item>
</template>
<script>
import { PBEHAVIOR_TYPE_TYPES } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
  setup() {
    return {
      types: [PBEHAVIOR_TYPE_TYPES.pause],
    };
  },
};
</script>
