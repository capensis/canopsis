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

    const currentOnlyTableInitialQuery = computed(() => ({
      page: 1,
      itemsPerPage: PAGINATION_LIMIT,
      entity_pattern: config.value.currentPattern,
      negative_entity_pattern: config.value.suggestionPattern,
    }));

    const suggestionOnlyTableInitialQuery = computed(() => ({
      page: 1,
      itemsPerPage: PAGINATION_LIMIT,
      entity_pattern: config.value.suggestionPattern,
      negative_entity_pattern: config.value.currentPattern,
    }));

    return {
      currentOnlyTableInitialQuery,
      suggestionOnlyTableInitialQuery,
    };
  },
};
</script>
