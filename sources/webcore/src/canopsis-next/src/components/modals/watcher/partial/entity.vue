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
            v-layout(row, align-center)
              div {{ $t('common.actionsLabel') }}:
              div(v-for="action in availableActions", :key="action.name")
                v-btn(
                @click.stop="action.action(action.name)",
                :disabled="!isActionBtnEnable(action.name)",
                depressed,
                small,
                light,
                )
                  v-icon {{ action.icon }}
            div(v-html="compiledTemplate")
</template>

<script>
import find from 'lodash/find';

import {
  MODALS,
  WATCHER_PBEHAVIOR_COLOR,
  WATCHER_STATES_COLORS,
  WEATHER_ICONS,
  WEATHER_ACK_EVENT_OUTPUT,
  EVENT_ENTITY_STYLE,
  EVENT_ENTITY_TYPES,
  ENTITIES_STATES,
  PBEHAVIOR_TYPES,
} from '@/constants';
import { compile } from '@/helpers/handlebars';

import modalMixin from '@/mixins/modal';
import weatherEventsMixin from '@/mixins/weather-event-actions';

export default {
  mixins: [modalMixin, weatherEventsMixin],
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
    return {
      state: this.entity.state.val,
      actionsClicked: [],
      actions: [
        {
          name: EVENT_ENTITY_TYPES.ack,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.ack].icon,
          action: this.prepareAckAction,
        },
        {
          name: EVENT_ENTITY_TYPES.declareTicket,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.declareTicket].icon,
          action: this.prepareDeclareTicketAction,
        },
        {
          name: EVENT_ENTITY_TYPES.validate,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.validate].icon,
          action: this.prepareValidateAction,
        },
        {
          name: EVENT_ENTITY_TYPES.invalidate,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.invalidate].icon,
          action: this.prepareInvalidateAction,
        },
        {
          name: EVENT_ENTITY_TYPES.pause,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.pause].icon,
          action: this.preparePauseAction,
        },
        {
          name: EVENT_ENTITY_TYPES.play,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.play].icon,
          action: this.preparePlayAction,
        },
      ],
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

    availableActions() {
      return this.actions.filter((action) => {
        if (
          this.entity.state.val !== ENTITIES_STATES.major &&
          (action.name === EVENT_ENTITY_TYPES.invalidate || action.name === EVENT_ENTITY_TYPES.validate)
        ) {
          return false;
        }

        if (
          (this.entity.state === ENTITIES_STATES.ok || this.entity.ack !== null) &&
          action.name === EVENT_ENTITY_TYPES.ack
        ) {
          return false;
        }

        if (this.isPaused && action.name === EVENT_ENTITY_TYPES.pause) {
          return false;
        }

        if (!this.isPaused && action.name === EVENT_ENTITY_TYPES.play) {
          return false;
        }

        return true;
      });
    },

    isActionBtnEnable() {
      return action => !this.actionsClicked.includes(action);
    },
  },

  methods: {
    prepareAckAction() {
      this.addAckActionToQueue({ entity: this.entity, output: WEATHER_ACK_EVENT_OUTPUT.ack });
      this.actionsClicked.push(EVENT_ENTITY_TYPES.ack);
    },

    prepareDeclareTicketAction() {
      this.showModal({
        name: MODALS.createWatcherDeclareTicketEvent,
        config: {
          action: (ticket) => {
            this.addAckActionToQueue({ entity: this.entity, output: WEATHER_ACK_EVENT_OUTPUT.ack });
            this.addDeclareTicketActionToQueue({ entity: this.entity, ticket });
            this.actionsClicked.push(EVENT_ENTITY_TYPES.declareTicket);
          },
        },
      });
    },

    prepareValidateAction() {
      this.addAckActionToQueue({ entity: this.entity, output: WEATHER_ACK_EVENT_OUTPUT.validateOk });
      this.addValidateActionToQueue({ entity: this.entity });
      this.actionsClicked.push(EVENT_ENTITY_TYPES.validate);
    },

    prepareInvalidateAction() {
      this.addAckActionToQueue({ entity: this.entity, output: WEATHER_ACK_EVENT_OUTPUT.ack });
      this.addInvalidateActionToQueue({ entity: this.entity });
      this.actionsClicked.push(EVENT_ENTITY_TYPES.invalidate);
    },

    preparePauseAction() {
      this.showModal({
        name: MODALS.createWatcherPauseEvent,
        config: {
          action: (pause) => {
            this.addPauseActionToQueue({
              entity: this.entity,
              comment: pause.comment,
              reason: pause.reason,
            });
            this.actionsClicked.push(EVENT_ENTITY_TYPES.pause);
          },
        },
      });
    },

    preparePlayAction() {
      this.addPlayActionToQueue({ entity: this.entity });
      this.actionsClicked.push(EVENT_ENTITY_TYPES.play);
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
