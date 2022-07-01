<template lang="pug">
  v-list-group
    template(#activator="")
      v-list-tile {{ $t('settings.filters') }}
    v-container
      v-layout(column)
        filter-selector(
          v-if="!hideSelector",
          v-field="value",
          :label="$t('filterSelector.defaultFilter')",
          :filters="filters"
        )
        filters-list(
          :filters="filters",
          :addable="addable",
          :editable="editable",
          @add="showCreateFilterModal",
          @edit="showEditFilterModal",
          @delete="showDeleteFilterModal"
        )
</template>

<script>
import { MODALS } from '@/constants';

import uuid from '@/helpers/uuid';

import { authMixin } from '@/mixins/auth';

import FilterSelector from '@/components/other/filter/filter-selector.vue';
import FiltersList from '@/components/other/filter/filters-list.vue';

export default {
  components: { FilterSelector, FiltersList },
  mixins: [authMixin],
  props: {
    widgetId: {
      type: String,
      required: false,
    },
    filters: {
      type: Array,
      default: () => [],
    },
    value: {
      type: String,
      default: null,
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
    hideSelector: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    showCreateFilterModal() {
      this.$modals.show({
        name: MODALS.createFilter,
        config: {
          title: this.$t('modals.createFilter.create.title'),
          withTitle: true,
          withAlarm: this.withAlarm,
          withEntity: this.withEntity,
          withPbehavior: this.withPbehavior,
          action: (newFilter) => {
            const filters = [
              ...this.filters,

              {
                ...newFilter,

                _id: uuid('filter'),
                widget: this.widgetId,
                is_private: false,
              },
            ];

            this.$emit('update:filters', filters);
          },
        },
      });
    },

    showEditFilterModal(filter) {
      this.$modals.show({
        name: MODALS.createFilter,
        config: {
          filter,

          title: this.$t('modals.createFilter.edit.title'),
          withTitle: true,
          withAlarm: this.withAlarm,
          withEntity: this.withEntity,
          withPbehavior: this.withPbehavior,
          action: (newFilter) => {
            const filters = this.filters.map(item => (item._id === filter._id ? newFilter : item));

            this.$emit('update:filters', filters);
          },
        },
      });
    },

    showDeleteFilterModal(filter) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => {
            const filters = this.filters.filter(({ _id: id }) => id !== filter._id);

            this.$emit('update:filters', filters);

            if (this.value === filter._id) {
              this.$emit('input', null);
            }
          },
        },
      });
    },
  },
};
</script>
