<template lang="pug">
  div(data-test="sharedActionsPanel")
    mq-layout(mq="xl")
      v-layout
        actions-panel-item(
          v-for="(action, index) in actions",
          v-bind="action",
          :key="`main-${index}`",
          @click="action.method",
          @action="$emit('action', $event)"
        )
        v-menu(v-show="dropDownActions && dropDownActions.length", bottom, left, @click.native.stop)
          v-btn(data-test="dropDownActionsButton", icon, slot="activator")
            v-icon more_vert
          v-list(data-test="dropDownActions")
            actions-panel-item(
              v-for="(action, index) in dropDownActions",
              v-bind="action",
              :key="`drop-down-${index}`",
              isDropDown,
              @click="action.method",
              @action="$emit('action', $event)"
            )
    mq-layout(mq="l")
      v-layout
        v-menu(
          v-show="actions.length + Object.keys(actions).length > 0",
          bottom,
          left,
          @click.native.stop
        )
          v-btn(icon, slot="activator")
            v-icon more_vert
          v-list
            actions-panel-item(
              v-for="(action, index) in actions",
              v-bind="action",
              :key="`mobile-main-${index}`",
              isDropDown,
              @click="action.method",
              @action="$emit('action', $event)"
            )
            actions-panel-item(
              v-for="(action, index) in dropDownActions",
              v-bind="action",
              :key="`mobile-drop-down-${index}`",
              isDropDown,
              @click="action.method",
              @action="$emit('action', $event)"
            )
</template>

<script>
import ActionsPanelItem from './actions-panel-item.vue';

/**
 * Component to regroup actions (actions-panel-item) for each alarm on the alarms list
 *
 * @module alarm
 *
 * @prop {Array} [actions=[]] - Actions object
 * @prop {Array} [dropDownActions=[]] - Drop down actions object
 */
export default {
  components: { ActionsPanelItem },
  props: {
    actions: {
      type: Array,
      default: () => [],
    },
    dropDownActions: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
