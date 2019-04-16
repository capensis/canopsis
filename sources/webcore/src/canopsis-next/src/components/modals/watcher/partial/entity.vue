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
              v-menu(bottom, offset-y, @click.native.stop, open-on-hover)
                v-btn(small, slot="activator")
                  v-icon(small) comment
                v-list
                  v-list-tile(v-for="pbehavior in pausePbehaviors", :key="pbehavior._id")
                    v-list-tile-title(@click="() => showEditCommentsModal(pbehavior)") {{ pbehavior._id }}
            entity-template(:entity="entity", :template="template")
</template>

<script>
import { find, isNull, pickBy, omit } from 'lodash';

import {
  MODALS,
  WATCHER_PBEHAVIOR_COLOR,
  WATCHER_STATES_COLORS,
  WEATHER_ICONS,
  EVENT_ENTITY_STYLE,
  EVENT_ENTITY_TYPES,
  ENTITIES_STATES,
  PBEHAVIOR_TYPES,
  WIDGETS_ACTIONS_TYPES,
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

    isPaused() {
      return this.entity.pbehavior.some(pbehavior => pbehavior.type_.toLowerCase() === PBEHAVIOR_TYPES.pause);
    },

    filteredActionsMap() {
      return pickBy(this.actionsMap, this.actionsAccessFilterHandler);
    },

    availableActions() {
      const { filteredActionsMap } = this;
      const actions = [filteredActionsMap.assocTicket];

      if (this.entity.state.val === ENTITIES_STATES.major) {
        actions.push(filteredActionsMap.validate, filteredActionsMap.invalidate);
      }

      if (this.entity.state !== ENTITIES_STATES.ok && isNull(this.entity.ack)) {
        actions.push(filteredActionsMap.ack);
      }

      if (this.entity.alarm_display_name) {
        actions.push(filteredActionsMap.cancel);
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

    pausePbehaviors() {
      return this.entity.pbehavior.filter(pbehavior => pbehavior.type_.toLowerCase() === PBEHAVIOR_TYPES.pause);
    },
  },
  methods: {
    showEditCommentsModal(pbehavior) {
      this.showModal({
        name: MODALS.editPbehaviorComments,
        config: {
          title: 'Edit comments',
          comments: pbehavior.comments,
          action: async (comments) => {
            const newComments = await comments.map(comment => omit(comment, ['key']));
            await this.updateSeveralPbehaviorComments({ pbehavior, comments: newComments });
            await this.fetchWatcherEntitiesList({ watcherId: this.watcherId });
          },
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
</style>
