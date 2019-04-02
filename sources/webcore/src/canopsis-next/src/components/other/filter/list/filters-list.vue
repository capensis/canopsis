<template lang="pug">
  div
    v-layout(wrap, justify-center)
      v-list
        v-list-tile(v-for="(filter, index) in filters", :key="filter.value")
          v-list-tile-content {{ filter.title }}
          v-list-tile-action(v-if="hasAccessToEditFilter")
            div
              v-btn.ma-1(icon, @click="showEditFilterModal(index)")
                v-icon edit
              v-btn.ma-1(icon, @click="showDeleteFilterModal(index)")
                v-icon delete
    v-btn.ml-0(
    v-if="hasAccessToAddFilter",
    color="primary",
    @click.prevent="showCreateFilterModal"
    ) {{ $t('common.add') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalMixin from '@/mixins/modal';

export default {
  mixins: [modalMixin],
  props: {
    filters: {
      type: Array,
      default: () => [],
    },
    hasAccessToAddFilter: {
      type: Boolean,
      default: true,
    },
    hasAccessToEditFilter: {
      type: Boolean,
      default: true,
    },
  },
  methods: {
    showCreateFilterModal() {
      this.showModal({
        name: MODALS.createFilter,
        config: {
          title: this.$t('modals.filter.create.title'),
          action: (newFilter) => {
            this.$emit('create:filter', newFilter);
            this.$emit('update:filters', [...this.filters, newFilter]);
          },
        },
      });
    },

    showEditFilterModal(index) {
      const filter = this.filters[index];

      this.showModal({
        name: MODALS.createFilter,
        config: {
          title: this.$t('modals.filter.edit.title'),
          filter,
          action: (newFilter) => {
            this.$emit('update:filter', newFilter, index);
            this.$emit('update:filters', this.filters.map((v, i) => (index === i ? newFilter : v)));
          },
        },
      });
    },

    showDeleteFilterModal(index) {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: () => {
            this.$emit('delete:filter', index);
            this.$emit('update:filters', this.filters.filter((v, i) => index !== i));
          },
        },
      });
    },
  },
};
</script>
