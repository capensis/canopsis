<template lang="pug">
  div
    v-alert(
      v-if="!filters.length",
      :value="true",
      type="info"
    ) {{ $t('modals.filter.emptyFilters') }}
    draggable(
      v-else,
      :value="filters",
      :options="draggableOptions",
      element="v-list",
      @change="changeFiltersOrdering"
    )
      filter-field(
        v-for="(filter, index) in filters",
        :filter="filters[index]",
        :key="filter.title",
        :has-access-to-edit="hasAccessToEditFilter",
        @edit="showEditFilterModal(index)",
        @delete="showDeleteFilterModal(index)"
      )
    v-btn.ml-0(
      v-if="hasAccessToAddFilter",
      color="primary",
      outline,
      @click.prevent="showCreateFilterModal"
    ) {{ $t('common.addFilter') }}
</template>

<script>
import Draggable from 'vuedraggable';

import { formArrayMixin } from '@/mixins/form';

import { dragDropChangePositionHandler } from '@/helpers/dragdrop';
import { VUETIFY_ANIMATION_DELAY } from '@/config';
import { MODALS, ENTITIES_TYPES } from '@/constants';

import FilterField from '@/components/other/filter/form/fields/filter-field.vue';

export default {
  components: { Draggable, FilterField },
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
      validator: value => [ENTITIES_TYPES.alarm, ENTITIES_TYPES.entity].includes(value),
    },
  },
  computed: {
    existingTitles() {
      return this.filters.map(({ title }) => title);
    },

    draggableOptions() {
      return {
        animation: VUETIFY_ANIMATION_DELAY,
        disabled: !this.hasAccessToEditFilter,
      };
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

    changeFiltersOrdering(event) {
      this.updateModel(dragDropChangePositionHandler(this.filters, event));
    },
  },
};
</script>
