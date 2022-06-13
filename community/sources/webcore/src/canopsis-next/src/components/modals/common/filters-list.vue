<template lang="pug">
  modal-wrapper(close)
    template(#title="")
      span {{ $t('common.filters') }}
    template(#text="")
      filters-list-component(
        :filters="filters",
        :pending="pending",
        :addable="config.hasAccessToAddFilter",
        :editable="config.hasAccessToEditFilter",
        @add="showCreateFilterModal",
        @edit="showEditFilterModal",
        @delete="showDeleteFilterModal"
      )
</template>

<script>
import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { entitiesWidgetMixin } from '@/mixins/entities/view/widget';
import { entitiesUserPreferenceMixin } from '@/mixins/entities/user-preference';

import FiltersListComponent from '@/components/other/filter/filters-list.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Confirmation modal
 */
export default {
  name: MODALS.filtersList,
  components: { FiltersListComponent, ModalWrapper },
  mixins: [
    modalInnerMixin,
    entitiesWidgetMixin,
    entitiesUserPreferenceMixin,
  ],
  data() {
    return {
      pending: false,
    };
  },
  computed: {
    private() {
      return this.config.private;
    },

    widgetId() {
      return this.config.widgetId;
    },

    widget() {
      return this.getWidgetById(this.widgetId);
    },

    userPreference() {
      return this.getUserPreferenceByWidgetId(this.widgetId);
    },

    filters() {
      return (this.private ? this.userPreference?.filters : this.widget?.filters) ?? [];
    },
  },
  mounted() {
    this.refreshFilters();
  },
  methods: {
    refreshFilters() {
      return this.config.private
        ? this.fetchUserPreference({ id: this.config.widgetId })
        : this.fetchWidget({ id: this.config.widgetId });
    },

    showCreateFilterModal() {
      this.$modals.show({
        name: MODALS.createFilter,
        config: {
          title: this.$t('modals.createFilter.create.title'),
          withTitle: true,
          withAlarm: this.withAlarm,
          withEntity: this.withEntity,
          withPbehavior: this.withPbehavior,
          action: async (newFilter) => {
            await this.createWidgetFilter({
              data: {
                ...newFilter,

                widget: this.widgetId,
                is_private: this.private,
              },
            });

            return this.refreshFilters();
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
          action: async (newFilter) => {
            await this.updateWidgetFilter({
              id: filter._id,
              data: newFilter,
            });

            return this.refreshFilters();
          },
        },
      });
    },

    showDeleteFilterModal(filter) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.removeWidgetFilter({
              id: filter._id,
            });

            return this.refreshFilters(); // TODO: check selected filter (discuss with backend team)
          },
        },
      });
    },
  },
};
</script>
