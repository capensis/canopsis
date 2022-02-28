<template lang="pug">
  div
    v-alert(
      v-if="!filters.length",
      :value="true",
      type="info"
    ) {{ $t('modals.filter.emptyFilters') }}
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
import { MODALS, ENTITIES_TYPES } from '@/constants';

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
    entitiesType: {
      type: String,
      default: ENTITIES_TYPES.alarm,
      validator: value => [ENTITIES_TYPES.alarm, ENTITIES_TYPES.entity].includes(value),
    },
  },
  computed: {
    existingTitles() {
      return this.filters.map(({ title }) => title);
    },
  },
  methods: {
    showCreateFilterModal() {
      this.$modals.show({
        name: MODALS.createFilter,
        config: {
          title: this.$t('modals.filter.create.title'),
          entitiesType: this.entitiesType,
          existingTitles: this.existingTitles,
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

          title: this.$t('modals.filter.edit.title'),
          entitiesType: this.entitiesType,
          existingTitles: this.existingTitles,
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
