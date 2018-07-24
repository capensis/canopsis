<template lang="pug">
  v-navigation-drawer(
  :temporary="$mq === 'mobile' || $mq === 'tablet'",
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
        v-icon.closeIcon(@click.stop="close", color="white") close
      v-divider
      div(
      :is="config.component",
      :widget="widget",
      :isNew="isNew"
      )
</template>

<script>
import AlarmSettingsFields from './alarm-settings-fields.vue';
import ContextSettingsFields from './context-settings-fields.vue';

/**
 * Settings wrapper component
 *
 * @prop {bool} value - controls visibility of current component
 * @prop {string} title - title for settings header
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
        listalarm: {
          title: this.$t('settings.titles.alarmListSettings'),
          component: 'alarm-settings-fields',
        },
        crudcontext: {
          title: this.$t('settings.titles.contextTableSettings'),
          component: 'context-settings-fields',
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
    close() {
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
