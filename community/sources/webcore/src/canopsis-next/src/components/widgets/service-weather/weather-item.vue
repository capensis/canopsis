<template lang="pug">
  v-card.white--text.cursor-pointer.weather__item(
    :class="itemClasses",
    :style="{ height: itemHeight + 'em', backgroundColor: color}",
    tile,
    @click.native="showAdditionalInfoModal"
  )
    v-layout.fill-height(row)
      v-flex.position-relative.fill-height
        v-icon.weather__item--background.white--text(size="5em") {{ icon }}
        v-layout(:class="{ blinking: isBlinking }", justify-start)
          v-runtime-template.watcherName.pa-3(:template="compiledTemplate")
        v-btn.helpBtn.ma-0(
          v-if="hasVariablesHelpAccess",
          icon,
          small,
          @click.stop="showVariablesHelpModal(watcher)"
        )
          v-icon help
        v-btn.pauseIcon(v-if="secondaryIcon", icon)
          v-icon(color="white") {{ secondaryIcon }}
      v-flex(v-if="isCountersEnabled", xs2)
        alarm-counters.fill-height(
          :counters="counters",
          :selected-types="selectedTypes"
        )
    v-btn.see-alarms-btn(
      v-if="isBothModalType && hasAlarmsListAccess",
      flat,
      @click.stop="showAlarmListModal"
    ) {{ $t('serviceWeather.seeAlarms') }}
</template>

<script>
import VRuntimeTemplate from 'v-runtime-template';

import {
  MODALS,
  USERS_RIGHTS,
  WIDGET_TYPES,
  WATCHER_STATES_COLORS,
  WEATHER_ICONS,
  SERVICE_WEATHER_WIDGET_MODAL_TYPES,
} from '@/constants';

import { compile } from '@/helpers/handlebars';
import { generateWidgetByType } from '@/helpers/entities';

import authMixin from '@/mixins/auth';
import entitiesWatcherEntityMixin from '@/mixins/entities/watcher-entity';

import { convertObjectToTreeview } from '@/helpers/treeview';

import AlarmCounters from './alarm-counters.vue';

export default {
  components: {
    AlarmCounters,
    VRuntimeTemplate,
  },
  mixins: [authMixin, entitiesWatcherEntityMixin],
  props: {
    watcher: {
      type: Object,
      required: true,
    },
    template: {
      type: String,
    },
    widget: {
      type: Object,
    },
  },
  asyncComputed: {
    compiledTemplate: {
      async get() {
        const compiledTemplate = await compile(this.template, { entity: this.watcher });

        return `<div>${compiledTemplate}</div>`;
      },
      default: '',
    },
  },
  computed: {
    hasMoreInfosAccess() {
      return this.checkAccess(USERS_RIGHTS.business.weather.actions.moreInfos);
    },

    hasAlarmsListAccess() {
      return this.checkAccess(USERS_RIGHTS.business.weather.actions.alarmsList);
    },

    hasVariablesHelpAccess() {
      return this.checkAccess(USERS_RIGHTS.business.weather.actions.variablesHelp);
    },

    color() {
      return WATCHER_STATES_COLORS[this.watcher.color];
    },

    icon() {
      return WEATHER_ICONS[this.watcher.icon];
    },

    secondaryIcon() {
      return WEATHER_ICONS[this.watcher.secondary_icon];
    },

    itemClasses() {
      const classes = [
        `mt-${this.widget.parameters.margin.top}`,
        `mr-${this.widget.parameters.margin.right}`,
        `mb-${this.widget.parameters.margin.bottom}`,
        `ml-${this.widget.parameters.margin.left}`,
      ];

      if (this.isBothModalType && this.hasAlarmsListAccess) {
        classes.push('v-card__with-see-alarms-btn');
      }

      return classes;
    },

    itemHeight() {
      return 4 + this.widget.parameters.heightFactor;
    },

    isBlinking() {
      return this.watcher.is_action_required;
    },

    isBothModalType() {
      return this.widget.parameters.modalType === SERVICE_WEATHER_WIDGET_MODAL_TYPES.both;
    },

    isAlarmListModalType() {
      return this.widget.parameters.modalType === SERVICE_WEATHER_WIDGET_MODAL_TYPES.alarmList;
    },

    counters() {
      return this.watcher.alarm_counters || [];
    },

    hasCounters() {
      return this.counters.length;
    },

    selectedTypes() {
      const { counters } = this.widget.parameters;

      return counters ? counters.types : [];
    },

    hasSelectedTypes() {
      return this.selectedTypes.length;
    },

    isCountersEnabled() {
      const { counters = {} } = this.widget.parameters;

      return counters.enabled && this.hasCounters && this.hasSelectedTypes;
    },
  },
  methods: {
    showAdditionalInfoModal(e) {
      if (e.target.tagName !== 'A' || !e.target.href) {
        if (this.isAlarmListModalType && this.hasAlarmsListAccess) {
          this.showAlarmListModal();
        } else if (!this.isAlarmListModalType && this.hasMoreInfosAccess) {
          this.showMainInfoModal();
        }
      }
    },

    showMainInfoModal() {
      this.$modals.show({
        name: MODALS.watcher,
        config: {
          color: this.color,
          watcher: this.watcher,
          entityTemplate: this.widget.parameters.entityTemplate,
          modalTemplate: this.widget.parameters.modalTemplate,
          itemsPerPage: this.widget.parameters.modalItemsPerPage,
        },
      });
    },

    async showAlarmListModal() {
      try {
        const widget = generateWidgetByType(WIDGET_TYPES.alarmList);

        const filter = { $and: [{ 'entity.impact': this.watcher._id }] };

        const watcherFilter = {
          title: this.watcher.name,
          filter,
        };

        const widgetParameters = {
          ...this.widget.parameters.alarmsList,

          mainFilter: watcherFilter,
          viewFilters: [watcherFilter],
        };

        this.$modals.show({
          name: MODALS.alarmsList,
          config: {
            widget: {
              ...widget,

              parameters: {
                ...widget.parameters,
                ...widgetParameters,
              },
            },
          },
        });
      } catch (err) {
        this.$popups.error({
          text: this.$t('errors.default'),
        });
      }
    },

    showVariablesHelpModal() {
      const entityFields = convertObjectToTreeview(this.watcher, 'entity');
      const variables = [entityFields];

      this.$modals.show({
        name: MODALS.variablesHelp,
        config: {
          variables,
        },
      });
    },
  },
};
</script>

<style lang="scss" scoped>
  @keyframes blink {
    0% { opacity: 1 }
    50% { opacity: 0.3 }
  }

  .blinking {
    animation: blink 2s linear infinite;
  }
</style>
