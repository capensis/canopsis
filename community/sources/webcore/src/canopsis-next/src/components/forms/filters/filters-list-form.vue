<template lang="pug">
  div
    v-alert(
      v-if="!filters.length",
      :value="true",
      type="info"
    ) {{ $t('modals.createFilter.emptyFilters') }}
    c-draggable-list-field(
      v-else,
      v-field="filters",
      :disabled="!editable",
      component="v-list"
    )
      filter-tile(
        v-for="(filter, index) in filters",
        :filter="filter",
        :key="filter.title",
        :editable="editable",
        @edit="showEditFilterModal(index)",
        @delete="showDeleteFilterModal(index)"
      )
    v-btn.ml-0(
      v-if="addable",
      color="primary",
      outline,
      @click.prevent="showCreateFilterModal"
    ) {{ $t('common.addFilter') }}
</template>

<script>
import { MODALS } from '@/constants';

import { formArrayMixin } from '@/mixins/form';

import FilterTile from './partials/filter-tile.vue';

export default {
  components: { FilterTile },
  mixins: [formArrayMixin],
  model: {
    prop: 'filters',
    event: 'input',
  },
  props: {
    filters: {
      type: Array,
      default: () => [],
    },
    addable: {
      type: Boolean,
      default: true,
    },
    editable: {
      type: Boolean,
      default: true,
    },
    alarm: {
      type: Boolean,
      default: false,
    },
    entity: {
      type: Boolean,
      default: false,
    },
    pbehavior: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    showCreateFilterModal() {
      this.$modals.show({
        name: MODALS.createFilter,
        config: {
          name: true,
          alarm: this.alarm,
          entity: this.entity,
          pbehavior: this.pbehavior,
          title: this.$t('modals.createFilter.create.title'),
          action: newFilter => this.addItemIntoArray(newFilter),
        },
      });
    },

    showEditFilterModal(index) {
      const filter = this.filters[index];

      this.$modals.show({
        name: MODALS.createFilter,
        config: {
          filter,

          name: true,
          alarm: this.alarm,
          entity: this.entity,
          pbehavior: this.pbehavior,
          title: this.$t('modals.createFilter.edit.title'),
          action: newFilter => this.updateItemInArray(index, { ...filter, ...newFilter }),
        },
      });
    },

    showDeleteFilterModal(index) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => this.removeItemFromArray(index),
        },
      });
    },
  },
};
</script>
