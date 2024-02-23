<template>
  <c-draggable-list-field
    v-field="views"
    :class="{ empty: isViewsEmpty }"
    :group="draggableGroup"
    class="views-panel secondary lighten-1"
  >
    <group-view-panel
      v-for="view in views"
      :key="view._id"
      :view="view"
    />
  </c-draggable-list-field>
</template>

<script>
import GroupViewPanel from './group-view-panel.vue';

export default {
  components: { GroupViewPanel },
  model: {
    prop: 'views',
    event: 'change',
  },
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
    & ::v-deep .panel-item-content {
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
