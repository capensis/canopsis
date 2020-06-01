<template lang="pug">
  div
    portal(:to="$constants.PORTALS_NAMES.additionalTopBarItems")
      window-size-field(v-model="size")
    grid-layout(
      :layout.sync="layout",
      :margin="[$constants.WIDGET_GRID_ROW_HEIGHT, $constants.WIDGET_GRID_ROW_HEIGHT]",
      :col-num="12",
      :row-height="$constants.WIDGET_GRID_ROW_HEIGHT",
      :style="windowSizeStyles",
      is-draggable,
      is-resizable,
      vertical-compact
    )
      grid-item(
        v-for="(item, index) in layout",
        :key="item.i",
        :x="item.x",
        :y="item.y",
        :w="item.w",
        :h="item.h",
        :fixedHeight="item.fixedHeight",
        :i="item.i",
        dragAllowFrom=".drag-handler",
        :ref="item.i"
      )
        div.wrapper
          div.drag-handler
            v-layout.controls
              v-btn-toggle(
                :value="item.fixedHeight",
                @change="changeFixedHeight($event, item.i, index)"
              )
                v-btn(small, :value="true")
                  v-icon lock
              widget-wrapper-menu(
                :widget="item.widget",
                :tab="tab",
                :updateTabMethod="updateTabMethod"
              )
          slot(:widget="item.widget")
</template>

<script>
import { omit } from 'lodash';

import GridItem from '@/components/other/grid/grid-item.vue';
import GridLayout from '@/components/other/grid/grid-layout.vue';
import WidgetWrapperMenu from '@/components/widgets/partials/widget-wrapper-menu.vue';
import WindowSizeField from '@/components/forms/fields/window-size.vue';

import { setSeveralFields } from '@/helpers/immutable';
import { WIDGET_GRID_SIZES_KEYS } from '@/constants';
import { MEDIA_QUERIES_BREAKPOINTS } from '@/config';

export default {
  components: {
    WindowSizeField,
    WidgetWrapperMenu,
    GridLayout,
    GridItem,
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
    const layout = this.getLayout();

    return {
      layout,
      size: WIDGET_GRID_SIZES_KEYS.desktop,
    };
  },
  computed: {
    windowSizeStyles() {
      return {
        maxWidth: `${this.windowMaxWidth}px`,
      };
    },
    windowMaxWidth() {
      return {
        [WIDGET_GRID_SIZES_KEYS.desktop]: MEDIA_QUERIES_BREAKPOINTS.l,
        [WIDGET_GRID_SIZES_KEYS.tablet]: MEDIA_QUERIES_BREAKPOINTS.t,
        [WIDGET_GRID_SIZES_KEYS.mobile]: MEDIA_QUERIES_BREAKPOINTS.m,
      }[this.size];
    },
  },
  watch: {
    'tab.widgets': function setLayout() {
      this.layout = this.getLayout(this.size);
    },
    size(size, oldSize) {
      this.layout = this.getLayout(size);
      this.saveTabWidgets(oldSize);
    },
    $mq() {
      this.size = this.getGridSizeByMediaQuery();
    },
  },
  beforeDestroy() {
    this.saveTabWidgets();
  },
  methods: {
    saveTabWidgets(size = WIDGET_GRID_SIZES_KEYS.desktop) {
      const fields = this.tab.widgets.reduce((acc, { gridParameters }, index) => {
        const params = this.layout[index];
        const gridSettings = gridParameters[size];

        acc[`widgets.${index}.gridParameters.${size}`] = {
          ...gridSettings,
          ...omit(params, ['i', 'widget']),
        };

        return acc;
      }, {});

      const newTab = setSeveralFields(this.tab, fields);

      this.updateTabMethod(newTab);
    },

    changeFixedHeight(value, id, index) {
      this.layout[index].fixedHeight = value;
    },

    getLayout(size = WIDGET_GRID_SIZES_KEYS.desktop) {
      return this.tab.widgets.map(widget => ({
        ...widget.gridParameters[size],

        i: widget._id,
        widget,
      }));
    },

    getGridSizeByMediaQuery() {
      return {
        xl: WIDGET_GRID_SIZES_KEYS.desktop,
        l: WIDGET_GRID_SIZES_KEYS.desktop,
        t: WIDGET_GRID_SIZES_KEYS.tablet,
        m: WIDGET_GRID_SIZES_KEYS.mobile,
      }[this.$mq];
    },
  },
};
</script>

<style lang="scss">
  .grid-layout-wrapper {
    padding-bottom: 500px;
  }
  .grid-wrapper {
    display: grid;
    grid-gap: 10px;
    grid-template-columns: repeat(12, calc(80% / 12));
  }
  .vue-grid-item {
    overflow: hidden;

    &:after {
      content: '';
      background-color: #888;
      position: absolute;
      left: 0;
      top: 36px;
      width: 100%;
      height: 100%;
      opacity: .4;
    }

    & > .vue-resizable-handle {
      z-index: 2;
    }
  }

  .wrapper {
    position: relative;
    overflow-y: auto;

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
      content: '';
      background-color: rgba(187, 187, 187, 0.4);
      width: 100%;
      height: 36px;
      transition: .2s ease-out;
      cursor: move;
      z-index: 2;

      &:hover {
        background-color: rgba(187, 187, 187, 0.2);
      }
    }

    & /deep/ .v-card {
      position: relative;
      min-height: 100%;
    }
  }
</style>
