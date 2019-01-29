<template lang="pug">
  span
    template(v-if="isDropDown")
      v-list-tile(@click.stop="method")
        v-list-tile-title
          v-icon.pr-3(:class="iconClass", left, small) {{ icon }}
          span.body-1 {{ title }}
    template(v-else)
      v-tooltip(bottom)
        v-btn.mx-1(slot="activator", flat, icon, @click.stop="method")
          v-icon(:class="iconClass") {{ icon }}
        span {{ title }}
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
    iconClass: {
      type: String,
      default: '',
    },
    isDropDown: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
