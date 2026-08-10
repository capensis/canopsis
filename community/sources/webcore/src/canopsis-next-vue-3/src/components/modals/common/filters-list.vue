<template>
  <c-modal-wrapper close>
    <template #title>
      <span>{{ t('common.filters') }}</span>
    </template>

    <template #text>
      <c-progress-overlay :pending="pending" />

      <filters-list-component
        :filters="filters"
        :pending="pending"
        :addable="config.addable"
        :editable="config.editable"
        @update:filters="updateFiltersPositions"
        @add="showCreateFilterModal"
        @edit="showEditFilterModal"
        @delete="showDeleteFilterModal"
      />
    </template>
  </c-modal-wrapper>
</template>

<script>
import {
  computed,
  onMounted,
  ref,
  watch,
} from 'vue';
import { pick } from 'lodash';

import { MODALS } from '@/constants';

import { mapIds } from '@/helpers/array';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePendingHandler } from '@/hooks/query/pending';
import { usePopups } from '@/hooks/popups';
import { usePatternsFields, usePatternsFieldsFetching } from '@/hooks/store/modules/patterns-fields';
import { useUserPreference } from '@/hooks/store/modules/user-preference';
import { useWidget } from '@/hooks/store/modules/widget';

import FiltersListComponent from '@/components/other/filter/filters-list.vue';

export default {
  name: MODALS.filtersList,

  components: { FiltersListComponent },

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
    const config = computed(() => props.modal.config ?? {});
    const {
      createWidgetFilter,
      updateWidgetFilter,
      removeWidgetFilter,
      updateWidgetFiltersPositions,
    } = useWidget();
    const {
      getUserPreferenceByWidgetId,
      fetchUserPreference,
    } = useUserPreference();

    const widgetId = computed(() => config.value.widget?._id);
    const widgetType = computed(() => config.value.widget?.type);
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
      widgetType: widgetType.value,
    }));

    /**
     * Sets the filters array with new filter values.
     *
     * @param {Array} [newFilters=[]] - Array of filter objects to set.
     */
    const setFilters = (newFilters = []) => {
      filters.value = newFilters;
    };

    /**
     * Refreshes filters by fetching user preference data.
     */
    const {
      pending,
      handler: refreshFilters,
    } = usePendingHandler(async () => Promise.all([
      fetchUserPreference({ id: widgetId.value }),
      fetchPatternsFields(),
    ]), true);

    const combinedPending = computed(() => pending.value || patternsFieldsPending.value);

    /**
     * Shows modal for creating a new filter.
     */
    const showCreateFilterModal = () => modals.show({
      name: MODALS.createFilter,
      config: {
        ...modalConfig.value,

        title: t('modals.createFilter.create.title'),
        corporate: true,
        action: async (newFilter) => {
          const result = await createWidgetFilter({
            data: {
              ...newFilter,

              widget: widgetId.value,
              is_user_preference: true,
            },
          });

          refreshFilters();

          return result;
        },
      },
    });

    /**
     * Shows modal for editing an existing filter.
     *
     * @param {Object} filter - Filter object to edit.
     */
    const showEditFilterModal = filter => modals.show({
      name: MODALS.createFilter,
      config: {
        ...modalConfig.value,

        filter,
        title: t('modals.createFilter.edit.title'),
        corporate: true,
        action: async (newFilter) => {
          const result = await updateWidgetFilter({
            id: filter._id,
            data: {
              ...newFilter,

              widget: widgetId.value,
            },
          });

          refreshFilters();

          return result;
        },
      },
    });

    /**
     * Shows confirmation modal for deleting a filter.
     *
     * @param {Object} filter - Filter object to delete.
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
     * Updates the positions of filters and refreshes user preference.
     *
     * @param {Array} newFilters - Array of filters with updated positions.
     */
    const updateFiltersPositions = async (newFilters) => {
      const oldFilters = filters.value;

      try {
        setFilters(newFilters);

        await updateWidgetFiltersPositions({
          data: mapIds(newFilters),
        });

        await fetchUserPreference({ id: widgetId.value });
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
      t,
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
