<template lang="pug">
  div.grid-layout-wrapper
    portal(:to="$constants.PORTALS_NAMES.additionalTopBarItems")
      window-size-field(v-model="size", color="white", light)
    grid-layout(
      ref="gridLayout",
      :layout.sync="layouts[size]",
      :margin="[$constants.WIDGET_GRID_ROW_HEIGHT, $constants.WIDGET_GRID_ROW_HEIGHT]",
      :col-num="$constants.WIDGET_GRID_COLUMNS_COUNT",
      :row-height="$constants.WIDGET_GRID_ROW_HEIGHT",
      :style="layoutStyles",
      is-draggable,
      is-resizable,
      vertical-compact,
      @layout-updated="updatedLayout"
    )
      grid-item(
        v-for="(layoutItem, index) in layouts[size]",
        :key="layoutItem.i",
        :x="layoutItem.x",
        :y="layoutItem.y",
        :w="layoutItem.w",
        :h="layoutItem.h",
        :i="layoutItem.i",
        :fixedHeight="layoutItem.fixedHeight",
        dragAllowFrom=".drag-handler"
      )
        div.wrapper
          div.drag-handler
            v-layout.controls
              v-btn-toggle.mr-2(
                :value="layoutItem.fixedHeight",
                @change="changeFixedHeight(index, $event)"
              )
                v-btn(small, :value="true")
                  v-icon lock
              widget-wrapper-menu(
                :widget="layoutItem.widget",
                :tab="tab",
                :updateTabMethod="updateTabMethod"
              )
          slot(:widget="layoutItem.widget")
</template>

<script>
import { get, omit } from 'lodash';

import {
  WIDGET_GRID_SIZES_KEYS,
  WIDGET_LAYOUT_MAX_WIDTHS,
  MQ_KEYS_TO_WIDGET_GRID_SIZES_KEYS_MAP,
} from '@/constants';

import { setSeveralFields } from '@/helpers/immutable';
import { getWidgetsLayoutBySize } from '@/helpers/grid-layout';

import WidgetWrapperMenu from '@/components/widgets/partials/widget-wrapper-menu.vue';
import WindowSizeField from '@/components/forms/fields/window-size.vue';

export default {
  components: {
    WindowSizeField,
    WidgetWrapperMenu,
  },
  props: {
    tab: {
      type: Object,
      required: true,
    },
    updateTabMethod: {
      type: Function,
      required: true,
    },
  },
  data() {
    const layouts = this.getLayoutsForAllSizes(this.tab.widgets);

    return {
      layouts,
      size: WIDGET_GRID_SIZES_KEYS.desktop,
    };
  },
  computed: {
    layoutStyles() {
      return {
        maxWidth: WIDGET_LAYOUT_MAX_WIDTHS[this.size],
      };
    },
  },
  watch: {
    'tab.widgets': function tabWidgets(widgets) {
      this.layouts = this.getLayoutsForAllSizes(widgets);
    },

    $mq(mq) {
      this.size = MQ_KEYS_TO_WIDGET_GRID_SIZES_KEYS_MAP[mq];
    },
  },
  methods: {
    /**
     * Change fixed height for special layout item
     *
     * @param {boolean} value
     * @param {number} index
     */
    changeFixedHeight(index, value = false) {
      this.$set(this.layouts[this.size][index], 'fixedHeight', value);
    },

    /**
     * Get layouts for all sizes
     *
     * @param {Array} widgets
     * @return {Array}
     */
    getLayoutsForAllSizes(widgets) {
      return Object.values(WIDGET_GRID_SIZES_KEYS).reduce((acc, size) => {
        const oldLayout = get(this, ['layout', size], []);

        acc[size] = getWidgetsLayoutBySize(widgets, oldLayout, size);

        return acc;
      }, {});
    },

    /**
     * Emit 'update:tab' event when layout will be updated
     */
    updatedLayout() {
      const fields = this.tab.widgets.reduce((acc, { gridParameters }, index) => {
        Object.entries(gridParameters).forEach(([size, gridSettings]) => {
          const layoutItem = this.layouts[size][index];

          acc[`widgets.${index}.gridParameters.${size}`] = {
            ...gridSettings,
            ...omit(layoutItem, ['i', 'widget', 'moved']),
          };
        });

        return acc;
      }, {});

      const newTab = setSeveralFields(this.tab, fields);

      this.$emit('update:tab', newTab);
    },
  },
};
</script>

<style lang="scss" scoped>
  .grid-layout-wrapper {
    padding-bottom: 500px;

    & /deep/ .vue-grid-layout {
      margin: auto;
      background-color: rgba(60, 60, 60, .05);
    }

    & /deep/ .vue-grid-item {
      overflow: hidden;
      transition: none !important;

      &:after {
        content: '';
        background-color: #888;
        position: absolute;
        left: 0;
        top: 36px;
        width: 100%;
        height: 100%;
        opacity: .4;
        z-index: 1;
      }

      & > .vue-resizable-handle {
        z-index: 2;
      }
    }
  }

  .wrapper {
    position: relative;
    overflow: hidden;

    .controls {
      position: absolute;
      right: 3px;
      top: 3px;
      z-index: 2;
    }

    &.fixed-height {
      height: 100%;
    }

    .drag-handler {
      position: absolute;
      background-color: rgba(0, 0, 0, .12);
      width: 100%;
      height: 36px;
      transition: .2s ease-out;
      cursor: move;
      z-index: 2;

      &:hover {
        background-color: rgba(0, 0, 0, .15);
      }
    }

    & /deep/ .v-card {
      position: relative;
      min-height: 100%;
    }
  }
</style>
