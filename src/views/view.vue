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
    settings(v-model="isSettingsOpen", :widget="activeWidgetSettings")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import Settings from '@/components/other/settings/settings.vue';
import AlarmListContainer from '@/containers/alarm-list.vue';
import EntitiesListContainer from '@/containers/entities-list.vue';
import viewMixin from '@/mixins/view';
import modalMixin from '@/mixins/modal/modal';
import { MODALS } from '@/constants';

const { mapActions } = createNamespacedHelpers('userPreference');

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
  async mounted() {
    await this.fetchView({ id: this.id });

    // TODO: fix it
    await Promise.all(this.widgetWrappers.map(widgetWrapper => this.fetchUserPreferences({
      params: {
        limit: 1,
        filter: {
          crecord_name: 'root',
          widget_id: widgetWrapper.widget.id,
          _id: `${widgetWrapper.widget.id}_root`, // TODO: change to real user
        },
      },
    })));
  },
  methods: {
    ...mapActions({ fetchUserPreferences: 'fetchList' }),

    openSettings(widgetId) {
      this.activeWidgetSettings = widgetId;
    },
    closeSettings() {
      this.activeWidgetSettings = null;
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
