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
        :auto-height="layoutItem.autoHeight",
        drag-allow-from=".drag-handler"
      )
        div.wrapper
          div.drag-handler
            v-layout.controls
              v-tooltip(bottom)
                template(#activator="{ on }")
                  v-btn.ma-0.mr-1(
                    v-on="on",
                    :color="layoutItem.autoHeight ? 'grey lighten-1' : 'transparent'",
                    icon,
                    small,
                    @click="toggleAutoHeight(index)"
                  )
                    v-icon(
                      :color="layoutItem.autoHeight ? 'black' : 'grey darken-1'",
                      small
                    ) lock
                span {{ $t('view.autoHeightButton') }}
              widget-wrapper-menu(
                :widget="layoutItem.widget",
                :tab="tab"
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

      this.updatedLayout();
    },

    $mq: {
      immediate: true,
      handler(mq) {
        this.size = MQ_KEYS_TO_WIDGET_GRID_SIZES_KEYS_MAP[mq];
      },
    },
  },
  methods: {
    /**
     * Toggle auto height flag for special layout item
     *
     * @param {number} index
     */
    toggleAutoHeight(index) {
      this.$set(this.layouts[this.size][index], 'autoHeight', !this.layouts[this.size][index].autoHeight);
    },

    /**
     * Get layouts for all sizes
     *
     * @param {Array} widgets
     * @return {Array}
     */
    getLayoutsForAllSizes(widgets) {
      return Object.values(WIDGET_GRID_SIZES_KEYS).reduce((acc, size) => {
        const oldLayout = get(this, ['layouts', size], []);

        acc[size] = getWidgetsLayoutBySize(widgets, oldLayout, size);

        return acc;
      }, {});
    },

    /**
     * Emit 'update:widgets-grid' event when layout will be updated
     */
    updatedLayout() {
      const widgetsGrid = Object.entries(this.layouts).reduce((acc, [size, layout]) => {
        layout.forEach((layoutItem) => {
          if (!acc[layoutItem.i]) {
            acc[layoutItem.i] = {};
          }

          acc[layoutItem.i][size] = omit(layoutItem, ['i', 'widget', 'moved']);
        });

        return acc;
      }, {});

      this.$emit('update:widgets-grid', widgetsGrid);
    },
  },
};
</script>

<style lang="scss" scoped>
  .grid-layout-wrapper {
    padding-bottom: 500px;

    & ::v-deep .vue-grid-layout {
      margin: auto;
      background-color: rgba(60, 60, 60, .05);
    }

    & ::v-deep .vue-grid-item {
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
        z-index: 2;
      }

      & > .vue-resizable-handle {
        z-index: 3;
      }
    }
  }

  .wrapper {
    position: relative;
    overflow: hidden;

    .controls {
      position: absolute;
      right: 4px;
      top: 4px;
      z-index: 2;
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

    & ::v-deep .v-card {
      position: relative;
      min-height: 100%;
    }
  }
</style>
