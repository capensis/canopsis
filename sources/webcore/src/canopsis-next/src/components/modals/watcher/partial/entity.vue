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
            v-layout(justify-space-between)
              v-flex(xs11)
                v-layout(justify-space-between)
                  v-layout(v-if="availableActions.length", row, align-center)
                    div {{ $t('common.actionsLabel') }}:
                    div(v-for="action in availableActions", :key="action.eventType")
                      v-tooltip(top)
                        v-btn(
                        slot="activator",
                        @click.stop="action.action",
                        :disabled="!isActionBtnEnable(action.eventType)",
                        depressed,
                        small,
                        light,
                        )
                          v-icon {{ action.icon }}
                        span {{ $t(`common.actions.${action.eventType}`) }}
              v-tooltip(v-if="hasActivePbehavior && hasAccessToManagePbehaviors", top)
                v-btn(small, @click="showPbehaviorsListModal", slot="activator")
                  v-icon(small) edit
                span {{ $t('modals.watcher.editPbehaviors') }}
            entity-template(:entity="entity", :template="template")
</template>

<script>
import { find, isNull, pickBy } from 'lodash';

import {
  CRUD_ACTIONS,
  MODALS,
  WATCHER_STATES_COLORS,
  WEATHER_ICONS,
  EVENT_ENTITY_STYLE,
  EVENT_ENTITY_TYPES,
  ENTITIES_STATES,
  ENTITIES_STATES_STYLES,
  PBEHAVIOR_TYPES,
  WIDGETS_ACTIONS_TYPES,
  USERS_RIGHTS,
  ENTITIES_STATUSES,
} from '@/constants';

import authMixin from '@/mixins/auth';
import modalMixin from '@/mixins/modal';
import widgetActionPanelWatcherEntityMixin from '@/mixins/widget/actions-panel/watcher-entity';
import entitiesPbehaviorCommentMixin from '@/mixins/entities/pbehavior/comment';
import entitiesWatcherEntityMixin from '@/mixins/entities/watcher-entity';

import EntityTemplate from './entity-template.vue';

export default {
  components: {
    EntityTemplate,
  },
  mixins: [
    authMixin,
    modalMixin,
    widgetActionPanelWatcherEntityMixin,
    entitiesPbehaviorCommentMixin,
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
  data() {
    const { weather: weatherActionsTypes } = WIDGETS_ACTIONS_TYPES;

    return {
      menu: false,
      state: this.entity.state.val,
      actionsClicked: [],
      actionsMap: {
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
      },
    };
  },
  computed: {
    color() {
      if (this.hasActivePbehavior || this.isWatcherOnPbehavior) {
        return WATCHER_STATES_COLORS.pause;
      }

      return ENTITIES_STATES_STYLES[this.state].color;
    },

    mainIcons() {
      const state = ENTITIES_STATES_STYLES[this.entity.state.val].text;
      const mainIcons = [];
      if (!this.isPaused && !this.hasActivePbehavior) {
        mainIcons.push(WEATHER_ICONS[state]);
      }

      const pausePbehavior = find(this.entity.pbehavior, { type_: PBEHAVIOR_TYPES.pause });
      const maintenancePbehavior = find(this.entity.pbehavior, { type_: PBEHAVIOR_TYPES.maintenance });
      const outOfSurveillancePbehavior = find(this.entity.pbehavior, { type_: PBEHAVIOR_TYPES.unmonitored });

      if (maintenancePbehavior) {
        mainIcons.push(WEATHER_ICONS.maintenance);
      }

      if (outOfSurveillancePbehavior) {
        mainIcons.push(WEATHER_ICONS.unmonitored);
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

      if (this.entity.status && this.entity.status.val === ENTITIES_STATUSES.cancelled) {
        extraIcons.push({
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.delete].icon,
          color: 'grey darken-1',
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

    isPaused() {
      return this.entity.pbehavior.some(pbehavior => pbehavior.type_.toLowerCase() === PBEHAVIOR_TYPES.pause);
    },

    filteredActionsMap() {
      return pickBy(this.actionsMap, this.actionsAccessFilterHandler);
    },

    availableActions() {
      const { filteredActionsMap } = this;
      const actions = [];

      if (this.entity.state !== ENTITIES_STATES.ok && isNull(this.entity.ack)) {
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

    isActionBtnEnable() {
      return action => !this.actionsClicked.includes(action);
    },

    pausePbehaviors() {
      return this.entity.pbehavior.filter(pbehavior => pbehavior.type_.toLowerCase() === PBEHAVIOR_TYPES.pause);
    },

    hasAccessToManagePbehaviors() {
      return this.checkAccess(USERS_RIGHTS.business.weather.actions.entityManagePbehaviors);
    },
  },
  methods: {
    showPbehaviorsListModal() {
      this.showModal({
        name: MODALS.pbehaviorList,
        config: {
          pbehaviors: this.pausePbehaviors,
          entityId: this.entity.entity_id,
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
