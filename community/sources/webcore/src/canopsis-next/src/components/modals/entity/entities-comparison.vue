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
        <entity-comparison-list
          :title="$t('modals.entitiesComparison.foundInCurrent')"
          :initial-query="currentOnlyTableInitialQuery"
        />
        <entity-comparison-list
          :title="$t('modals.entitiesComparison.foundInSuggestion')"
          :initial-query="suggestionOnlyTableInitialQuery"
        />
      </v-layout>
    </template>
  </modal-wrapper>
</template>

<script>
import { computed } from 'vue';

import { PAGINATION_LIMIT } from '@/config';
import { MODALS } from '@/constants';

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
    const { config } = useInnerModal(props);

    const currentPattern = computed(() => JSON.stringify(config.value.currentPattern));
    const suggestionPattern = computed(() => JSON.stringify(config.value.suggestionPattern));

    const currentOnlyTableInitialQuery = computed(() => ({
      sort: [],
      page: 1,
      itemsPerPage: PAGINATION_LIMIT,
      entity_pattern: currentPattern.value,
      negative_entity_pattern: suggestionPattern.value,
    }));

    const suggestionOnlyTableInitialQuery = computed(() => ({
      sort: [],
      page: 1,
      itemsPerPage: PAGINATION_LIMIT,
      entity_pattern: suggestionPattern.value,
      negative_entity_pattern: currentPattern.value,
    }));

    return {
      currentOnlyTableInitialQuery,
      suggestionOnlyTableInitialQuery,
    };
  },
};
</script>
