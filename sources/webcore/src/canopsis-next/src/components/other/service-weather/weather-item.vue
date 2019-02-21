<template lang="pug">
  v-card.white--text.cursor-pointer(
  :class="getItemClasses",
  :style="{ height: itemHeight + 'em', backgroundColor: format.color}",
  tile,
  @click.native="showAdditionalInfoModal"
  )
    v-btn.helpBtn.ma-0(@click.stop="showVariablesHelpModal(watcher)", v-if="isEditingMode", icon, small)
      v-icon help
    div(:class="{ blinking: isBlinking }")
      v-layout(justify-start)
        v-icon.px-3.py-2.white--text(size="2em") {{ format.icon }}
        div.watcherName.pt-3(v-html="compiledTemplate")
        v-btn.pauseIcon(v-if="watcher.active_pb_some && !watcher.active_pb_all", icon)
          v-icon(color="white") {{ secondaryIcon }}
</template>

<script>
import { find } from 'lodash';

import {
  MODALS,
  USERS_RIGHTS,
  WIDGET_TYPES,
  WATCHER_STATES_COLORS,
  WATCHER_PBEHAVIOR_COLOR,
  PBEHAVIOR_TYPES,
  WEATHER_ICONS,
  SERVICE_WEATHER_WIDGET_MODAL_TYPES,
} from '@/constants';

import { compile } from '@/helpers/handlebars';
import { generateWidgetByType } from '@/helpers/entities';
import { prepareFilterWithFieldsPrefix } from '@/helpers/filter';

import authMixin from '@/mixins/auth';
import modalMixin from '@/mixins/modal';
import popupMixin from '@/mixins/popup';
import entitiesWatcherEntityMixin from '@/mixins/entities/watcher-entity';

import convertObjectFieldToTreeBranch from '@/helpers/treeview';

export default {
  mixins: [authMixin, modalMixin, popupMixin, entitiesWatcherEntityMixin],
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
    isEditingMode: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    hasMoreInfosAccess() {
      return this.checkAccess(USERS_RIGHTS.business.weather.actions.moreInfos);
    },
    hasAlarmsListAccess() {
      return this.checkAccess(USERS_RIGHTS.business.weather.actions.alarmsList);
    },
    isPaused() {
      return this.watcher.active_pb_all;
    },
    hasWatcherPbehavior() {
      return this.watcher.active_pb_watcher;
    },
    isPbehavior() {
      return this.watcher.pbehavior.some(pbehavior => pbehavior.isActive);
    },
    format() {
      if (!this.isPaused && !this.hasWatcherPbehavior) {
        const state = this.watcher.state.val;

        return {
          icon: WEATHER_ICONS[state],
          color: WATCHER_STATES_COLORS[state],
        };
      }

      const pbehaviors = this.hasWatcherPbehavior ? this.watcher.watcher_pbehavior : this.watcher.pbehavior;

      const maintenancePbehavior = find(pbehaviors, { type_: PBEHAVIOR_TYPES.maintenance });
      const outOfSurveillancePbehavior = find(pbehaviors, { type_: PBEHAVIOR_TYPES.outOfSurveillance });

      let icon = WEATHER_ICONS.pause;

      if (maintenancePbehavior) {
        icon = WEATHER_ICONS.maintenance;
      } else if (outOfSurveillancePbehavior) {
        icon = WEATHER_ICONS.outOfSurveillance;
      }

      return {
        color: WATCHER_PBEHAVIOR_COLOR,
        icon,
      };
    },
    secondaryIcon() {
      if (this.watcher.pbehavior.some(value => value.type_ === PBEHAVIOR_TYPES.maintenance)) {
        return WEATHER_ICONS.maintenance;
      } else if (this.watcher.pbehavior.every(value => value.type_ === PBEHAVIOR_TYPES.outOfSurveillance)) {
        return WEATHER_ICONS.outOfSurveillance;
      }

      return WEATHER_ICONS.pause;
    },
    compiledTemplate() {
      return compile(this.template, { entity: this.watcher });
    },
    getItemClasses() {
      return [
        `mt-${this.widget.parameters.margin.top}`,
        `mr-${this.widget.parameters.margin.right}`,
        `mb-${this.widget.parameters.margin.bottom}`,
        `ml-${this.widget.parameters.margin.left}`,
      ];
    },
    itemHeight() {
      return 4 + this.widget.parameters.heightFactor;
    },
    isBlinking() {
      return (
        this.watcher.alerts_not_ack
        && !this.hasWatcherPbehavior
        && !this.isPbehavior
      );
    },
  },
  methods: {
    showAdditionalInfoModal() {
      const isAlarmListModalType = this.widget.parameters.modalType === SERVICE_WEATHER_WIDGET_MODAL_TYPES.alarmList;

      if (isAlarmListModalType && this.hasAlarmsListAccess) {
        this.showAlarmListModal();
      } else if (!isAlarmListModalType && this.hasMoreInfosAccess) {
        this.showMainInfoModal();
      }
    },

    showMainInfoModal() {
      this.showModal({
        name: MODALS.watcher,
        config: {
          watcherId: this.watcher.entity_id,
          entityTemplate: this.widget.parameters.entityTemplate,
          modalTemplate: this.widget.parameters.modalTemplate,
        },
      });
    },

    async showAlarmListModal() {
      try {
        const initialFilter = JSON.parse(this.watcher.mfilter);
        const newFilter = prepareFilterWithFieldsPrefix(initialFilter, 'entity.');
        const widget = generateWidgetByType(WIDGET_TYPES.alarmList);
        const watcherFilter = {
          title: this.watcher.display_name,
          filter: newFilter,
        };

        const widgetParameters = {
          ...this.widget.parameters.alarmsList,

          mainFilter: watcherFilter,
          viewFilters: [watcherFilter],
        };

        this.showModal({
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
        this.addErrorPopup({
          text: this.$t('errors.default'),
        });
      }
    },

    showVariablesHelpModal() {
      const entityFields = convertObjectFieldToTreeBranch(this.watcher, 'entity');
      const variables = [entityFields];

      this.showModal({
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
  .pauseIcon {
    position: absolute;
    right: 0;
    bottom: 0;
    cursor: inherit;
  }

  .watcherName {
    color: white;
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    line-height: 1.2em;
  }

  @keyframes blink {
    0% { opacity: 1 }
    50% { opacity: 0.3 }
  }

  .blinking {
    animation: blink 2s linear infinite;
  }

  .cursor-pointer {
    cursor: pointer;
  }

  .helpBtn {
    position: absolute;
    right: 0.2em;
    top: 0;
    z-index: 1;
  }
</style>
