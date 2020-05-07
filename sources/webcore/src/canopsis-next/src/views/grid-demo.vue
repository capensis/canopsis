<template lang="pug">
  div.grid-layout-wrapper
    grid-layout(
      :layout.sync="layout",
      :col-num="12",
      :row-height="20",
      :is-draggable="draggable",
      :is-resizable="resizable",
      :vertical-compact="true",
      :use-css-transforms="true",
      :responsive="false"
    )
      grid-item(
        v-for="item in layout",
        :key="item.i._id",
        :x="item.x",
        :y="item.y",
        :w="item.w",
        :h="item.h",
        :i="item.i",
        dragAllowFrom=".drag-handler"
      )
        div.wrapper
          div.drag-handler
          widget-wrapper(
            :widget="item.i.widget",
            :tab="{ _id: '123' }",
            :isEditingMode="true",
            :row="{}",
            :updateTabMethod="() => {}",
            ref="widgets"
          )
</template>

<script>
import { GridLayout, GridItem } from 'vue-grid-layout';

import { WIDGET_TYPES } from '@/constants';

import { generateWidgetByType } from '@/helpers/entities';

import WidgetWrapper from '@/components/widgets/widget-wrapper.vue';

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
        i: { index, widget },
      })),
      draggable: true,
      resizable: true,
      index: 0,
      widgets,
    };
  },
};
</script>

<style lang="scss">
  .grid-layout-wrapper {
    padding-bottom: 500px;
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
    height: 100%;
    overflow-y: auto;

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
