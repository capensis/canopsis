<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.filterEditor') }}
    div.text-xs-center
      v-btn(@click="openFilterModal") {{ $t('modals.filter.create.title') }}
      v-btn(@click="deleteFilter", icon)
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
