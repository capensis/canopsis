<template>
  <template v-if="!item.parentKey || opened.includes(item.parentKey)">
    <tr>
      <td class="service-dependencies-row__icons">
        <div :style="indentStyle" class="service-dependencies-row__icons-content">
          <v-btn
            v-if="item.children && !item.loadMore && !item.cycle"
            :icon="isOpened ? 'keyboard_arrow_down' : 'keyboard_arrow_right'"
            variant="text"
            size="x-small"
            @click="$emit('toggle', item)"
          />
          <span v-else class="service-dependencies-row__indent" />
          <service-dependencies-expand
            :item="item"
            :pending="pendingByIds[item.parentId]"
            @load="$emit('load', $event)"
            @show="$emit('show', $event)"
          />
        </div>
      </td>
      <td v-for="column in columns" :key="column.value">
        <service-dependencies-entity-cell v-if="item.entity" :item="item" :column="column" />
      </td>
    </tr>
    <service-dependencies-row
      v-for="child in item.children"
      :key="child.key"
      :item="child"
      :columns="columns"
      :level="level + 1"
      :opened="opened"
      :pending-by-ids="pendingByIds"
      @toggle="$emit('toggle', $event)"
      @load="$emit('load', $event)"
      @show="$emit('show', $event)"
    />
  </template>
</template>

<script>
import ServiceDependenciesEntityCell from './service-dependencies-entity-cell.vue';
import ServiceDependenciesExpand from './service-dependencies-expand.vue';

export default {
  name: 'ServiceDependenciesRow',

  components: {
    ServiceDependenciesEntityCell,
    ServiceDependenciesExpand,
  },

  props: {
    item: {
      type: Object,
      required: true,
    },
    columns: {
      type: Array,
      default: () => [],
    },
    level: {
      type: Number,
      default: 0,
    },
    opened: {
      type: Array,
      default: () => [],
    },
    pendingByIds: {
      type: Object,
      default: () => ({}),
    },
  },

  emits: ['load', 'show', 'toggle'],

  computed: {
    isOpened() {
      return this.opened.includes(this.item.key);
    },

    indentStyle() {
      return {
        marginInlineStart: `${this.level * 24}px`,
      };
    },
  },
};
</script>

<style lang="scss" scoped>
.service-dependencies-row {
  &__icons {
    width: 1%;
    white-space: nowrap;
  }

  &__icons-content {
    display: inline-flex;
    align-items: center;
    flex-wrap: nowrap;
    gap: 2px;
  }

  &__indent {
    display: inline-block;
    width: 32px;
    flex: 0 0 32px;
  }
}
</style>
