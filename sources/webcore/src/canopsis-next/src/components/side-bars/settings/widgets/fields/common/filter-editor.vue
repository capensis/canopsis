<template lang="pug">
  v-container.pa-3(fluid)
    v-layout(align-center, justify-space-between)
      div.subheading {{ $t('settings.filterEditor') }}
        .font-italic.caption.ml-1(v-show="!required") ({{ $t('common.optional') }})
      div
        v-btn.primary(
          small,
          @click="openFilterModal"
        )
          span(v-show="isValueEmpty") {{ $t('common.create') }}
          span(v-show="!isValueEmpty") {{ $t('common.edit') }}
        v-btn.error(v-show="!isValueEmpty", small, @click="deleteFilter")
          v-icon delete
</template>

<script>
import { isEmpty } from 'lodash';

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
    required: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isValueEmpty() {
      return isEmpty(this.value);
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
