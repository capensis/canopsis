<template>
  <filters-list
    :filters="filters"
    :addable="addable"
    :editable="editable"
    :name="name"
    :required="required"
    @update:filters="$emit('update:filters', $event)"
    @add="showCreateFilterModal"
    @edit="showEditFilterModal"
    @delete="showDeleteFilterModal"
  />
</template>

<script>
import { pick } from 'lodash';

import { MODALS } from '@/constants';

import { uid } from '@/helpers/uid';

import FiltersList from '@/components/other/filter/filters-list.vue';

export default {
  components: { FiltersList },

  props: {
    widgetId: {
      type: String,
      required: false,
    },
    widgetType: {
      type: String,
      required: false,
    },
    filters: {
      type: Array,
      default: () => [],
    },
    addable: {
      type: Boolean,
      default: false,
    },
    editable: {
      type: Boolean,
      default: false,
    },
    withAlarm: {
      type: Boolean,
      default: false,
    },
    withEntity: {
      type: Boolean,
      default: false,
    },
    withPbehavior: {
      type: Boolean,
      default: false,
    },
    withServiceWeather: {
      type: Boolean,
      default: false,
    },
    entityTypes: {
      type: Array,
      required: false,
    },
    entityCountersType: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'filters',
    },
    required: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:filters'],

  computed: {
    modalConfig() {
      return {
        ...pick(this, [
          'withAlarm',
          'withEntity',
          'withPbehavior',
          'withServiceWeather',
          'entityTypes',
        'entityCountersType',
      ]),

        widgetType: this.widgetType,
        withTitle: true,
        withoutLink: true,
      };
    },
  },

  methods: {
    updateFilters(filters) {
      this.$emit('update:filters', filters);
    },

    showCreateFilterModal() {
      this.$modals.show({
        name: MODALS.createFilter,
        config: {
          ...this.modalConfig,

          title: this.$t('modals.createFilter.create.title'),
          action: newFilter => this.updateFilters([
            ...this.filters,
            {
              ...newFilter,
              _id: uid('filter'),
              widget: this.widgetId,
              is_user_preference: false,
            },
          ]),
        },
      });
    },

    showEditFilterModal(filter, index) {
      this.$modals.show({
        name: MODALS.createFilter,
        config: {
          ...this.modalConfig,

          filter,
          title: this.$t('modals.createFilter.edit.title'),
          action: (newFilter) => {
            const filters = [...this.filters];

            filters[index] = {
              ...newFilter,
              widget: this.widgetId,
              _id: filter._id,
            };

            this.updateFilters(filters);
          },
        },
      });
    },

    showDeleteFilterModal(filter, index) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => this.updateFilters(this.filters.filter((_, filterIndex) => filterIndex !== index)),
        },
      });
    },
  },
};
</script>
