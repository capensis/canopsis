<template lang="pug">
  div
    v-list
      v-list-tile.pa-0(v-for="(filter, index) in filters", :key="filter.title")
        v-layout
          v-flex(xs12)
            v-list-tile-content {{ filter.title }}
          v-list-tile-action(v-if="hasAccessToEditFilter")
            v-layout
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
import { MODALS, ENTITIES_TYPES } from '@/constants';

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
    entitiesType: {
      type: String,
      default: ENTITIES_TYPES.alarm,
      validator: value => [ENTITIES_TYPES.alarm, ENTITIES_TYPES.entity, ENTITIES_TYPES.pbehavior].includes(value),
    },
  },
  computed: {
    existingTitles() {
      return this.filters.map(({ title }) => title);
    },
  },
  methods: {
    showCreateFilterModal() {
      this.showModal({
        name: MODALS.createFilter,
        config: {
          title: this.$t('modals.filter.create.title'),
          entitiesType: this.entitiesType,
          existingTitles: this.existingTitles,
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
          entitiesType: this.entitiesType,
          existingTitles: this.existingTitles,
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
