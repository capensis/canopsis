<template lang="pug">
  v-container
    div
      div(v-for="row in rows", :key="row.title")
        h1 {{ row.title }}
        div(v-for="widget in row.widgets")
          h2 {{ widget.type }}
          component(
          :is="widgetsMap[widget.type]",
          :widget="widget",
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

import { WIDGET_TYPES } from '@/constants';
import uid from '@/helpers/uid';

import AlarmsList from '@/components/other/alarm/alarms-list-new.vue';
import EntitiesList from '@/components/other/context/entities-list.vue';
import Weather from '@/components/other/service-weather/weather.vue';

import entitiesViewV3Mixin from '@/mixins/entities/view-v3/view-v3';

export default {
  components: {
    AlarmsList,
    EntitiesList,
    Weather,
  },
  mixins: [entitiesViewV3Mixin],
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
  mounted() {
    this.fetchView({ id: this.id });
  },
  methods: {
    async refreshView() {
      await this.fetchView({ id: this.id });

      this.widgetKeyPrefix = uid();
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
