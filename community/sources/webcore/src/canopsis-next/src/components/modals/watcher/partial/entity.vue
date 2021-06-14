<template lang="pug">
  .weather-watcher-entity-expansion-panel
    v-expansion-panel(dark)
      v-expansion-panel-content(:style="{ backgroundColor: color }")
        v-layout(slot="header", justify-space-between, align-center)
          v-flex.pa-2(v-for="(icon, index) in mainIcons", :key="index")
            v-icon(color="white", small) {{ icon }}
          v-flex.pl-1.white--text.subheading.entity-title(xs12)
            v-layout(align-center)
              div.mr-1.entityName(
                v-resize-text="{ maxFontSize: '16px' }",
              ) {{ { entity } | get(entityNameField, false, entityNameField) }}
              v-btn.mx-1.white(v-for="icon in extraIcons", :key="icon.icon", :color="icon.color", small, dark, icon)
                v-icon(small) {{ icon.icon }}
        v-card(color="white black--text")
          v-card-text
            v-layout(justify-space-between)
              v-flex(xs11)
                v-layout(justify-space-between)
                  v-layout(v-if="availableActions.length", row, align-center)
                    div {{ $t('common.actionsLabel') }}:
                    div(v-for="action in availableActions", :key="action.eventType")
                      v-tooltip(top)
                        v-btn(
                          slot="activator",
                          :disabled="isActionDisabled(action.eventType)",
                          depressed,
                          small,
                          light,
                          @click.stop="action.action"
                        )
                          v-icon {{ action.icon }}
                        span {{ $t(`common.actions.${action.eventType}`) }}
              v-tooltip(v-if="hasActivePbehavior && hasAccessToManagePbehaviors", top)
                v-btn(slot="activator", @click="showPbehaviorsListModal", small)
                  v-icon(small) list
                span {{ $t('modals.watcher.editPbehaviors') }}
            entity-template(:entity="entity", :template="template")
</template>

<script>
import { isNull, pickBy } from 'lodash';

import {
  CRUD_ACTIONS,
  MODALS,
  WATCHER_STATES_COLORS,
  WEATHER_ICONS,
  EVENT_ENTITY_STYLE,
  EVENT_ENTITY_TYPES,
  ENTITIES_STATES,
  ENTITIES_STATES_STYLES,
  WIDGETS_ACTIONS_TYPES,
  USERS_RIGHTS,
  ENTITIES_STATUSES,
  PBEHAVIOR_TYPE_TYPES,
} from '@/constants';

import authMixin from '@/mixins/auth';
import widgetActionPanelWatcherEntityMixin from '@/mixins/widget/actions-panel/watcher-entity';
import entitiesWatcherEntityMixin from '@/mixins/entities/watcher-entity';

import EntityTemplate from './entity-template.vue';

export default {
  inject: ['$eventsQueue'],
  components: {
    EntityTemplate,
  },
  mixins: [
    authMixin,
    widgetActionPanelWatcherEntityMixin,
    entitiesWatcherEntityMixin,
  ],
  props: {
    watcherId: {
      type: String,
      required: true,
    },
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
    isWatcherOnPbehavior: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    actionsClicked() {
      return this.$eventsQueue.queue
        .reduce((acc, { entityId, type }) => {
          if (entityId === this.entity._id) {
            acc.push(type);
          }

          return acc;
        }, []);
    },

    actionsMap() {
      const { weather: weatherActionsTypes } = WIDGETS_ACTIONS_TYPES;

      return {
        ack: {
          type: weatherActionsTypes.entityAck,
          eventType: EVENT_ENTITY_TYPES.ack,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.ack].icon,
          action: this.prepareAckAction,
        },
        assocTicket: {
          type: weatherActionsTypes.entityAssocTicket,
          eventType: EVENT_ENTITY_TYPES.assocTicket,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.assocTicket].icon,
          action: this.prepareAssocTicketAction,
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
        cancel: {
          type: weatherActionsTypes.entityCancel,
          eventType: EVENT_ENTITY_TYPES.cancel,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.delete].icon,
          action: this.prepareCancelAction,
        },
        comment: {
          type: weatherActionsTypes.entityComment,
          eventType: EVENT_ENTITY_TYPES.comment,
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.comment].icon,
          action: this.prepareCommentAction,
        },
      };
    },

    color() {
      if (this.hasActivePbehavior || this.isWatcherOnPbehavior) {
        return WATCHER_STATES_COLORS.pause;
      }

      return ENTITIES_STATES_STYLES[this.entity.state.val].color;
    },

    mainIcons() {
      const mainIcons = [];

      if (!this.isPaused && !this.hasActivePbehavior) {
        mainIcons.push(WEATHER_ICONS[this.entity.icon]);
      }

      mainIcons.push(...this.entity.pbehaviors.map(({ type }) => type.icon_name));

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

      if (this.entity.status && this.entity.status.val === ENTITIES_STATUSES.cancelled) {
        extraIcons.push({
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.delete].icon,
          color: 'grey darken-1',
        });
      }

      return extraIcons;
    },

    hasActivePbehavior() {
      return this.entity.pbehaviors.some(pbehavior => pbehavior.type.type === PBEHAVIOR_TYPE_TYPES.active);
    },

    isPaused() {
      return this.entity.pbehaviors.some(pbehavior => pbehavior.type.type === PBEHAVIOR_TYPE_TYPES.pause);
    },

    filteredActionsMap() {
      return pickBy(this.actionsMap, this.actionsAccessFilterHandler);
    },

    availableActions() {
      const { filteredActionsMap } = this;
      const actions = [filteredActionsMap.comment];

      if (this.entity.state.val !== ENTITIES_STATES.ok && isNull(this.entity.ack)) {
        actions.push(filteredActionsMap.ack);
      }

      actions.push(filteredActionsMap.assocTicket);

      if (this.entity.state.val === ENTITIES_STATES.major) {
        actions.push(filteredActionsMap.validate, filteredActionsMap.invalidate);
      }

      if (this.isPaused) {
        actions.push(filteredActionsMap.play);
      } else {
        actions.push(filteredActionsMap.pause);
      }

      if (
        this.entity.alarm_display_name &&
        (!this.entity.status || this.entity.status.val !== ENTITIES_STATUSES.cancelled)
      ) {
        actions.push(filteredActionsMap.cancel);
      }

      return actions.filter(action => !!action);
    },

    pausePbehaviors() {
      return this.entity.pbehaviors.filter(pbehavior => pbehavior.type.type === PBEHAVIOR_TYPE_TYPES.pause);
    },

    hasAccessToManagePbehaviors() {
      return this.checkAccess(USERS_RIGHTS.business.weather.actions.entityManagePbehaviors);
    },
  },
  methods: {
    isActionDisabled(action) {
      return this.actionsClicked.includes(action);
    },

    showPbehaviorsListModal() {
      this.$modals.show({
        name: MODALS.pbehaviorList,
        config: {
          pbehaviors: this.pausePbehaviors,
          entityId: this.entity._id,
          onlyActive: true,
          availableActions: [CRUD_ACTIONS.delete, CRUD_ACTIONS.update],
        },
      });
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

  .pbehavior {
    cursor: pointer;
  }
</style>
