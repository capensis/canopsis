<template>
  <v-list-item class="filter-tile pa-0">
    <c-layout class="filter-tile__content" align="center">
      <c-layout
        class="filter-tile__title"
        align="center"
      >
        <v-icon
          v-if="editable"
          class="draggable ml-0 mr-3 action-drag-handler"
          icon="drag_indicator"
          size="small"
        />

        <v-list-item-title>{{ filter.title }}</v-list-item-title>
      </c-layout>

      <template v-if="editable">
        <c-layout
          class="filter-tile__actions"
          align="center"
        >
          <c-action-btn type="edit" @click="$emit('edit')" />
          <c-action-btn type="delete" @click="$emit('delete')" />
        </c-layout>
      </template>
    </c-layout>
  </v-list-item>
</template>

<script>
export default {
  props: {
    filter: {
      type: Object,
      default: () => ({}),
    },
    editable: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['edit', 'delete'],
};
</script>

<style lang="scss" scoped>
.filter-tile {
  width: 100%;

  &__content {
    width: 100%;
  }

  &__title {
    flex: 1 1 auto;
    min-width: 0;
  }

  &__actions {
    flex: 0 0 auto;
    margin-left: auto;
  }
}
</style>
