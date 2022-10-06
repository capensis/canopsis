<template lang="pug">
  div
    mq-layout(mq="l+")
      actions-panel-item(
        v-for="(action, index) in actions",
        v-bind="action",
        :key="`multiple-${index}`"
      )
    mq-layout(:mq="['m', 't']")
      v-menu(bottom, left, @click.native.stop="")
        template(#activator="{ on, attrs }")
          v-btn(v-bind="attrs", v-on="on", icon)
            v-icon more_vert
        v-list
          actions-panel-item(
            v-for="(action, index) in actions",
            v-bind="action",
            is-drop-down,
            :key="`mobile-multiple-${index}`"
          )
</template>

<script>
import ActionsPanelItem from './actions-panel-item.vue';

/**
 * Panel regrouping mass actions icons
 *
 * @module alarm
 *
 * @prop {Array} [itemIds] - Items selected for the mass action
 */
export default {
  components: { ActionsPanelItem },
  props: {
    actions: {
      type: Array,
      required: true,
    },
  },
};
</script>
