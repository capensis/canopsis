<template>
  <filters-list
    v-field="filters"
    :addable="addable"
    :editable="editable"
    :name="name"
    :required="required"
    @add="showCreateFilterModal"
    @edit="showEditFilterModal"
    @delete="showDeleteFilterModal"
  />
</template>

<script>
import { pick } from 'lodash';

import { MODALS } from '@/constants';

import { uuid } from '@/helpers/uuid';

import { formArrayMixin } from '@/mixins/form';

import FiltersList from '@/components/other/filter/filters-list.vue';

export default {
  components: { FiltersList },
  mixins: [formArrayMixin],
  model: {
    prop: 'filters',
    event: 'input',
  },
  props: {
    widgetId: {
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

        withTitle: true,
      };
    },
  },
  methods: {
    showCreateFilterModal() {
      this.$modals.show({
        name: MODALS.createFilter,
        config: {
          ...this.modalConfig,

          title: this.$t('modals.createFilter.create.title'),
          action: newFilter => this.addItemIntoArray({
            ...newFilter,

            _id: uuid('filter'),
            widget: this.widgetId,
            is_user_preference: false,
          }),
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
          action: newFilter => this.updateItemInArray(index, {
            ...newFilter,

            widget: this.widgetId,
            _id: filter._id,
          }),
        },
      });
    },

    showDeleteFilterModal(filter, index) {
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
