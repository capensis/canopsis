<template lang="pug">
  c-advanced-data-table.white(
    :pagination="pagination",
    :items="filters",
    :loading="pending",
    :headers="headers",
    :total-items="totalItems",
    search,
    advanced-pagination,
    hide-actions,
    expand,
    @update:pagination="$emit('update:pagination', $event)"
  )
    template(#created="props") {{ props.item.created | date }}
    template(#updated="props") {{ props.item.updated | date }}
    template(#actions="props")
      v-layout(row)
        c-action-btn(
          v-if="updatable",
          type="edit",
          @click="$emit('edit', props.item)"
        )
        c-action-btn(
          v-if="duplicable",
          type="duplicate",
          @click="$emit('duplicate', props.item)"
        )
        c-action-btn(
          v-if="removable",
          type="delete",
          @click="$emit('remove', props.item._id)"
        )
    template(#expand="props")
      kpi-filters-expand-item(:filter="props.item")
</template>

<script>
import KpiFiltersExpandItem from './partials/kpi-filters-expand-item.vue';

export default {
  components: { KpiFiltersExpandItem },
  props: {
    filters: {
      type: Array,
      required: true,
    },
    pagination: {
      type: Object,
      required: true,
    },
    totalItems: {
      type: Number,
      required: false,
    },
    pending: {
      type: Boolean,
      required: true,
    },
    removable: {
      type: Boolean,
      default: false,
    },
    updatable: {
      type: Boolean,
      default: false,
    },
    duplicable: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    headers() {
      return [
        { text: this.$t('common.title'), value: 'name' },
        { text: this.$t('common.created'), value: 'created' },
        { text: this.$t('common.lastModifiedOn'), value: 'updated' },
        { text: this.$t('common.actionsLabel'), value: 'actions', sortable: false },
      ];
    },
  },
};
</script>
