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
import { pick } from 'lodash';

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
    widgetId() {
      return this.config.widgetId;
    },

    userPreference() {
      return this.getUserPreferenceByWidgetId(this.widgetId);
    },

    filters() {
      return this.userPreference?.filters ?? [];
    },

    modalConfig() {
      return {
        ...pick(this.config, ['withAlarm', 'withEntity', 'withPbehavior', 'withServiceWeather', 'entityTypes']),

        withTitle: true,
      };
    },
  },
  mounted() {
    this.refreshFilters();
  },
  methods: {
    refreshFilters() {
      return this.fetchUserPreference({ id: this.config.widgetId });
    },

    showCreateFilterModal() {
      this.$modals.show({
        name: MODALS.createFilter,
        config: {
          ...this.modalConfig,

          title: this.$t('modals.createFilter.create.title'),
          corporate: true,
          action: async (newFilter) => {
            await this.createWidgetFilter({
              data: {
                ...newFilter,

                widget: this.widgetId,
                is_private: true,
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
          ...this.modalConfig,

          filter,
          title: this.$t('modals.createFilter.edit.title'),
          corporate: true,
          action: async (newFilter) => {
            await this.updateWidgetFilter({
              id: filter._id,
              data: {
                ...newFilter,

                widget: this.widgetId,
              },
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

            return this.refreshFilters();
          },
        },
      });
    },
  },
};
</script>
