<template lang="pug">
  v-container(fluid)
    v-layout
      v-btn(icon, @click="showSettings")
        v-icon settings
    fade-transition
      v-layout(wrap, v-show="!watchersPending")
        v-flex(v-for="item in watchers", :key="item._id", :class="flexSize")
          weather-item(:watcher="item", :widget="widget", :template="widget.block_template")
    fade-transition
      v-layout(column v-show="watchersPending")
        v-flex(xs12)
          v-layout(justify-center)
            v-progress-circular(indeterminate, color="primary")
</template>

<script>
import entitiesWatcherMixin from '@/mixins/entities/watcher';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';
import widgetQueryMixin from '@/mixins/widget/query';
import sideBarMixin from '@/mixins/side-bar/side-bar';
import FadeTransition from '@/components/transition/fade.vue';

import { SIDE_BARS } from '@/constants';

import WeatherItem from './weather-item.vue';

export default {
  components: {
    WeatherItem,
    FadeTransition,
  },
  mixins: [entitiesWatcherMixin, entitiesUserPreferenceMixin, widgetQueryMixin, sideBarMixin],
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  computed: {
    flexSize() {
      return [
        `xs${this.widget.columnSM}`,
        `md${this.widget.columnMD}`,
        `lg${this.widget.columnLG}`,
      ];
    },
  },
  methods: {
    fetchList() {
      this.fetchWatchersList({
        filter: this.widget.filter,
        params: this.getQuery(),
        widgetId: this.widget.id,
      });
    },

    showSettings() {
      this.showSideBar({
        name: SIDE_BARS.weatherSettings,
        config: {
          widget: this.widget,
        },
      });
    },
  },
};
</script>
