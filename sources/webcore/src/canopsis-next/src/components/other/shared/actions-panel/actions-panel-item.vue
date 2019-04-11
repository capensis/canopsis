<template lang="pug" functional>
  span
    template(v-if="props.isDropDown")
      v-list-tile(@click.stop="props.method")
        v-list-tile-title
          v-icon.pr-3(:class="props.iconClass", left, small) {{ props.icon }}
          span.body-1 {{ props.title }}
    template(v-else)
      v-tooltip(bottom)
        v-btn.mx-1(slot="activator", flat, icon, @click.stop="props.method")
          v-icon(:class="props.iconClass") {{ props.icon }}
        span {{ props.title }}
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
 * @prop {string} [iconClass=''] - Action icon className
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
