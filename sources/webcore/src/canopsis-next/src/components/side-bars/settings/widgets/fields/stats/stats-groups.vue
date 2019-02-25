<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.statsGroups.title') }}
      .font-italic.caption.ml-1 ({{ $t('settings.statsGroups.required') }})
    v-container
      v-alert(:value="errors.has('groups')", type="error") {{ $t('settings.statsGroups.required') }}
      v-btn(@click="addGroup") {{ $t('settings.statsGroups.manageGroups') }}
      v-list.secondary(dark)
        v-list-tile(v-for="(group, index) in groups", :key="index")
          v-list-tile-content {{ group.title }}
          v-list-tile-action
            v-layout
              v-btn.primary.mx-1(@click="editGroup(group, index)", fab, small, depressed)
                v-icon edit
              v-btn.error(@click.stop="deleteGroup(index)", fab, small, depressed)
                v-icon delete
</template>

<script>
import { MODALS } from '@/constants';

import modalMixin from '@/mixins/modal';
import formArrayMixin from '@/mixins/form/array';

export default {
  inject: ['$validator'],
  mixins: [modalMixin, formArrayMixin],
  model: {
    prop: 'groups',
    event: 'input',
  },
  props: {
    groups: {
      type: Array,
      default: () => [],
    },
  },
  watch: {
    groups(value) {
      this.$validator.validate('groups', value);
    },
  },
  created() {
    this.$validator.attach({
      name: 'groups',
      rules: 'required',
      getter: () => this.groups,
      context: () => this,
    });
  },
  methods: {
    addGroup() {
      this.showModal({
        name: MODALS.manageHistogramGroups,
        config: {
          title: 'modals.manageHistogramGroups.title.add',
          action: newGroup => this.addItemIntoArray(newGroup),
        },
      });
    },
    editGroup(group, index) {
      this.showModal({
        name: MODALS.manageHistogramGroups,
        config: {
          title: 'modals.manageHistogramGroups.title.edit',
          group,
          action: newGroup => this.updateItemInArray(index, newGroup),
        },
      });
    },
    deleteGroup(index) {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: () => this.removeItemFromArray(index),
        },
      });
    },
  },
};
</script>

