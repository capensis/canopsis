<template lang="pug">
  div
    mq-layout(mq="xl")
      v-layout
        actions-panel-item(
        v-for="(action, index) in actions",
        v-bind="action",
        :key="`main-${index}`"
        )
        v-menu(v-show="dropDownActions && dropDownActions.length", bottom, left, @click.native.stop)
          v-btn(icon, slot="activator")
            v-icon more_vert
          v-list
            actions-panel-item(
            v-for="(action, index) in dropDownActions",
            v-bind="action",
            isDropDown,
            :key="`drop-down-${index}`"
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
            isDropDown,
            :key="`mobile-main-${index}`"
            )
            actions-panel-item(
            v-for="(action, index) in dropDownActions",
            v-bind="action",
            isDropDown,
            :key="`mobile-drop-down-${index}`"
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
