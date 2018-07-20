<template lang="pug">
  v-container
    div
      div(v-for="widgetWrapper in widgetWrappers", :key="widgetWrapper._id")
        div(:is="widgetsMap[widgetWrapper.widget.xtype]", :widget="widgetWrapper.widget", @openSettings="openSettings")
    v-speed-dial.fab(
    direction="top",
    :open-on-hover="true",
    transition="scale-transition"
    )
      v-btn(slot="activator", v-model="fab", color="green darken-3", dark, fab)
        v-icon add
      v-tooltip(left)
        v-btn(slot="activator", fab, dark, small, color="indigo", @click.prevent="")
          v-icon widgets
        span widget
    v-navigation-drawer(v-model="settingsIsOpen", fixed, temporary, right, width="400")
      v-toolbar(color="blue darken-4")
        v-list
          v-list-tile
            v-list-tile-title.white--text.text-xs-center Title
        v-icon.closeIcon(@click.stop="closeSettings", color="white") close
      v-divider
      alarm-settings-fields
</template>

<script>
import AlarmSettingsFields from '@/components/other/settings/alarm-settings-fields.vue';
import AlarmListContainer from '@/containers/alarm-list.vue';
import EntitiesListContainer from '@/containers/entities-list.vue';
import viewMixin from '@/mixins/view';
import modalMixin from '@/mixins/modal/modal';
import { MODALS } from '@/constants';

export default {
  components: {
    AlarmListContainer,
    EntitiesListContainer,
    AlarmSettingsFields,
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
      settingsIsOpen: false,
      fab: false,
      widgetsMap: {
        listalarm: 'alarm-list-container',
        crudcontext: 'entities-list-container',
      },
    };
  },
  mounted() {
    this.fetchView({ id: this.id });
  },
  methods: {
    openSettings(widgetId) {
      this.settingsIsOpen = true;
      this.activeWidgetSettings = widgetId;
    },
    closeSettings() {
      this.settingsIsOpen = false;
      this.activeWidgetSettings = null;
    },
    showInsertWidgetModal() {
      this.showModal({ name: MODALS.insertWidget });
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
