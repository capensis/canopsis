<template lang="pug">
  v-container.pa-3(fluid)
    v-layout(align-center, justify-space-between)
      div.subheading {{ $t('settings.infoPopup.title') }}
      v-layout(justify-end)
        v-btn.primary(
          small,
          @click="edit"
        ) {{ $t('common.create') }}/{{ $t('common.edit') }}
</template>

<script>
import { MODALS } from '@/constants';

import { formMixin } from '@/mixins/form';

export default {
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
