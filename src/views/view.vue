<template lang="pug">
  v-container
    div
      div(v-for="widget in widgets", :key="widget._id")
        h2 {{ widget.title }} {{ widgetsMap[widget.xtype] }}
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
    settings(v-model="isSettingsOpen", :widget="widgetSettings", :isNew="isWidgetNew")
</template>

<script>
import Settings from '@/components/other/settings/settings.vue';
import AlarmsListContainer from '@/containers/alarms-list.vue';
import EntitiesListContainer from '@/containers/entities-list.vue';
import WeatherContainer from '@/containers/weather.vue';

import { MODALS, WIDGET_TYPES } from '@/constants';
import modalMixin from '@/mixins/modal/modal';
import entitiesViewMixin from '@/mixins/entities/view';
import entitiesWidgetMixin from '@/mixins/entities/widget';

export default {
  components: {
    AlarmsListContainer,
    EntitiesListContainer,
    WeatherContainer,
    Settings,
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
      widgetSettings: null,
      isWidgetNew: false,
      fab: false,
      widgetsMap: {
        [WIDGET_TYPES.alarmList]: 'alarms-list-container',
        [WIDGET_TYPES.context]: 'entities-list-container',
        [WIDGET_TYPES.weather]: 'weather-container',
      },
    };
  },
  computed: {
    isSettingsOpen: {
      get() {
        return !!this.widgetSettings;
      },
      set(value) {
        if (!value) {
          this.closeSettings();
        }
      },
    },
  },
  mounted() {
    this.fetchView({ id: this.id });
  },
  methods: {
    openSettings(widget, isNew) {
      this.widgetSettings = widget;
      this.isWidgetNew = isNew;
    },
    closeSettings() {
      this.widgetSettings = null;
      this.isWidgetNew = false;
    },
    showCreateWidgetModal() {
      this.showModal({
        name: MODALS.createWidget,
        config: {
          action: widget => this.openSettings(widget, true),
        },
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
