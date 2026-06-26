<template>
  <c-draggable-list-field
    :model-value="views"
    :class="{ empty: isViewsEmpty }"
    :group="draggableGroup"
    class="views-panel"
    @update:model-value="$emit('change', $event)"
  >
    <template #item="{ element: view }">
      <group-view-panel :view="view"
      />
    </template>
  </c-draggable-list-field>
</template>

<script>
import GroupViewPanel from './group-view-panel.vue';

export default {
  components: { GroupViewPanel },
  props: {
    views: {
      type: Array,
      required: true,
    },
    put: {
      type: Boolean,
      default: false,
    },
    pull: {
      type: [Boolean, String],
      default: false,
    },
  },

  emits: ['change'],

  computed: {
    isViewsEmpty() {
      return this.views.length === 0;
    },

    draggableGroupName() {
      return 'views';
    },

    draggableGroup() {
      return {
        name: this.draggableGroupName,
        put: this.put ? [this.draggableGroupName] : false,
        pull: this.pull ? [this.draggableGroupName] : false,
      };
    },
  },
};
</script>

<style lang="scss" scoped>
.views-panel {
  background-color: var(--v-secondary-lighten1);

  :deep(.panel-item-content) {
    cursor: move;
  }

  &.empty {
    &:after {
      content: '';
      display: block;
      height: 48px;
      border: 4px dashed #4f6479;
      border-radius: 5px;
      position: relative;
    }
  }
}
</style>
