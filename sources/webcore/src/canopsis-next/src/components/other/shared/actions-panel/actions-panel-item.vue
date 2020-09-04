<template lang="pug">
  span(data-test="actionsPanelItem")
    v-list-tile(v-if="isDropDown", @click.stop="onAction")
      v-list-tile-title
        v-icon.pr-3(:class="iconClass", left, small) {{ icon }}
        span.body-1 {{ title }}
    v-tooltip(v-else, bottom)
      v-btn.mx-1(slot="activator", flat, icon, @click.stop="onAction")
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
 * @prop {string} [iconClass=''] - Action icon className
 * @prop {boolean} [isDropDown=false] - Boolean to decide whether to show a dropdown with actions, or actions separately
 */
export default {
  props: {
    type: {
      type: String,
      required: true,
    },
    title: {
      type: String,
      required: true,
    },
    icon: {
      type: String,
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
  methods: {
    onAction() {
      this.$emit('click');
      this.$emit('action', this.type);
    },
  },
};
</script>
