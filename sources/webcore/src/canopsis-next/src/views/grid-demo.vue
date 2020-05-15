<template lang="pug">
  div.grid-layout-wrapper
    v-layout
      v-switch(v-model="editMode", label="Edit grid")
    grid-layout(
      v-if="editMode",
      :layout.sync="layout",
      :col-num="12",
      :row-height="20",
      :is-draggable="draggable",
      :is-resizable="resizable",
      :vertical-compact="true",
      :responsive="true"
    )
      grid-item(
        v-for="item in layout",
        :key="item.i._id",
        :x="item.x",
        :y="item.y",
        :w="item.w",
        :h="item.h",
        :i="item.i",
        dragAllowFrom=".drag-handler",
        @resized="resized"
      )
        div.wrapper(:class="{ 'fixed-height': item.i.fixedHeight }")
          div.drag-handler
          v-btn-toggle.lock-icon(v-model="item.i.fixedHeight", @change="change($event, item.i.index)")
            v-btn(small, :value="true")
              v-icon lock
          widget-wrapper(
            :widget="item.i.widget",
            :tab="{ _id: '123' }",
            :isEditingMode="true",
            :row="{}",
            :updateTabMethod="() => {}",
            ref="widgets"
          )
    .grid-wrapper(v-else)
      div.grid-item(
        v-for="item in layout",
        :key="item.i._id",
        :style="getWidgetFlexStyle(item)"
      )
        widget-wrapper(
          :widget="item.i.widget",
          :tab="{ _id: '123' }",
          :row="{}",
          :updateTabMethod="() => {}",
          ref="widgets"
        )
</template>

<script>
import { GridLayout } from 'vue-grid-layout';

import { WIDGET_TYPES } from '@/constants';

import { generateWidgetByType } from '@/helpers/entities';

import WidgetWrapper from '@/components/widgets/widget-wrapper.vue';

import GridItem from './grid-item.vue';

export default {
  components: {
    GridLayout,
    GridItem,
    WidgetWrapper,
  },
  data() {
    const defaultHeight = 20;
    const defaultWidth = 6;
    const widgets = [
      generateWidgetByType(WIDGET_TYPES.alarmList),
      generateWidgetByType(WIDGET_TYPES.context),
      generateWidgetByType(WIDGET_TYPES.statsCurves),
      generateWidgetByType(WIDGET_TYPES.alarmList),
      generateWidgetByType(WIDGET_TYPES.statsHistogram),
    ];

    return {
      layout: widgets.map((widget, index) => ({
        x: (index % 2) * defaultWidth,
        y: index * defaultHeight,
        w: defaultWidth,
        h: defaultHeight,
        i: { index, widget, fixedHeight: true },
      })),
      draggable: true,
      resizable: true,
      editMode: true,
      index: 0,
      widgets,
    };
  },
  methods: {
    resized(i, h, previousW, height, width, autoSize) {
      if (!autoSize && !i.fixedHeight) {
        this.$nextTick(() => this.$refs.widgets[i.index].$parent.autoSizeHeight());
      }
    },

    getWidgetFlexStyle(widget) {
      return {
        gridColumnStart: widget.x + 1,
        gridColumnEnd: widget.x + 1 + widget.w,
        gridRowStart: widget.y + 1,
        gridRowEnd: widget.y + widget.h,
      };
    },

    change(value, index) {
      if (!value) {
        this.$nextTick(() => this.$refs.widgets[index].$parent.autoSizeHeight());
      }
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
