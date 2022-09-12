<template lang="pug">
  div
    v-alert(
      :value="!pending && !filters.length",
      type="info"
    ) {{ $t('modals.createFilter.emptyFilters') }}
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
