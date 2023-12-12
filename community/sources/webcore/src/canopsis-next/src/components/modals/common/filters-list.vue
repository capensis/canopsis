<template lang="pug">
  modal-wrapper(close)
    template(#title="")
      span {{ $t('common.filters') }}
    template(#text="")
      c-progress-overlay(:pending="pending")
      filters-list-component(
        :filters="filters",
        :pending="pending",
        :addable="config.addable",
        :editable="config.editable",
        @input="updateFiltersPositions",
        @add="showCreateFilterModal",
        @edit="showEditFilterModal",
        @delete="showDeleteFilterModal"
      )
</template>

<script>
import { pick } from 'lodash';

import { MODALS } from '@/constants';

import { mapIds } from '@/helpers/array';

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
      pending: true,
      filters: [],
    };
  },
  computed: {
    widgetId() {
      return this.config.widgetId;
    },

    userPreference() {
      return this.getUserPreferenceByWidgetId(this.widgetId);
    },

    modalConfig() {
      return {
        ...pick(this.config, [
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
  watch: {
    'userPreference.filters': function handler(filters) {
      this.setFilters(filters);
    },
  },
  mounted() {
    this.refreshFilters();
  },
  methods: {
    setFilters(filters = []) {
      this.filters = filters;
    },

    async refreshFilters() {
      this.pending = true;

      await this.fetchUserPreference({ id: this.config.widgetId });

      this.pending = false;
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
                is_user_preference: true,
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

    async updateFiltersPositions(filters) {
      const oldFilters = this.filters;

      try {
        this.setFilters(filters);

        await this.updateWidgetFiltersPositions({
          data: mapIds(filters),
        });
      } catch (err) {
        console.error(err);

        this.$popups.error({ text: this.$t('errors.default') });

        this.setFilters(oldFilters);
      }
    },
  },
};
</script>
