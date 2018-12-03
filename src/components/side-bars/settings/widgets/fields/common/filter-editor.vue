<template lang="pug">
v-container.pa-3(fluid)
  v-layout(align-center, justify-space-between)
    div.subheading {{ $t('settings.filterEditor') }}
      .font-italic.caption.ml-1 ({{ $t('common.optionnal') }})
    div
      v-btn.primary(
      small,
      @click="openFilterModal"
      ) {{ $t('common.create') }}/{{ $t('common.edit') }}
      v-btn.error(small, @click="deleteFilter")
        v-icon delete
</template>

<script>
import isEmpty from 'lodash/isEmpty';

import modalMixin from '@/mixins/modal';

export default {
  mixins: [modalMixin],
  props: {
    value: {
      type: Object,
    },
  },
  computed: {
    openFilterButtonText() {
      if (isEmpty(this.value)) {
        return this.$t('modals.filter.create.title');
      }

      return this.$t('modals.filter.edit.title');
    },
  },
  methods: {
    openFilterModal() {
      this.showModal({
        name: this.$constants.MODALS.createFilter,
        config: {
          title: 'modals.filter.create.title',
          filter: this.value || {},
          action: newFilter => this.$emit('input', newFilter),
        },
      });
    },
    deleteFilter() {
      this.showModal({
        name: this.$constants.MODALS.confirmation,
        config: {
          action: () => this.$emit('input', {}),
        },
      });
    },
  },
};
</script>
