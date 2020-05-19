<template lang="pug">
  dashboard(id="tab", ref="dashboard")
    dash-layout(
      v-bind="layout",
      :debug="true",
      :key="layout.breakpoint"
    )
      dash-item(
        v-for="item in layout.items",
        v-bind.sync="item",
        :key="item.id"
      )
        div.wrapper
          widget-wrapper(
            :widget="item.widget",
            :tab="tab",
            :updateTabMethod="updateTabMethod"
          )
</template>

<script>
import { omit } from 'lodash';
import { Dashboard, DashLayout, DashItem } from 'vue-responsive-dash';

import WidgetWrapper from '@/components/widgets/widget-wrapper.vue';
import { setSeveralFields } from '@/helpers/immutable';

export default {
  components: {
    Dashboard, DashLayout, DashItem, WidgetWrapper,
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
    const items = this.tab.widgets.map(widget => ({
      ...widget.gridParameters.desktop,
      width: widget.gridParameters.desktop.w,
      height: widget.gridParameters.desktop.h,

      id: widget._id,
      widget,
    }));

    return {
      layout: {
        breakpoint: 'xl',
        numberOfCols: 12,
        rowHeight: 20,
        useCssTransforms: true,
        colWidth: (this.$parent.$el.offsetWidth - 130) / 12,
        items,
      },
    };
  },
  created() {
    // this.$nextTick(() => this.$refs.dashboard.onResize({ detail: { width: this.$parent.$el.offsetWidth } }));
  },
  mounted() {
    this.layout.colWidth = null;
  },
  beforeDestroy() {
    const fields = this.tab.widgets.reduce((acc, { gridParameters }, index) => {
      acc[`widgets.${index}.gridParameters.desktop`] = {
        ...gridParameters.desktop,
        ...omit(this.layout.items[index], ['i', 'widget', 'width', 'height']),

        w: this.layout.items[index].width,
        h: this.layout.items[index].height,
      };

      return acc;
    }, {});

    const newTab = setSeveralFields(this.tab, fields);

    this.updateTabMethod(newTab);
  },
  methods: {
    resized() {},
  },
};
</script>

<style lang="scss">
  .grid-layout-wrapper {
    padding-bottom: 500px;
  }

  #tab {
    position: relative;
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
