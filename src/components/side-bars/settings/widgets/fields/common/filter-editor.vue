<template lang="pug">
v-container.pa-3(fluid)
  v-layout(align-center, justify-space-between)
    div.subheading {{ $t('settings.filterEditor') }}
      .font-italic.caption.ml-1 ({{ $t('common.optionnal') }})
    div
      v-btn.green.darken-4.white--text.caption.my-0.mx-1(
      small,
      @click="openFilterModal"
      ) {{ $t('common.create') }}/{{ $t('common.edit') }}
      v-btn.red.darken-4.white--text.caption.my-0.mx-1(small, @click="deleteFilter")
        v-icon delete
</template>

<script>
import { MODALS } from '@/constants';
import modalMixin from '@/mixins/modal/modal';

export default {
  mixins: [modalMixin],
  props: {
    value: {
      type: Object,
    },
  },
  data() {
    return {
      item: {},
    };
  },
  methods: {
    openFilterModal() {
      this.showModal({
        name: MODALS.createFilter,
        config: {
          title: 'modals.filter.create.title',
          filter: this.value || {},
          action: newFilter => this.$emit('input', newFilter),
        },
      });
    },
    deleteFilter() {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: () => this.$emit('input', {}),
        },
      });
    },
  },
};
</script>
