<template>
  <c-draggable-list-field
    v-field="groups"
    :class="{ empty: isGroupsEmpty }"
    :component-data="{ props: expansionPanelsProps }"
    :group="draggableGroup"
    class="groups-panel secondary"
    component="v-expansion-panels"
  >
    <group-panel
      v-for="(group, groupIndex) in groups"
      :key="group._id"
      :group="group"
    >
      <draggable-group-views
        :views="group.views"
        :put="viewPut"
        :pull="viewPull"
        @change="changeViewsHandler(groupIndex, $event)"
      />
    </group-panel>
  </c-draggable-list-field>
</template>

<script>
import DraggableGroupViews from './draggable-group-views.vue';
import GroupPanel from './group-panel.vue';

export default {
  components: { DraggableGroupViews, GroupPanel },
  model: {
    prop: 'groups',
    event: 'change',
  },
  props: {
    groups: {
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
    viewPut: {
      type: Boolean,
      default: false,
    },
    viewPull: {
      type: [Boolean, String],
      default: false,
    },
  },
  computed: {
    draggableGroupName() {
      return 'groups';
    },

    draggableGroup() {
      return {
        name: this.draggableGroupName,
        put: this.put ? [this.draggableGroupName] : false,
        pull: this.pull ? [this.draggableGroupName] : false,
      };
    },

    expansionPanelsProps() {
      return {
        multiple: true,
        dark: true,
        accordion: true,
        flat: true,
        tile: true,
      };
    },

    isGroupsEmpty() {
      return this.groups.length === 0;
    },
  },
  methods: {
    changeViewsHandler(groupIndex, views) {
      const group = this.groups[groupIndex];

      this.$emit('change:group', groupIndex, { ...group, views });
    },
  },
};
</script>

<style lang="scss" scoped>
  .groups-panel {
    cursor: move;

    & ::v-deep .v-expansion-panel__header {
      cursor: move;
    }

    &.empty {
      &:after {
        content: '';
        display: block;
        height: 48px;
        width: 100%;
        border: 4px dashed #3c5365;
        border-radius: 5px;
        position: relative;
      }
    }
  }
</style>
