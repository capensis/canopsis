<template lang="pug">
  .weather-watcher-entity-expansion-panel
    v-expansion-panel(dark)
      v-expansion-panel-content(:style="{ backgroundColor: color }")
        v-layout(slot="header", justify-space-between, align-center)
          v-flex.pa-2(v-for="(icon, index) in mainIcons", :key="index")
            v-icon(color="white", small) {{ icon }}
          v-flex.pl-1.white--text.subheading.entity-title(
          xs12,
          )
            v-layout(align-center)
              div.mr-1.entityName(
              v-resize-text="{maxFontSize: '16px'}",
              ) {{ { entity } | get(entityNameField, false, entityNameField) }}
              v-btn.mx-1.white(v-for="icon in extraIcons", :key="icon.icon", :color="icon.color", small, dark, icon)
                v-icon(small) {{ icon.icon }}
        v-card(color="white black--text")
          v-card-text
            v-layout(v-if="availableActions.length", row, align-center)
              div {{ $t('common.actionsLabel') }}:
              div(v-for="action in availableActions", :key="action.eventType")
                v-btn(
                @click.stop="action.action",
                :disabled="!isActionBtnEnable(action.eventType)",
                depressed,
                small,
                light,
                )
                  v-icon {{ action.icon }}
            div(v-html="compiledTemplate")
</template>

<script>
import { find, isNull, pickBy } from 'lodash';

import {
  WATCHER_PBEHAVIOR_COLOR,
  WATCHER_STATES_COLORS,
  WEATHER_ICONS,
  EVENT_ENTITY_STYLE,
  EVENT_ENTITY_TYPES,
  ENTITIES_STATES,
  PBEHAVIOR_TYPES,
  WIDGETS_ACTIONS_TYPES,
} from '@/constants';
import { compile } from '@/helpers/handlebars';

import authMixin from '@/mixins/auth';
import modalMixin from '@/mixins/modal';
import widgetActionPanelWatcherEntityMixin from '@/mixins/widget/actions-panel/watcher-entity';

export default {
  mixins: [
    authMixin,
    modalMixin,
    widgetActionPanelWatcherEntityMixin,
  ],
  props: {
    entity: {
      type: Object,
      required: true,
    },
    entityNameField: {
      type: String,
      default: 'name',
    },
    template: {
      type: String,
    },
  },
  data() {
    const { weather: weatherActionsTypes } = WIDGETS_ACTIONS_TYPES;

    return {
      state: this.entity.state.val,
      actionsClicked: [],
      actionsMap: {
        ack: {
          type: weatherActionsTypes.entityAck,
          eventType: EVENT_ENTITY_TYPES.ack,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.ack].icon,
          action: this.prepareAckAction,
        },
        declareTicket: {
          type: weatherActionsTypes.entityDeclareTicker,
          eventType: EVENT_ENTITY_TYPES.declareTicket,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.declareTicket].icon,
          action: this.prepareDeclareTicketAction,
        },
        validate: {
          type: weatherActionsTypes.entityValidate,
          eventType: EVENT_ENTITY_TYPES.validate,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.validate].icon,
          action: this.prepareValidateAction,
        },
        invalidate: {
          type: weatherActionsTypes.entityInvalidate,
          eventType: EVENT_ENTITY_TYPES.invalidate,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.invalidate].icon,
          action: this.prepareInvalidateAction,
        },
        pause: {
          type: weatherActionsTypes.entityPause,
          eventType: EVENT_ENTITY_TYPES.pause,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.pause].icon,
          action: this.preparePauseAction,
        },
        play: {
          type: weatherActionsTypes.entityPlay,
          eventType: EVENT_ENTITY_TYPES.play,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.play].icon,
          action: this.preparePlayAction,
        },
      },
    };
  },
  computed: {
    color() {
      if (this.hasActivePbehavior) {
        return WATCHER_PBEHAVIOR_COLOR;
      }

      return WATCHER_STATES_COLORS[this.state];
    },

    mainIcons() {
      const mainIcons = [];
      if (!this.isPaused && !this.hasActivePbehavior) {
        mainIcons.push(WEATHER_ICONS[this.entity.state.val]);
      }

      const pausePbehavior = find(this.entity.pbehavior, { type_: PBEHAVIOR_TYPES.pause });
      const maintenancePbehavior = find(this.entity.pbehavior, { type_: PBEHAVIOR_TYPES.maintenance });
      const outOfSurveillancePbehavior = find(this.entity.pbehavior, { type_: PBEHAVIOR_TYPES.outOfSurveillance });

      if (maintenancePbehavior) {
        mainIcons.push(WEATHER_ICONS.maintenance);
      }

      if (outOfSurveillancePbehavior) {
        mainIcons.push(WEATHER_ICONS.outOfSurveillance);
      }

      if (pausePbehavior) {
        mainIcons.push(WEATHER_ICONS.pause);
      }

      return mainIcons;
    },

    extraIcons() {
      const extraIcons = [];

      if (this.entity.ack) {
        extraIcons.push({
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.fastAck].icon,
          color: 'purple',
        });
      }

      if (this.entity.ticket) {
        extraIcons.push({
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.assocTicket].icon,
          color: 'blue',
        });
      }

      return extraIcons;
    },

    hasActivePbehavior() {
      if (!this.entity.pbehavior || !this.entity.pbehavior.length) {
        return false;
      }

      return this.entity.pbehavior.filter(value => value.isActive).length;
    },

    compiledTemplate() {
      return compile(this.template, { entity: this.entity });
    },

    isPaused() {
      return this.entity.pbehavior.some(pbehavior => pbehavior.type_.toLowerCase() === PBEHAVIOR_TYPES.pause);
    },

    filteredActionsMap() {
      return pickBy(this.actionsMap, this.actionsAccessFilterHandler);
    },

    availableActions() {
      const { filteredActionsMap } = this;
      const actions = [filteredActionsMap.declareTicket];

      if (this.entity.state.val === ENTITIES_STATES.major) {
        actions.push(filteredActionsMap.validate, filteredActionsMap.invalidate);
      }

      if (this.entity.state !== ENTITIES_STATES.ok && isNull(this.entity.ack)) {
        actions.push(filteredActionsMap.ack);
      }

      if (this.isPaused) {
        actions.push(filteredActionsMap.play);
      } else {
        actions.push(filteredActionsMap.pause);
      }

      return actions.filter(action => !!action);
    },

    isActionBtnEnable() {
      return action => !this.actionsClicked.includes(action);
    },
  },
};
</script>

<style lang="scss" scoped>
  .weather-watcher-entity-expansion-panel /deep/ .v-expansion-panel__header {
    padding: 0;
    height: auto;

    .entity-title {
      line-height: 52px;
    }

    .actions-button-wrapper {
      float: right;
    }
  }

  .entityName {
    line-height: 1.5em;
  }
</style>
