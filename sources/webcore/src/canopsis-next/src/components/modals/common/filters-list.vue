<template lang="pug">
  modal-wrapper(data-test="filtersListModal")
    template(slot="title")
      span {{ $t('common.filters') }}
    template(slot="text")
      filters-list-component(
        :filters.sync="filters",
        :hasAccessToAddFilter="config.hasAccessToAddFilter",
        :hasAccessToEditFilter="config.hasAccessToEditFilter",
        @create:filter="createFilter",
        @update:filter="updateFilter",
        @delete:filter="deleteFilter"
      )
</template>

<script>
import { cloneDeep } from 'lodash';

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
    return {
      filters: cloneDeep(this.modal.config.filters || []),
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
        this.actions.create(newFilter);
      }
    },

    updateFilter(newFilter, index) {
      if (this.actions.update) {
        this.actions.update(newFilter, index);
      }
    },

    deleteFilter(index) {
      if (this.actions.delete) {
        this.actions.delete(index);
      }
    },
  },
};
</script>
