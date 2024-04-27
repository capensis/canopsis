<template lang="pug">
  div.view(:id="`view-tab-${tab._id}`")
    portal(v-if="editing", :to="$constants.PORTALS_NAMES.additionalTopBarItems")
      window-size-field(v-model="size", color="white", light)
    grid-layout(
      v-model="layouts[size]",
      :margin="[$constants.WIDGET_GRID_ROW_HEIGHT, $constants.WIDGET_GRID_ROW_HEIGHT]",
      :columns-count="$constants.WIDGET_GRID_COLUMNS_COUNT",
      :row-height="$constants.WIDGET_GRID_ROW_HEIGHT",
      :style="layoutStyle",
      :disabled="!editing || !visible"
    )
      template(#item="{ on, item }")
        widget-edit-drag-handler(
          v-if="editing",
          v-on="on",
          :widget="item.widget",
          :auto-height="item.autoHeight",
          :tab="tab"
        )
        widget-wrapper(
          :widget="item.widget",
          :tab="tab",
          :kiosk="kiosk",
          :editing="editing",
          :visible="visible"
        )
</template>

<script>
import { isEqual } from 'lodash';

import { WIDGET_GRID_SIZES_KEYS, MQ_KEYS_TO_WIDGET_GRID_SIZES_KEYS_MAP, WIDGET_LAYOUT_MAX_WIDTHS } from '@/constants';

import { widgetsToLayoutsWithCompact, layoutsToWidgetsGrid } from '@/helpers/entities/widget/grid';

import { queryMixin } from '@/mixins/query';
import { activeViewMixin } from '@/mixins/active-view';
import { entitiesWidgetMixin } from '@/mixins/entities/view/widget';

import GridLayout from '@/components/common/grid/grid-layout.vue';
import WidgetWrapper from '@/components/widgets/widget-wrapper.vue';
import WindowSizeField from '@/components/forms/fields/window-size.vue';
import WidgetEditDragHandler from '@/components/widgets/widget-edit-drag-handler.vue';

export default {
  components: {
    GridLayout,
    WidgetWrapper,
    WindowSizeField,
    WidgetEditDragHandler,
  },
  mixins: [
    queryMixin,
    activeViewMixin,
    entitiesWidgetMixin,
  ],
  props: {
    tab: {
      type: Object,
      required: true,
    },
    kiosk: {
      type: Boolean,
      default: false,
    },
    visible: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    const layouts = widgetsToLayoutsWithCompact(this.tab.widgets);

    return {
      layouts,

      size: WIDGET_GRID_SIZES_KEYS.desktop,
      widgetsGrid: layoutsToWidgetsGrid(layouts),
    };
  },
  computed: {
    layoutStyle() {
      return {
        maxWidth: WIDGET_LAYOUT_MAX_WIDTHS[this.size],
      };
    },
  },
  watch: {
    'tab.widgets': function tabWidgets(widgets) {
      this.layouts = widgetsToLayoutsWithCompact(widgets, this.layouts);
    },

    $mq: {
      immediate: true,
      handler(mq) {
        this.size = MQ_KEYS_TO_WIDGET_GRID_SIZES_KEYS_MAP[mq];
      },
    },

    editing() {
      this.size = MQ_KEYS_TO_WIDGET_GRID_SIZES_KEYS_MAP[this.$mq];
    },

    visible: {
      immediate: true,
      handler(visible) {
        if (visible) {
          this.registerEditingOffHandler(this.updatePositions);
        } else {
          this.unregisterEditingOffHandler(this.updatePositions);
        }
      },
    },
  },
  beforeDestroy() {
    this.removeWidgetsQueries();
    this.unregisterEditingOffHandler(this.updatePositions);
  },
  methods: {
    async updatePositions() {
      try {
        const newWidgetsGrid = layoutsToWidgetsGrid(this.layouts);

        if (isEqual(this.widgetsGrid, newWidgetsGrid) || !newWidgetsGrid.length) {
          return;
        }

        this.widgetsGrid = newWidgetsGrid;

        await this.updateWidgetGridPositions({ data: newWidgetsGrid });
        await this.fetchActiveView();
      } catch (err) {
        console.error(err);

        this.$popups.error({ text: this.$t('errors.default') });
      }
    },

    /**
     * Remove queries which was created for all widgets
     */
    removeWidgetsQueries() {
      this.tab.widgets.forEach(({ _id: id }) => this.removeQuery({ id }));
    },
  },
};
</script>

<style lang="scss" scoped>
  .full-screen {
    .hide-on-full-screen {
      display: none;
    }
  }
</style>
