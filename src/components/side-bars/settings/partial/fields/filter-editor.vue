<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.filterEditor') }}
    div.text-xs-center
      v-btn(@click="openFilterModal") {{ $t('modals.filter.create.title') }}
</template>

<script>
import { MODALS } from '@/constants';
import modalMixin from '@/mixins/modal/modal';

export default {
  mixins: [modalMixin],
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
          filter: {},
          action: newFilter => this.$emit('update:filter', newFilter),
        },
      });
    },
  },
};
</script>
