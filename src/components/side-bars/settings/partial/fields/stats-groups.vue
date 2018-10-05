<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.statsGroups.title') }}
    v-container.pt-0
      v-btn.ma-0(@click="showGroupsModal") {{ $t('settings.statsGroups.manageGroups') }}
</template>

<script>
import modalMixin from '@/mixins/modal/modal';
import { MODALS } from '@/constants';

export default {
  mixins: [modalMixin],
  props: {
    groups: {
      type: Array,
    },
  },
  methods: {
    showGroupsModal() {
      this.showModal({
        name: MODALS.manageHistogramGroups,
        config: {
          groups: this.groups,
          action: newGroups => this.$emit('update:groups', newGroups),
        },
      });
    },
  },
};
</script>

