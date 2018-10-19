<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.statsGroups.title') }}
      .font-italic.caption.ml-1 ({{ $t('settings.statsGroups.required') }})
    v-container
      v-btn(@click="addGroup") {{ $t('settings.statsGroups.manageGroups') }}
      v-list(dark)
        v-list-tile(v-for="(group, index) in value", :key="index")
          v-list-tile-content {{ group.title }}
          v-list-tile-action
            v-layout
              v-btn.green.darken-4.white--text.mx-1(@click="editGroup(group, index)", fab, small, depressed)
                v-icon edit
              v-btn.red.darken-4.white--text.mx-1(@click.stop="deleteGroup(index)", fab, small, depressed)
                v-icon delete
</template>

<script>
import pullAt from 'lodash/pullAt';
import modalMixin from '@/mixins/modal/modal';
import { MODALS } from '@/constants';

export default {
  mixins: [modalMixin],
  props: {
    value: {
      type: Array,
    },
  },
  methods: {
    addGroup() {
      this.showModal({
        name: MODALS.manageHistogramGroups,
        config: {
          title: 'modals.manageHistogramGroups.title.add',
          action: (newGroup) => {
            const groups = [...this.value];
            groups.push(newGroup);
            this.$emit('input', groups);
          },
        },
      });
    },
    editGroup(group, index) {
      this.showModal({
        name: MODALS.manageHistogramGroups,
        config: {
          title: 'modals.manageHistogramGroups.title.edit',
          group,
          action: (newGroup) => {
            const groups = [...this.value];
            groups[index] = newGroup;
            this.$emit('input', groups);
          },
        },
      });
    },
    deleteGroup(index) {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: () => {
            const groups = [...this.value];
            pullAt(groups, index);
            this.$emit('input', groups);
          },
        },
      });
    },
  },
};
</script>

