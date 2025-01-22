<template>
  <widget-settings-item :title="$t('settings.infoPopup.title')">
    <v-btn
      class="primary"
      small
      @click="edit"
    >
      {{ $t('common.create') }}/{{ $t('common.edit') }}
    </v-btn>
  </widget-settings-item>
</template>

<script>
import { MODALS } from '@/constants';

import { formMixin } from '@/mixins/form';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  mixins: [formMixin],
  model: {
    prop: 'popups',
    event: 'input',
  },
  props: {
    popups: {
      type: [Array, Object],
      default: () => [],
    },
    columns: {
      type: Array,
      default: () => [],
    },
  },
  methods: {
    edit() {
      this.$modals.show({
        name: MODALS.infoPopupSetting,
        config: {
          infoPopups: this.popups,
          columns: this.columns,
          action: popups => this.updateModel(popups),
        },
      });
    },
  },
};
</script>
