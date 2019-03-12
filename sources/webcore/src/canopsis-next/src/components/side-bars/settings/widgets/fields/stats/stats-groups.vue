<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.statsGroups.title') }}
      .font-italic.caption.ml-1 ({{ $t('settings.statsGroups.required') }})
    v-container
      v-alert(:value="errors.has('groups')", type="error") {{ $t('settings.statsGroups.required') }}
      v-btn(@click="showAddGroupModal") {{ $t('settings.statsGroups.manageGroups') }}
      v-list.secondary(dark)
        v-list-tile(v-for="(group, index) in groups", :key="index")
          v-list-tile-content {{ group.title }}
          v-list-tile-action
            v-layout
              v-btn.primary.mx-1(@click="showEditGroupModal(group, index)", fab, small, depressed)
                v-icon edit
              v-btn.error(@click.stop="showDeleteGroupModal(index)", fab, small, depressed)
                v-icon delete
</template>

<script>
import { MODALS } from '@/constants';

import modalMixin from '@/mixins/modal';
import formMixinArray from '@/mixins/form/array';

export default {
  inject: ['$validator'],
  filters: {
    groupToFilter(group = {}) {
      return {
        title: group.title || '',
        filter: group.filter && group.filter.filter ? group.filter.filter : {},
      };
    },
    filterToGroup(filter = {}) {
      return {
        title: filter.title || '',
        filter: {
          filter: filter.filter || {},
        },
      };
    },
  },
  mixins: [modalMixin, formMixinArray],
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
    showAddGroupModal() {
      const defaultFilter = { title: '', filter: {} };

      this.showModal({
        name: MODALS.createFilter,
        config: {
          title: this.$t('modals.manageHistogramGroups.title.add'),
          filter: defaultFilter,
          action: newFilter => this.addItemIntoArray(this.$options.filters.filterToGroup(newFilter)),
        },
      });
    },

    showEditGroupModal(group, index) {
      this.showModal({
        name: MODALS.createFilter,
        config: {
          title: this.$t('modals.manageHistogramGroups.title.edit'),
          filter: this.$options.filters.groupToFilter(group),
          action: newFilter => this.updateItemInArray(index, this.$options.filters.filterToGroup(newFilter)),
        },
      });
    },

    showDeleteGroupModal(index) {
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

