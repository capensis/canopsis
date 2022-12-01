<template lang="pug">
  div
    v-alert(
      :value="!pending && !filters.length",
      type="info"
    ) {{ $t('modals.createFilter.emptyFilters') }}
    c-draggable-list-field(
      v-field="filters",
      :disabled="!editable",
      handle=".action-drag-handler"
    )
      filter-tile(
        v-for="filter in filters",
        :filter="filter",
        :key="filter._id",
        :editable="editable",
        @edit="$emit('edit', filter)",
        @delete="$emit('delete', filter)"
      )
    v-btn.ml-0(
      v-if="addable",
      color="primary",
      outline,
      @click.prevent="$emit('add', $event)"
    ) {{ $t('common.addFilter') }}
</template>

<script>
import { entitiesWidgetMixin } from '@/mixins/entities/view/widget';

import FilterTile from './partials/filter-tile.vue';

export default {
  components: { FilterTile },
  mixins: [entitiesWidgetMixin],
  model: {
    prop: 'filters',
    event: 'input',
  },
  props: {
    filters: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
    },
    addable: {
      type: Boolean,
      default: true,
    },
    editable: {
      type: Boolean,
      default: true,
    },
  },
};
</script>
