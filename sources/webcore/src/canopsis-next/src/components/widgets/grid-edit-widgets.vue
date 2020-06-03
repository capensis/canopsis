<template lang="pug">
  div.grid-layout-wrapper
    portal(:to="$constants.PORTALS_NAMES.additionalTopBarItems")
      window-size-field(v-model="size", color="white", light)
    grid-layout(
      ref="gridLayout",
      :layout.sync="layouts[size]",
      :margin="[$constants.WIDGET_GRID_ROW_HEIGHT, $constants.WIDGET_GRID_ROW_HEIGHT]",
      :col-num="12",
      :row-height="$constants.WIDGET_GRID_ROW_HEIGHT",
      :style="layoutStyles",
      is-draggable,
      is-resizable,
      vertical-compact
    )
      grid-item(
        v-for="(item, index) in layouts[size]",
        :key="item.i",
        :x="item.x",
        :y="item.y",
        :w="item.w",
        :h="item.h",
        :i="item.i",
        :fixedHeight="item.fixedHeight",
        dragAllowFrom=".drag-handler"
      )
        div.wrapper
          div.drag-handler
            v-layout.controls
              v-btn-toggle.mr-2(
                :value="item.fixedHeight",
                @change="changeFixedHeight($event, index)"
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

import WidgetWrapperMenu from '@/components/widgets/partials/widget-wrapper-menu.vue';
import WindowSizeField from '@/components/forms/fields/window-size.vue';

import { setSeveralFields } from '@/helpers/immutable';
import { WIDGET_GRID_SIZES_KEYS } from '@/constants';
import { MEDIA_QUERIES_BREAKPOINTS } from '@/config';

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
    const layouts = this.getLayouts();

    return {
      layouts,
      size: WIDGET_GRID_SIZES_KEYS.desktop,
    };
  },
  computed: {
    layoutStyles() {
      return {
        maxWidth: this.layoutMaxWidth,
      };
    },

    layoutMaxWidth() {
      return {
        [WIDGET_GRID_SIZES_KEYS.desktop]: '100%',
        [WIDGET_GRID_SIZES_KEYS.tablet]: `${MEDIA_QUERIES_BREAKPOINTS.t}px`,
        [WIDGET_GRID_SIZES_KEYS.mobile]: `${MEDIA_QUERIES_BREAKPOINTS.m}px`,
      }[this.size];
    },
  },
  watch: {
    'tab.widgets': function setLayout() {
      this.layouts = this.getLayouts();
    },
    $mq() {
      this.size = this.getGridSizeByMediaQuery();
    },
  },
  beforeDestroy() {
    this.saveTabWidgets();
  },
  methods: {
    saveTabWidgets() {
      const fields = this.tab.widgets.reduce((acc, { gridParameters }, index) => {
        Object.entries(gridParameters).forEach(([size, gridSettings]) => {
          const params = this.layouts[size][index];

          acc[`widgets.${index}.gridParameters.${size}`] = {
            ...gridSettings,
            ...omit(params, ['i', 'widget']),
          };
        });

        return acc;
      }, {});

      const newTab = setSeveralFields(this.tab, fields);

      this.updateTabMethod(newTab);
    },

    changeFixedHeight(value, index) {
      this.$set(this.layouts[this.size][index], 'fixedHeight', value || false);
    },

    getLayout(size = WIDGET_GRID_SIZES_KEYS.desktop) {
      return this.tab.widgets.map(widget => ({
        ...widget.gridParameters[size],

        i: widget._id,
        widget,
      }));
    },

    getLayouts() {
      return Object.values(WIDGET_GRID_SIZES_KEYS).reduce((acc, size) => {
        acc[size] = this.getLayout(size);

        return acc;
      }, {});
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

    .vue-grid-layout {
      margin: auto;
      background-color: rgba(60, 60, 60, .05);
    }
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
      z-index: 1;
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
