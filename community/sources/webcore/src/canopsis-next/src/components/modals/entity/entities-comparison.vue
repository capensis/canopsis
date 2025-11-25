<template>
  <modal-wrapper close>
    <template #title="">
      <span>{{ $t('modals.entitiesComparison.title') }}</span>
    </template>
    <template #text="">
      <v-layout class="gap-4" column>
        <c-alert :value="true" type="warning">
          <div v-html="$t('modals.entitiesComparison.infoMessage')" class="pre-wrap" />
        </c-alert>
        <entity-comparison-list :title="$t('modals.entitiesComparison.foundInCurrent', { count: 3 })" />
        <entity-comparison-list :title="$t('modals.entitiesComparison.foundInSuggestion', { count: 5 })" />
      </v-layout>
    </template>
  </modal-wrapper>
</template>

<script>
import { ref, computed } from 'vue';

import { MODALS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';

import EntityComparisonList from '@/components/other/entity/partials/entity-comparison-list.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.entitiesComparison,
  components: {
    ModalWrapper,
    EntityComparisonList,
  },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { config } = useInnerModal(props);

    const openedPanels = ref([0, 1]);
    const loading = ref(false);
    const currentOnlyTableOptions = ref({
      page: 1,
      itemsPerPage: 10,
    });
    const suggestionOnlyTableOptions = ref({
      page: 1,
      itemsPerPage: 10,
    });

    const headers = computed(() => [
      { text: t('modals.entitiesComparison.name'), value: 'name', sortable: false },
      { text: t('modals.entitiesComparison.type'), value: 'type', sortable: false },
      { text: t('common.actionsLabel'), value: 'actions', sortable: false },
    ]);

    const currentOnlyEntities = computed(() => {
      const { currentOnly = [] } = config.value || {};
      return currentOnly;
    });

    const suggestionOnlyEntities = computed(() => {
      const { suggestionOnly = [] } = config.value || {};
      return suggestionOnly;
    });

    const updateCurrentOnlyTablePage = (page) => {
      currentOnlyTableOptions.value.page = page;
    };

    const updateCurrentOnlyTableItemsPerPage = (itemsPerPage) => {
      currentOnlyTableOptions.value.itemsPerPage = itemsPerPage;
    };

    const updateSuggestionOnlyTablePage = (page) => {
      suggestionOnlyTableOptions.value.page = page;
    };

    const updateSuggestionOnlyTableItemsPerPage = (itemsPerPage) => {
      suggestionOnlyTableOptions.value.itemsPerPage = itemsPerPage;
    };

    const handleEdit = () => {
      // TODO: Implement edit action
    };

    const handlePause = () => {
      // TODO: Implement pause action
    };

    const handleInfo = () => {
      // TODO: Implement info action
    };

    return {
      openedPanels,
      loading,
      currentOnlyTableOptions,
      suggestionOnlyTableOptions,
      headers,
      currentOnlyEntities,
      suggestionOnlyEntities,
      updateCurrentOnlyTablePage,
      updateCurrentOnlyTableItemsPerPage,
      updateSuggestionOnlyTablePage,
      updateSuggestionOnlyTableItemsPerPage,
      handleEdit,
      handlePause,
      handleInfo,
    };
  },
};
</script>
