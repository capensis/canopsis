<template lang="pug">
  v-container
    div
      div(v-for="widget in widgets", :key="`${widgetKeyPrefix}_${widget.id}`")
        h2 {{ widget.title }}
        component(
        :is="widgetsMap[widget.xtype]",
        :widget="widget",
        )
    .fab
      v-btn(@click="refreshView", icon, color="info", dark, fab)
        v-icon refresh
      v-speed-dial(
      :open-on-hover="true",
      transition="scale-transition",
      direction="left",
      )
        v-btn(slot="activator", color="green darken-3", dark, fab)
          v-icon add
        v-tooltip(left)
          v-btn(slot="activator", fab, dark, small, color="indigo", @click.prevent="showCreateWidgetModal")
            v-icon widgets
          span {{ $t('common.widget') }}
</template>

<script>
import { MODALS, WIDGET_TYPES } from '@/constants';
import uid from '@/helpers/uid';

import AlarmsList from '@/components/other/alarm/alarms-list.vue';
import EntitiesList from '@/components/other/context/entities-list.vue';
import Weather from '@/components/other/service-weather/weather.vue';

import modalMixin from '@/mixins/modal/modal';
import entitiesViewMixin from '@/mixins/entities/view';
import entitiesWidgetMixin from '@/mixins/entities/widget';

export default {
  components: {
    AlarmsList,
    EntitiesList,
    Weather,
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
      widgetKeyPrefix: uid(),
      widgetsMap: {
        [WIDGET_TYPES.alarmList]: 'alarms-list',
        [WIDGET_TYPES.context]: 'entities-list',
        [WIDGET_TYPES.weather]: 'weather',
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

    async refreshView() {
      await this.fetchView({ id: this.id });

      this.widgetKeyPrefix = uid();
    },
  },
};
</script>
