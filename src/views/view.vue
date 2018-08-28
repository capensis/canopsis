<template lang="pug">
  v-container
    div
      div(v-for="widget in widgets", :key="widget._id")
        h2 {{ widget.title }}
        div(
        :is="widgetsMap[widget.xtype]",
        :widget="widget",
        @openSettings="openSettings(widget)"
        )
    v-speed-dial.fab(
    direction="top",
    :open-on-hover="true",
    transition="scale-transition"
    )
      v-btn(slot="activator", v-model="fab", color="green darken-3", dark, fab)
        v-icon add
      v-tooltip(left)
        v-btn(slot="activator", fab, dark, small, color="indigo", @click.prevent="showCreateWidgetModal")
          v-icon widgets
        span {{ $t('common.widget') }}
</template>

<script>
import { MODALS, WIDGET_TYPES } from '@/constants';

import AlarmsListContainer from '@/containers/alarms-list.vue';
import EntitiesListContainer from '@/containers/entities-list.vue';
import WeatherContainer from '@/containers/weather.vue';

import modalMixin from '@/mixins/modal/modal';
import entitiesViewMixin from '@/mixins/entities/view';
import entitiesWidgetMixin from '@/mixins/entities/widget';

export default {
  components: {
    AlarmsListContainer,
    EntitiesListContainer,
    WeatherContainer,
  },
  mixins: [
    modalMixin,
    entitiesViewMixin,
    entitiesWidgetMixin,
  ],
  props: {
    id: {
      type: [String, Number],
      required: true,
    },
  },
  data() {
    return {
      fab: false,
      widgetsMap: {
        [WIDGET_TYPES.alarmList]: 'alarms-list-container',
        [WIDGET_TYPES.context]: 'entities-list-container',
        [WIDGET_TYPES.weather]: 'weather-container',
      },
    };
  },
  mounted() {
    this.fetchView({ id: this.id });
  },
  methods: {
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
