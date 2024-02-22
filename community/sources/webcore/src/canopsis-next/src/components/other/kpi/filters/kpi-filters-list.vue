<template>
  <c-advanced-data-table
    :options="options"
    :items="filters"
    :loading="pending"
    :headers="headers"
    :total-items="totalItems"
    search
    advanced-pagination
    hide-actions
    expand
    @update:options="$emit('update:options', $event)"
  >
    <template #created="{ item }">
      {{ item.created | date }}
    </template>
    <template #updated="{ item }">
      {{ item.updated | date }}
    </template>
    <template #actions="{ item }">
      <v-layout>
        <c-action-btn
          v-if="updatable"
          type="edit"
          @click="$emit('edit', item)"
        />
        <c-action-btn
          v-if="duplicable"
          type="duplicate"
          @click="$emit('duplicate', item)"
        />
        <c-action-btn
          v-if="removable"
          type="delete"
          @click="$emit('remove', item._id)"
        />
      </v-layout>
    </template>
    <template #expand="{ item }">
      <kpi-filters-expand-item :filter="item" />
    </template>
  </c-advanced-data-table>
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
    options: {
      type: Object,
      required: true,
    },
    totalItems: {
      type: Number,
      required: false,
    },
    pending: {
      type: Boolean,
      default: false,
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
        { text: this.$t('common.name'), value: 'name' },
        { text: this.$t('common.created'), value: 'created' },
        { text: this.$t('common.lastModifiedOn'), value: 'updated' },
        { text: this.$t('common.actionsLabel'), value: 'actions', sortable: false },
      ];
    },
  },
};
</script>
