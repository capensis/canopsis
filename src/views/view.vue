<template lang="pug">
  v-container
    div
      div(v-for="widgetWrapper in widgetWrappers", :key="widgetWrapper._id")
        div(
        :is="widgetsMap[widgetWrapper.widget.xtype]",
        :widget="widgetWrapper.widget",
        @openSettings="openSettings(widgetWrapper.widget)"
        )
    v-speed-dial.fab(
    direction="top",
    :open-on-hover="true",
    transition="scale-transition"
    )
      v-btn(slot="activator", v-model="fab", color="green darken-3", dark, fab)
        v-icon add
      v-tooltip(left)
        v-btn(slot="activator", fab, dark, small, color="indigo", @click.prevent="showInsertWidgetModal")
          v-icon widgets
        span widget
    settings(v-model="isSettingsOpen", :widget="activeWidgetSettings", :isWidgetNew="isActiveWidgetNew")
</template>

<script>
import Settings from '@/components/other/settings/settings.vue';
import AlarmListContainer from '@/containers/alarm-list.vue';
import EntitiesListContainer from '@/containers/entities-list.vue';
import viewMixin from '@/mixins/view';
import modalMixin from '@/mixins/modal/modal';
import { MODALS } from '@/constants';

export default {
  components: {
    AlarmListContainer,
    EntitiesListContainer,
    Settings,
  },
  mixins: [
    viewMixin,
    modalMixin,
  ],
  props: {
    id: {
      type: [String, Number],
      required: true,
    },
  },
  data() {
    return {
      activeWidgetSettings: null,
      isActiveWidgetNew: false,
      fab: false,
      widgetsMap: {
        listalarm: 'alarm-list-container',
        crudcontext: 'entities-list-container',
      },
    };
  },
  computed: {
    isSettingsOpen: {
      get() {
        return !!this.activeWidgetSettings;
      },
      set(value) {
        if (!value) {
          this.activeWidgetSettings = null;
        }
      },
    },
  },
  mounted() {
    this.fetchView({ id: this.id });
  },
  methods: {
    openSettings(widget, isNew) {
      this.activeWidgetSettings = widget;
      this.isActiveWidgetNew = isNew;
    },
    closeSettings() {
      this.activeWidgetSettings = null;
      this.isActiveWidgetNew = false;
    },
    showInsertWidgetModal() {
      this.showModal({
        name: MODALS.insertWidget,
        config: {
          action: widget => this.activeWidgetSettings = widget,
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
