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
import { MODALS } from '@/constants';

import modalMixin from '@/mixins/modal';

export default {
  mixins: [modalMixin],
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    hiddenFields: {
      type: Array,
      default: () => [],
    },
  },
  methods: {
    openFilterModal() {
      this.showModal({
        name: MODALS.createFilter,
        config: {
          title: this.$t('modals.filter.create.title'),
          filter: this.value,
          hiddenFields: this.hiddenFields,
          action: filterObject => this.$emit('input', filterObject),
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
