<template lang="pug">
  div.grid-layout-wrapper
    portal(v-if="editing", :to="$constants.PORTALS_NAMES.additionalTopBarItems")
      window-size-field(v-model="size", color="white", light)
    c-grid-layout(
      v-model="layouts[size]",
      :margin="[$constants.WIDGET_GRID_ROW_HEIGHT, $constants.WIDGET_GRID_ROW_HEIGHT]",
      :columns-count="$constants.WIDGET_GRID_COLUMNS_COUNT",
      :row-height="$constants.WIDGET_GRID_ROW_HEIGHT",
      :style="layoutStyles",
      :disabled="!editing",
      auto-size,
      @input="updatedLayout"
    )
      template(#default="{ bind, on }")
        c-grid-item(
          v-for="(layoutItem, index) in layouts[size]",
          v-bind="bind",
          v-on="on",
          :key="layoutItem.i",
          :x="layoutItem.x",
          :y="layoutItem.y",
          :w="layoutItem.w",
          :h="layoutItem.h",
          :i="layoutItem.i",
          :auto-height="layoutItem.autoHeight",
          :widget="layoutItem.widget",
          resizable
        )
          template(#default="{ on }")
            widget-edit-drag-handler(
              v-on="on",
              :widget="layoutItem.widget",
              :tab="tab",
              :auto-height="layoutItem.autoHeight",
              @toggle="toggleAutoHeight(index)"
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

import { getWidgetsLayoutBySize } from '@/helpers/entities/widget/layout';

import CGridLayout from '@/components/common/grid/c-grid-layout.vue';
import CGridItem from '@/components/common/grid/c-grid-item.vue';
import WidgetWrapperMenu from '@/components/widgets/partials/widget-wrapper-menu.vue';
import WindowSizeField from '@/components/forms/fields/window-size.vue';
import WidgetEditDragHandler from '@/components/widgets/widget-edit-drag-handler.vue';

export default {
  components: {
    WidgetEditDragHandler,
    CGridLayout,
    CGridItem,
    WindowSizeField,
    WidgetWrapperMenu,
  },
  props: {
    tab: {
      type: Object,
      required: true,
    },
    editing: {
      type: Boolean,
      default: false,
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
  }

  .wrapper {
    position: relative;
    overflow: hidden;

    & ::v-deep .v-card {
      position: relative;
      min-height: 100%;
    }
  }
</style>
