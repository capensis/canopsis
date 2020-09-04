<template lang="pug">
  div(data-test="massActionsPanel")
    mq-layout(mq="l+")
      actions-panel-item(
        v-for="(action, index) in actions",
        v-bind="action",
        :key="`multiple-${index}`",
        @click="action.method"
      )
    mq-layout(mq="t")
      v-menu(bottom, left, @click.native.stop)
        v-btn(icon, slot="activator")
          v-icon more_vert
        v-list
          actions-panel-item(
            v-for="(action, index) in actions",
            v-bind="action",
            :key="`mobile-multiple-${index}`",
            isDropDown,
            @click="action.method"
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
