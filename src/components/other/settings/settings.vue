<template lang="pug">
  v-navigation-drawer(
  :temporary="$mq | mq({ m: true, l: false })",
  :value="value",
  @input="$emit('input', $event)",
  width="400",
  temporary,
  fixed,
  right
  )
    div(v-if="widget")
      v-toolbar(color="blue darken-4")
        v-list
          v-list-tile
            v-list-tile-title.white--text.text-xs-center {{ config.title }}
        v-icon.closeIcon(@click.stop="closeSettings", color="white") close
      v-divider
      div(
      :is="config.component",
      :widget="widget",
      :isNew="isNew",
      @closeSettings="closeSettings"
      )
</template>

<script>
import { WIDGET_TYPES } from '@/constants';

import AlarmSettingsFields from './alarm-settings-fields.vue';
import ContextSettingsFields from './context-settings-fields.vue';

/**
 * Settings component
 *
 * @prop {bool} value - controls visibility of current component
 * @prop {Object} widget - active widget
 * @prop {bool} isNew - is widget new
 */
export default {
  components: {
    AlarmSettingsFields,
    ContextSettingsFields,
  },
  props: {
    value: {
      type: Boolean,
      default: false,
    },
    widget: {
      type: Object,
      default: () => ({}),
    },
    isNew: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      widgetsSettingsMap: {
        [WIDGET_TYPES.alarmList]: {
          title: this.$t('settings.titles.alarmListSettings'),
          component: 'alarm-settings-fields',
        },
        [WIDGET_TYPES.context]: {
          title: this.$t('settings.titles.contextTableSettings'),
          component: 'context-settings-fields',
        },
        [WIDGET_TYPES.weather]: {
          title: this.$t('settings.titles.weatherSettings'),
          component: 'weather-settings-fields',
        },
      },
    };
  },
  computed: {
    config() {
      if (this.widget && this.widget.xtype) {
        return this.widgetsSettingsMap[this.widget.xtype];
      }

      return {};
    },
  },
  methods: {
    closeSettings() {
      this.$emit('input', false);
    },
  },
};
</script>

<style scoped>
  .closeIcon:hover {
    cursor: pointer;
  }
</style>
