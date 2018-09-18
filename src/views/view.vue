<template lang="pug">
  v-container
    div
      div(v-for="row in rows", :key="row._id")
        h1 {{ row.title }}
        div(v-for="widget in row.widgets")
          h2 {{ widget.type }}
          component(
          :is="widgetsMap[widget.type]",
          :widget="widget",
          :key="`${widgetKeyPrefix}_${widget._id}`"
          )
    .fab
      v-btn(@click="refreshView", icon, color="info", dark, fab)
        v-icon refresh
      v-speed-dial(
      direction="left",
      :open-on-hover="true",
      transition="scale-transition"
      )
        v-btn(slot="activator", color="green darken-3", dark, fab)
          v-icon add
        v-tooltip(left)
          v-btn(slot="activator", fab, dark, small, color="indigo", @click.prevent="showCreateWidgetModal")
            v-icon widgets
          span {{ $t('common.widget') }}
</template>

<script>
import get from 'lodash/get';

import { WIDGET_TYPES, MODALS } from '@/constants';
import uid from '@/helpers/uid';

import AlarmsList from '@/components/other/alarm/alarms-list.vue';
import EntitiesList from '@/components/other/context/entities-list.vue';
import Weather from '@/components/other/service-weather/weather.vue';

import modalMixin from '@/mixins/modal/modal';
import entitiesViewMixin from '@/mixins/entities/view';

export default {
  components: {
    AlarmsList,
    EntitiesList,
    Weather,
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
      widgetsMap: {
        [WIDGET_TYPES.alarmList]: 'alarms-list',
        [WIDGET_TYPES.context]: 'entities-list',
        [WIDGET_TYPES.weather]: 'weather',
      },
      widgetKeyPrefix: uid(),
    };
  },
  computed: {
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
