<template lang="pug">
  v-container.pa-3(fluid)
    v-layout(align-center, justify-space-between)
      div.subheading {{ $t('settings.infoPopup.title') }}
      v-layout(justify-end)
        v-btn.primary(
        data-test="infoPopupButton",
        small,
        @click="edit"
        ) {{ $t('common.create') }}/{{ $t('common.edit') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalMixin from '@/mixins/modal';

export default {
  mixins: [
    modalMixin,
  ],
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
      this.showModal({
        name: MODALS.infoPopupSetting,
        config: {
          infoPopups: this.popups,
          columns: this.columns,
          action: popups => this.$emit('input', popups),
        },
      });
    },
  },
};
</script>
