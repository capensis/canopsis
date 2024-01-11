<template>
  <widget-settings-flat-item :title="$t('settings.infoPopup.title')">
    <template #actions>
      <v-btn
        class="primary"
        small
        @click="edit"
      >
        {{ $t('common.create') }}/{{ $t('common.edit') }}
      </v-btn>
    </template>
  </widget-settings-flat-item>
</template>

<script>
import { MODALS } from '@/constants';

import { formMixin } from '@/mixins/form';

import WidgetSettingsFlatItem from '@/components/sidebars/partials/widget-settings-flat-item.vue';

export default {
  components: { WidgetSettingsFlatItem },
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
