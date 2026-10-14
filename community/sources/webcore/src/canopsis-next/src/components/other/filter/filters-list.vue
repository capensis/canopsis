<template>
  <div>
    <c-alert v-if="!pending && !filters.length" :type="required ? 'error' : 'info'">
      {{ $t('modals.createFilter.emptyFilters') }}
    </c-alert>

    <c-draggable-list-field
      :model-value="filters"
      :disabled="!editable"
      handle=".action-drag-handler"
      @update:model-value="$emit('update:filters', $event)"
    >
      <template #item="{ element, index }">
        <filter-tile
          :filter="element"
          :editable="editable"
          @edit="$emit('edit', element, index)"
          @delete="$emit('delete', element, index)"
        />
      </template>
    </c-draggable-list-field>

    <v-btn
      v-if="addable"
      class="my-2"
      color="primary"
      variant="outlined"
      @click.prevent="$emit('add', $event)"
    >
      {{ $t('common.addFilter') }}
    </v-btn>
  </div>
</template>

<script>
import FilterTile from './partials/filter-tile.vue';

export default {
  components: { FilterTile },

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
    required: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'filters',
    },
  },

  emits: ['update:filters', 'add', 'edit', 'delete'],
};
</script>
