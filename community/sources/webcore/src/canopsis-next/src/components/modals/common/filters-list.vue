<template>
  <modal-wrapper close>
    <template #title="">
      <span>{{ $t('common.filters') }}</span>
    </template>
    <template #text="">
      <c-progress-overlay :pending="combinedPending" />
      <filters-list-component
        :filters="filters"
        :pending="combinedPending"
        :addable="config.addable"
        :editable="config.editable"
        @input="updateFiltersPositions"
        @add="showCreateFilterModal"
        @edit="showEditFilterModal"
        @delete="showDeleteFilterModal"
      />
    </template>
  </modal-wrapper>
</template>

<script>
import { ref, computed, watch, onMounted } from 'vue';
import { pick } from 'lodash';

import { MODALS } from '@/constants';

import { mapIds } from '@/helpers/array';

import { useInnerModal, useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useI18n } from '@/hooks/i18n';
import { usePendingHandler } from '@/hooks/query/pending';
import { useWidget } from '@/hooks/store/modules/widget';
import { useUserPreference } from '@/hooks/store/modules/user-preference';
import { usePatternsFields, usePatternsFieldsFetching } from '@/hooks/store/modules/patterns-fields';

import FiltersListComponent from '@/components/other/filter/filters-list.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Confirmation modal
 */
export default {
  name: MODALS.filtersList,
  components: { FiltersListComponent, ModalWrapper },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const filters = ref([]);

    const { t } = useI18n();
    const modals = useModals();
    const popups = usePopups();
    const { config } = useInnerModal(props);
    const {
      createWidgetFilter,
      updateWidgetFilter,
      removeWidgetFilter,
      updateWidgetFiltersPositions,
    } = useWidget();
    const {
      getUserPreferenceByWidgetId,
      fetchUserPreferenceItem,
    } = useUserPreference();

    const widgetId = computed(() => config.value.widgetId);

    const userPreference = computed(() => getUserPreferenceByWidgetId.value(widgetId.value));

    const { fetchWidgetFilterPatternFields } = usePatternsFields();
    const {
      pending: patternsFieldsPending,
      alarmAttributes,
      entityAttributes,
      pbehaviorAttributes,
      weatherServiceAttributes,
      fetchPatternsFields,
    } = usePatternsFieldsFetching(() => fetchWidgetFilterPatternFields(), true);

    const modalConfig = computed(() => ({
      ...pick(config.value, [
        'withAlarm',
        'withEntity',
        'withPbehavior',
        'withServiceWeather',
        'entityTypes',
        'entityCountersType',
      ]),

      withTitle: true,
      alarmAttributes: alarmAttributes.value,
      entityAttributes: entityAttributes.value,
      pbehaviorAttributes: pbehaviorAttributes.value,
      weatherServiceAttributes: weatherServiceAttributes.value,
    }));

    /**
     * Sets the filters array with new filter values
     *
     * @param {Array} [newFilters=[]] - Array of filter objects to set
     */
    const setFilters = (newFilters = []) => filters.value = newFilters;

    /**
     * Refreshes filters by fetching user preference data
     */
    const {
      pending,
      handler: refreshFilters,
    } = usePendingHandler(async () => Promise.all([
      fetchUserPreferenceItem({ id: config.value.widgetId }),
      fetchPatternsFields(),
    ]), true);

    const combinedPending = computed(() => pending.value || patternsFieldsPending.value);

    /**
     * Shows modal for creating a new filter
     */
    const showCreateFilterModal = () => modals.show({
      name: MODALS.createFilter,
      config: {
        ...modalConfig.value,

        title: t('modals.createFilter.create.title'),
        corporate: true,
        action: async (newFilter) => {
          await createWidgetFilter({
            data: {
              ...newFilter,

              widget: widgetId.value,
              is_user_preference: true,
            },
          });

          return refreshFilters();
        },
      },
    });

    /**
     * Shows modal for editing an existing filter
     *
     * @param {Object} filter - Filter object to edit
     */
    const showEditFilterModal = filter => modals.show({
      name: MODALS.createFilter,
      config: {
        ...modalConfig.value,

        filter,
        title: t('modals.createFilter.edit.title'),
        corporate: true,
        action: async (newFilter) => {
          await updateWidgetFilter({
            id: filter._id,
            data: {
              ...newFilter,

              widget: widgetId.value,
            },
          });

          return refreshFilters();
        },
      },
    });

    /**
     * Shows confirmation modal for deleting a filter
     *
     * @param {Object} filter - Filter object to delete
     */
    const showDeleteFilterModal = filter => modals.show({
      name: MODALS.confirmation,
      config: {
        action: async () => {
          await removeWidgetFilter({
            id: filter._id,
          });

          return refreshFilters();
        },
      },
    });

    /**
     * Updates the positions of filters and refreshes user preference
     *
     * @param {Array} newFilters - Array of filters with updated positions
     */
    const updateFiltersPositions = async (newFilters) => {
      const oldFilters = filters.value;

      try {
        setFilters(newFilters);

        await updateWidgetFiltersPositions({
          data: mapIds(newFilters),
        });

        await fetchUserPreferenceItem({ id: config.value.widgetId });
      } catch (err) {
        console.error(err);

        popups.error({ text: t('errors.default') });

        setFilters(oldFilters);
      }
    };

    watch(
      () => userPreference.value?.filters,
      (newFilters) => {
        setFilters(newFilters);
      },
    );

    onMounted(refreshFilters);

    return {
      pending: combinedPending,
      filters,
      config,
      modalConfig,
      setFilters,
      refreshFilters,
      showCreateFilterModal,
      showEditFilterModal,
      showDeleteFilterModal,
      updateFiltersPositions,
    };
  },
};
</script>
