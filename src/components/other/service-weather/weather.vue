<template lang="pug">
  div
    v-layout
      v-btn(icon, @click="showSettings")
        v-icon settings
    v-fade-transition
      v-layout(v-show="!watchersPending", wrap)
        v-flex(v-for="item in watchers", :key="item._id", :class="flexSize")
          weather-item(:watcher="item", :widget="widget", :template="widget.parameters.blockTemplate")
    v-fade-transition
      v-layout(v-show="watchersPending", column)
        v-flex(xs12)
          v-layout(justify-center)
            v-progress-circular(indeterminate, color="primary")
</template>

<script>
import entitiesWatcherMixin from '@/mixins/entities/watcher';
import widgetQueryMixin from '@/mixins/widget/query';
import sideBarMixin from '@/mixins/side-bar/side-bar';

import WeatherItem from './weather-item.vue';

export default {
  components: {
    WeatherItem,
  },
  mixins: [
    entitiesWatcherMixin,
    widgetQueryMixin,
    sideBarMixin,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
    rowId: {
      type: String,
      required: true,
    },
  },
  computed: {
    flexSize() {
      return [
        `xs${this.widget.parameters.columnSM}`,
        `md${this.widget.parameters.columnMD}`,
        `lg${this.widget.parameters.columnLG}`,
      ];
    },
  },
  methods: {
    showSettings() {
      this.showSideBar({
        name: this.$constants.SIDE_BARS.weatherSettings,
        config: {
          widget: this.widget,
          rowId: this.rowId,
        },
      });
    },

    fetchList() {
      this.fetchWatchersList({
        filter: this.widget.parameters.mfilter.filter,
        params: this.getQuery(),
        widgetId: this.widget._id,
      });
    },
  },
};
</script>
