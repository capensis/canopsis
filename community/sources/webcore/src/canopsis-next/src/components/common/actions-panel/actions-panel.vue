<template lang="pug">
  div.actions-panel(:class="{ 'actions-panel--dense': dense }")
    mq-layout(mq="xl")
      v-layout(row, align-center)
        actions-panel-item(
          v-for="(action, index) in actions",
          v-bind="action",
          :key="`main-${index}`"
        )
        v-menu(v-if="dropDownActions.length", bottom, left, @click.native.stop="")
          template(#activator="{ on }")
            v-btn.mr-0(v-on="on", icon)
              v-icon more_vert
          v-list
            actions-panel-item(
              v-for="(action, index) in dropDownActions",
              v-bind="action",
              is-drop-down,
              :key="`drop-down-${index}`"
            )
    mq-layout(:mq="['m', 't', 'l']")
      v-layout
        v-menu(
          v-if="actions.length || dropDownActions.length",
          bottom,
          left,
          @click.native.stop=""
        )
          template(#activator="{ on }")
            v-btn(v-on="on", icon)
              v-icon more_vert
          v-list
            actions-panel-item(
              v-for="(action, index) in actions",
              v-bind="action",
              is-drop-down,
              :key="`mobile-main-${index}`"
            )
            actions-panel-item(
              v-for="(action, index) in dropDownActions",
              v-bind="action",
              is-drop-down,
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
    dense: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style lang="scss">
.actions-panel {
  &--dense {
    .v-btn--icon {
      width: 24px;
      height: 24px;

      .v-icon {
        font-size: 20px;
      }
    }
  }
}
</style>
