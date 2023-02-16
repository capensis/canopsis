<template lang="pug">
  div
    mq-layout(mq="xl")
      v-layout(row, align-center)
        actions-panel-item(
          v-for="action in inlineActions",
          v-bind="action",
          :key="action.type"
        )
        v-menu(v-if="dropdownActions.length", bottom, left, @click.native.stop="")
          template(#activator="{ on }")
            v-btn.mr-0(v-on="on", icon)
              v-icon more_vert
          v-list
            actions-panel-item(
              v-for="action in dropdownActions",
              v-bind="action",
              is-drop-down,
              :key="action.type"
            )
    mq-layout(:mq="['m', 't', 'l']")
      v-layout
        v-menu(
          v-if="actions.length",
          bottom,
          left,
          @click.native.stop=""
        )
          template(#activator="{ on }")
            v-btn(v-on="on", icon)
              v-icon more_vert
          v-list
            actions-panel-item(
              v-for="action in actions",
              v-bind="action",
              is-drop-down,
              :key="action.type"
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
 */
export default {
  components: { ActionsPanelItem },
  props: {
    actions: {
      type: Array,
      default: () => [],
    },
    inlineCount: {
      type: Number,
      default: 3,
    },
  },
  computed: {
    inlineActions() {
      return this.actions.slice(0, this.inlineCount);
    },

    dropdownActions() {
      return this.actions.slice(this.inlineCount);
    },
  },
};
</script>
