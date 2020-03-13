<template lang="pug">
  v-card.white--text.cursor-pointer(
    :class="itemClasses",
    :style="{ height: itemHeight + 'em', backgroundColor: color}",
    tile,
    @click.native="showAdditionalInfoModal"
  )
    v-btn.helpBtn.ma-0(
      v-if="isEditingMode && hasVariablesHelpAccess",
      icon,
      small,
      @click.stop="showVariablesHelpModal(watcher)"
    )
      v-icon help
    div(:class="{ blinking: isBlinking }")
      v-layout(justify-start)
        v-icon.px-3.py-2.white--text(size="2em") {{ icon }}
        v-runtime-template.watcherName.pt-3(:template="compiledTemplate")
        v-btn.pauseIcon(v-if="secondaryIcon", icon)
          v-icon(color="white") {{ secondaryIcon }}
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

export default {
  components: {
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
    isEditingMode: {
      type: Boolean,
      default: false,
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
      return WATCHER_STATES_COLORS[this.watcher.tileColor];
    },

    icon() {
      return WEATHER_ICONS[this.watcher.tileIcon];
    },

    secondaryIcon() {
      return WEATHER_ICONS[this.watcher.tileSecondaryIcon];
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
      return this.watcher.isActionRequired;
    },

    isBothModalType() {
      return this.widget.parameters.modalType === SERVICE_WEATHER_WIDGET_MODAL_TYPES.both;
    },

    isAlarmListModalType() {
      return this.widget.parameters.modalType === SERVICE_WEATHER_WIDGET_MODAL_TYPES.alarmList;
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

        const filter = { $and: [{ 'entity.impact': this.watcher.entity_id }] };

        const watcherFilter = {
          title: this.watcher.display_name,
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
  $seeAlarmBtnHeight: 18px;

  .v-card__with-see-alarms-btn {
    padding-bottom: $seeAlarmBtnHeight;

    .see-alarms-btn {
      position: absolute;
      bottom: 0;
      width: 100%;
      font-size: .6em;
      height: $seeAlarmBtnHeight;
      color: white;
      margin: 0;
      background-color: rgba(0, 0, 0, .2);

      &.v-btn--active:before, &.v-btn:focus:before, &.v-btn:hover:before {
        background-color: rgba(0, 0, 0, .5);
      }
    }
  }

  .pauseIcon {
    position: absolute;
    right: 0;
    bottom: 1em;
    cursor: inherit;
  }

  .watcherName {
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    line-height: 1.2em;

    &, & /deep/ a {
      color: white;
    }
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
