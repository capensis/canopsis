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
      vertical-compact,
      @layout-updated="updatedLayout"
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
import { get, omit } from 'lodash';

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
    const layouts = this.getLayouts(this.tab.widgets);

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
    'tab.widgets': function updateLayouts(widgets) {
      this.layouts = this.getLayouts(widgets, true);
    },

    $mq() {
      this.size = this.getGridSizeByMediaQuery();
    },
  },
  methods: {
    updatedLayout() {
      const fields = this.tab.widgets.reduce((acc, { gridParameters }, index) => {
        Object.entries(gridParameters).forEach(([size, gridSettings]) => {
          const params = this.layouts[size][index];

          acc[`widgets.${index}.gridParameters.${size}`] = {
            ...gridSettings,
            ...omit(params, ['i', 'widget', 'moved']),
          };
        });

        return acc;
      }, {});

      const newTab = setSeveralFields(this.tab, fields);

      this.$emit('update:tab', newTab);
    },

    changeFixedHeight(value, index) {
      this.$set(this.layouts[this.size][index], 'fixedHeight', value || false);
    },

    getLayout(widgets, size = WIDGET_GRID_SIZES_KEYS.desktop, updating = false) {
      return widgets.map((widget) => {
        const oldLayout = updating &&
          get(this, ['layouts', size], []).find(({ i }) => i === widget._id);

        const layout = oldLayout ?
          omit(oldLayout, ['i', 'widget']) :
          { ...widget.gridParameters[size] };

        layout.i = widget._id;
        layout.widget = widget;

        return layout;
      });
    },

    getLayouts(widgets, updating = false) {
      return Object.values(WIDGET_GRID_SIZES_KEYS).reduce((acc, size) => {
        acc[size] = this.getLayout(widgets, size, updating);

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
