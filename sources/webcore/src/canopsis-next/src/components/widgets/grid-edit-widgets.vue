<template lang="pug">
  grid-layout(
    :layout.sync="layout",
    :margin="[$constants.WIDGET_GRID_ROW_HEIGHT, $constants.WIDGET_GRID_ROW_HEIGHT]",
    :col-num="12",
    :row-height="$constants.WIDGET_GRID_ROW_HEIGHT",
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
      @resized="resized",
      :ref="item.i"
    )
      div.wrapper
        div.drag-handler
        v-btn-toggle.lock-icon(
          :value="item.i.fixedHeight",
          @change="changeFixedHeight($event, item.i, index)"
        )
          v-btn(small, :value="true")
            v-icon lock
        widget-wrapper(
          :widget="item.widget",
          :tab="tab",
          :updateTabMethod="updateTabMethod"
        )
</template>

<script>
import { omit } from 'lodash';

import GridItem from '@/components/other/grid/grid-item.vue';
import GridLayout from '@/components/other/grid/grid-layout.vue';
import WidgetWrapper from '@/components/widgets/widget-wrapper.vue';

import { setSeveralFields } from '@/helpers/immutable';

export default {
  components: { GridLayout, WidgetWrapper, GridItem },
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
    const layout = this.tab.widgets.map(widget => ({
      ...widget.gridParameters.desktop,

      i: widget._id,
      widget,
    }));

    return {
      layout,
    };
  },
  beforeDestroy() {
    const fields = this.tab.widgets.reduce((acc, { gridParameters }, index) => {
      acc[`widgets.${index}.gridParameters.desktop`] = {
        ...gridParameters.desktop,
        ...omit(this.layout[index], ['i', 'widget']),
      };

      return acc;
    }, {});

    const newTab = setSeveralFields(this.tab, fields);

    this.updateTabMethod(newTab);
  },
  methods: {
    resized(id, h, previousW, height, width, autoSize) {
      const [widgetLayout] = this.$refs[id];

      if (!autoSize && !widgetLayout.fixedHeight) {
        this.autoSizeWidgetHeight(id);
      }
    },

    changeFixedHeight(value, id, index) {
      if (!value) {
        this.autoSizeWidgetHeight(id);
      }

      this.layout[index].fixedHeight = value;
    },

    autoSizeWidgetHeight(id) {
      const [widgetLayout] = this.$refs[id];

      this.$nextTick(() => widgetLayout.autoSizeHeight());
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

    .lock-icon {
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
      background-color: #BBB;
      position: absolute;
      left: 0;
      top: 0;
      width: 100%;
      height: 36px;
      opacity: .4;
      transition: .2s ease-out;
      cursor: move;
      z-index: 2;

      &:hover {
        opacity: .2;
      }
    }

    & /deep/ .v-card {
      position: relative;
      min-height: 100%;
    }
  }
</style>
