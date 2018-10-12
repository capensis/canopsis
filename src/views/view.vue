<template lang="pug">
  v-container
    div
      v-layout(v-for="row in rows", :key="row._id", row, wrap)
        v-flex(xs12)
          h2 {{ row.title }}
        v-flex(
        v-for="widget in row.widgets",
        :key="`${widgetKeyPrefix}_${widget._id}`",
        :class="getWidgetFlexClass(widget)"
        )
          h3 {{ widget.title }}
          component(
          :is="widgetsComponentsMap[widget.type]",
          :widget="widget",
          :rowId="row._id"
          )
    .fab
      v-speed-dial(
      direction="top",
      :open-on-hover="true",
      transition="scale-transition"
      )
        v-btn(slot="activator", color="green darken-3", dark, fab)
          v-icon menu
        v-tooltip(left)
          v-btn(slot="activator", @click="refreshView", color="info", dark, fab, small)
            v-icon refresh
          span {{ $t('common.refresh') }}
        v-tooltip(left)
          v-btn(slot="activator", fab, dark, small, color="indigo", @click.prevent="showCreateWidgetModal")
            v-icon widgets
          span {{ $t('common.addWidget') }}
        v-tooltip(left)
          v-btn(slot="activator", color="", dark, fab, small)
            v-icon edit
          span {{ $t('common.toggleEditView') }}
</template>

<script>
import get from 'lodash/get';

import { WIDGET_TYPES, MODALS } from '@/constants';
import uid from '@/helpers/uid';

import AlarmsList from '@/components/other/alarm/alarms-list.vue';
import EntitiesList from '@/components/other/context/entities-list.vue';
import Weather from '@/components/other/service-weather/weather.vue';
import StatsHistogram from '@/components/other/stats/histogram/stats-histogram-wrapper.vue';
import StatsTable from '@/components/other/stats/stats-table.vue';

import modalMixin from '@/mixins/modal/modal';
import entitiesViewMixin from '@/mixins/entities/view';

export default {
  components: {
    AlarmsList,
    EntitiesList,
    Weather,
    StatsHistogram,
    StatsTable,
  },
  mixins: [
    modalMixin,
    entitiesViewMixin,
  ],
  props: {
    id: {
      type: [String, Number],
      required: true,
    },
  },
  data() {
    return {
      widgetsComponentsMap: {
        [WIDGET_TYPES.alarmList]: 'alarms-list',
        [WIDGET_TYPES.context]: 'entities-list',
        [WIDGET_TYPES.weather]: 'weather',
        [WIDGET_TYPES.statsHistogram]: 'stats-histogram',
        [WIDGET_TYPES.statsTable]: 'stats-table',
      },
      widgetKeyPrefix: uid(),
    };
  },
  computed: {
    getWidgetFlexClass() {
      return widget => [
        `xs${widget.size.sm}`,
        `md${widget.size.md}`,
        `lg${widget.size.lg}`,
      ];
    },
    rows() {
      return get(this.view, 'rows', []);
    },
  },
  created() {
    this.fetchView({ id: this.id });
  },
  methods: {
    async refreshView() {
      await this.fetchView({ id: this.id });

      this.widgetKeyPrefix = uid();
    },

    showCreateWidgetModal() {
      this.showModal({
        name: MODALS.createWidget,
      });
    },
  },
};
</script>

<style scoped>
  .fab {
    position: fixed;
    bottom: 0;
    right: 0;
  }
</style>
