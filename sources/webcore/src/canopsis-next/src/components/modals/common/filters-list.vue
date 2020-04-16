<template lang="pug">
  modal-wrapper(data-test="filtersListModal")
    template(slot="title")
      span {{ $t('common.filters') }}
    template(slot="text")
      filters-list-component(
        :filters="filters",
        :entitiesType="config.entitiesType",
        :hasAccessToAddFilter="config.hasAccessToAddFilter",
        :hasAccessToEditFilter="config.hasAccessToEditFilter",
        @create:filter="createFilter",
        @update:filter="updateFilter",
        @delete:filter="deleteFilter",
        @update:filters="updateFilters"
      )
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import FiltersListComponent from '@/components/other/filter/list/filters-list.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Confirmation modal
 */
export default {
  name: MODALS.filtersList,
  components: { FiltersListComponent, ModalWrapper },
  mixins: [modalInnerMixin],
  data() {
    const { filters = [] } = this.modal.config;

    return {
      filters: [...filters],
    };
  },
  computed: {
    actions() {
      return this.config.actions || {};
    },
  },
  methods: {
    createFilter(newFilter) {
      if (this.actions.create) {
        this.filters = [...this.actions.create(newFilter)];
      }
    },

    updateFilter(newFilter, index) {
      if (this.actions.update) {
        this.filters = [...this.actions.update(newFilter, index)];
      }
    },

    deleteFilter(index) {
      if (this.actions.delete) {
        this.filters = [...this.actions.delete(index)];
      }
    },

    updateFilters(filters) {
      if (this.actions.updateList) {
        this.filters = [...this.actions.updateList(filters)];
      }
    },
  },
};
</script>
