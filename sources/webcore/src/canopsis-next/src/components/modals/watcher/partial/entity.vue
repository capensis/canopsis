<template lang="pug">
  .weather-watcher-entity-expansion-panel
    v-expansion-panel
      v-expansion-panel-content(hide-actions)
        v-layout.pa-2(slot="header", :class="entityClass", justify-space-between)
          span.pl-1.white--text.subheading.entity-title {{ entity | get(entityNameField) }}
          v-layout(justify-end)
            div(v-for="action in availableActions", :key="action.name")
              v-btn.secondary(
              @click.stop="action.action(action.name)",
              :disabled="!isActionBtnEnable(action.name)",
              small,
              fab,
              depressed
              )
                v-icon {{ action.icon }}
        v-card
          v-card-text
            div(v-html="compiledTemplate")
        v-divider
</template>

<script>
import {
  MODALS,
  WATCHER_PBEHAVIOR_COLOR,
  WATCHER_STATES_COLORS,
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
    console.log(this.entity);
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
    entityClass() {
      if (this.hasActivePbehavior) {
        return WATCHER_PBEHAVIOR_COLOR;
      }

      return WATCHER_STATES_COLORS[this.state];
    },

    hasActivePbehavior() {
      if (!this.entity.pbehavior || !this.entity.pbehavior.length) {
        return false;
      }

      return this.entity.pbehavior.filter((value) => {
        const start = value.dtstart * 1000;
        const end = value.dtend * 1000;
        const now = Date.now();

        return start <= now && now < end;
      }).length;
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
</style>
