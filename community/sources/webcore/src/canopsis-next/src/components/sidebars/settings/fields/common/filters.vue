<template lang="pug">
  widget-settings-item(:title="$t('settings.filters')")
    v-layout(column)
      filter-selector(
        v-if="!hideSelector",
        v-field="value",
        :label="$t('filter.selector.defaultFilter')",
        :filters="filters",
        hide-multiply
      )
      filters-list(
        :filters="filters",
        :addable="addable",
        :editable="editable",
        @input="$emit('update:filters', $event)",
        @add="showCreateFilterModal",
        @edit="showEditFilterModal",
        @delete="showDeleteFilterModal"
      )
</template>

<script>
import { pick } from 'lodash';

import { MODALS } from '@/constants';

import uuid from '@/helpers/uuid';

import { authMixin } from '@/mixins/auth';

import FilterSelector from '@/components/other/filter/filter-selector.vue';
import FiltersList from '@/components/other/filter/filters-list.vue';
import WidgetSettingsItem from '@/components/sidebars/settings/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem, FilterSelector, FiltersList },
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
    withServiceWeather: {
      type: Boolean,
      default: false,
    },
    hideSelector: {
      type: Boolean,
      default: false,
    },
    entityTypes: {
      type: Array,
      required: false,
    },
  },
  computed: {
    modalConfig() {
      return {
        ...pick(this, ['withAlarm', 'withEntity', 'withPbehavior', 'withServiceWeather', 'entityTypes']),

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
          ...this.modalConfig,

          filter,
          title: this.$t('modals.createFilter.edit.title'),
          action: (newFilter) => {
            const preparedNewFilter = {
              ...newFilter,

              widget: this.widgetId,
              _id: filter._id,
            };

            const filters = this.filters.map(item => (item._id === filter._id ? preparedNewFilter : item));

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
