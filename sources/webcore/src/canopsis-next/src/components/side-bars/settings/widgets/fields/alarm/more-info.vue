<template lang="pug">
  v-container.pa-3(fluid)
    v-layout(align-center, justify-space-between)
      div.subheading {{$t('settings.moreInfosModal')}}
      v-layout(justify-end)
        v-btn.primary(
        small,
        @click="openTextEditorModal"
        ) {{ $t('common.create') }}/{{ $t('common.edit') }}
        v-btn.error(v-if="value", small, @click="deleteMoreInfoTemplate")
          v-icon delete
</template>

<script>
import { MODALS } from '@/constants';

import modalMixin from '@/mixins/modal';

export default {
  mixins: [modalMixin],
  props: {
    value: {
      type: String,
      default: '',
    },
  },
  methods: {
    openTextEditorModal() {
      this.showModal({
        name: MODALS.textEditor,
        config: {
          text: this.value,
          action: value => this.$emit('input', value),
        },
      });
    },

    deleteMoreInfoTemplate() {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: () => this.$emit('input', ''),
        },
      });
    },
  },
};
</script>

