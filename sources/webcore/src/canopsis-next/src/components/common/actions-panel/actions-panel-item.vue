<template lang="pug">
  span(data-test="actionsPanelItem")
    template(v-if="isDropDown")
      v-list-tile(:disabled="disabled", @click.stop="method")
        v-list-tile-title
          v-icon.pr-3(
            :color="iconColor",
            :disabled="disabled",
            left,
            small
          ) {{ icon }}
          span.body-1 {{ title }}
    template(v-else)
      c-action-btn(
        :tooltip="title",
        :disabled="disabled",
        :icon="icon",
        :color="iconColor",
        @click="method"
      )
</template>


<script>
/**
 * Component showing an action icon
 *
 * @module alarm
 *
 * @prop {string} title - Action title
 * @prop {string} icon - Action icon
 * @prop {Function} method - Action to execute when user clicks on the action's icon
 * @prop {string} [iconColor=''] - Action icon color
 * @prop {boolean} [isDropDown=false] - Boolean to decide whether to show a dropdown with actions, or actions separately
 */
export default {
  props: {
    title: {
      type: String,
      required: true,
    },
    icon: {
      type: String,
      required: true,
    },
    method: {
      type: Function,
      required: true,
    },
    iconColor: {
      type: String,
      default: '',
    },
    isDropDown: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
