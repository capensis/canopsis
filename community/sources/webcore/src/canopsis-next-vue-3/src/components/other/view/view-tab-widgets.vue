<template>
  <div :id="`view-tab-${tab._id}`" class="view">
    <Teleport
      v-if="activeViewEditing && visible && additionalTopBarItemsTarget"
      :to="additionalTopBarItemsTarget"
    >
      <window-size-field v-model="size" color="white" />
    </Teleport>

    <the-mass-actions-panels v-if="!isDefaultMode" fixed />

    <grid-layout
      v-model="layouts[size]"
      :margin="[WIDGET_GRID_ROW_HEIGHT, WIDGET_GRID_ROW_HEIGHT]"
      :columns-count="WIDGET_GRID_COLUMNS_COUNT"
      :row-height="WIDGET_GRID_ROW_HEIGHT"
      :style="layoutStyle"
      :disabled="!activeViewEditing || !visible"
    >
      <template #item="{ on, item }">
        <widget-edit-drag-handler
          v-if="activeViewEditing"
          v-bind="on"
          :widget="item.widget"
          :auto-height="item.autoHeight"
          :tab="tab"
        />

        <widget-wrapper :widget="item.widget" :tab="tab" :visible="visible" />
      </template>
    </grid-layout>
  </div>
</template>

<script>
import { isEqual } from 'lodash';
import { mapActions } from 'pinia';

import {
  TELEPORTS_NAMES,
  WIDGET_GRID_COLUMNS_COUNT,
  WIDGET_GRID_ROW_HEIGHT,
  WIDGET_GRID_SIZES_KEYS,
  DISPLAY_NAME_TO_WIDGET_GRID_SIZES_KEYS_MAP,
  WIDGET_LAYOUT_MAX_WIDTHS,
  VIEW_SCREEN_MODES,
} from '@/constants';

import { useWidgetStore } from '@/stores/modules/widget';

import { widgetsToLayoutsWithCompact, layoutsToWidgetsGrid } from '@/helpers/entities/widget/grid';

import { activeViewMixin } from '@/mixins/active-view';

import TheMassActionsPanels from '@/components/common/actions-panel/the-mass-actions-panels.vue';
import GridLayout from '@/components/common/grid/grid-layout.vue';
import WindowSizeField from '@/components/forms/fields/window-size.vue';
import WidgetEditDragHandler from '@/components/widgets/widget-edit-drag-handler.vue';
import WidgetWrapper from '@/components/widgets/widget-wrapper.vue';

export default {
  components: {
    GridLayout,
    TheMassActionsPanels,
    WidgetEditDragHandler,
    WidgetWrapper,
    WindowSizeField,
  },
  mixins: [activeViewMixin],
  props: {
    tab: {
      type: Object,
      required: true,
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
      additionalTopBarItemsTarget: `[data-teleport-target="${TELEPORTS_NAMES.additionalTopBarItems}"]`,
      WIDGET_GRID_COLUMNS_COUNT,
      WIDGET_GRID_ROW_HEIGHT,
    };
  },
  computed: {
    layoutStyle() {
      return {
        maxWidth: WIDGET_LAYOUT_MAX_WIDTHS[this.size],
      };
    },

    displayName() {
      return this.$vuetify.display.name;
    },

    displayGridSize() {
      return DISPLAY_NAME_TO_WIDGET_GRID_SIZES_KEYS_MAP[this.displayName] ?? WIDGET_GRID_SIZES_KEYS.desktop;
    },

    isDefaultMode() {
      return this.activeViewScreenMode === VIEW_SCREEN_MODES.default;
    },
  },
  watch: {
    'tab.widgets': function tabWidgets(widgets) {
      this.layouts = widgetsToLayoutsWithCompact(widgets, this.layouts);
    },

    displayGridSize: {
      immediate: true,
      handler(size) {
        this.size = size;
      },
    },

    activeViewEditing() {
      this.size = this.displayGridSize;
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
  beforeUnmount() {
    this.unregisterEditingOffHandler(this.updatePositions);
  },
  methods: {
    ...mapActions(useWidgetStore, {
      updateWidgetGridPositions: 'updateGridPositions',
    }),

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
