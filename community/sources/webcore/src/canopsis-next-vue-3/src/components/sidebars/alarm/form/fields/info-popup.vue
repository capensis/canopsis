<template>
  <widget-settings-item :title="$t('settings.infoPopup.title')">
    <v-btn
      color="primary"
      size="small"
      @click="edit"
    >
      {{ $t('common.create') }}/{{ $t('common.edit') }}
    </v-btn>
  </widget-settings-item>
</template>

<script>
import { MODALS } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },

  props: {
    modelValue: {
      type: [Array, Object],
      default: () => [],
    },
    columns: {
      type: Array,
      default: () => [],
    },
    variables: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['update:modelValue'],

  methods: {
    edit() {
      this.$modals.show({
        name: MODALS.infoPopupSetting,
        config: {
          infoPopups: this.modelValue,
          columns: this.columns,
          variables: this.variables,
          action: popups => this.$emit('update:modelValue', popups),
        },
      });
    },
  },
};
</script>
