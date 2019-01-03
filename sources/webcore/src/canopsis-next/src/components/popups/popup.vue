<template lang="pug">
  v-alert(v-model="isVisible", :type="type", transition="fade-transition", dismissible)
    div(v-html="text")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { POPUP_AUTO_CLOSE_DELAY } from '@/config';

const { mapActions } = createNamespacedHelpers('popup');

/**
* Popup component
*
* @prop {String} [id] - Id of the popup
* @prop {String} [type] - Type of the popup (info, error, ...)
* @prop {String} [text] - Text displayed in the popup
* @prop {Number,Boolean} [autoClose] - Auto close delay
*/
export default {
  props: {
    id: {
      type: String,
      required: true,
    },
    type: {
      type: String,
      default: 'error',
    },
    text: {
      type: String,
      default: '',
    },
    autoClose: {
      type: [Number, Boolean],
      default: POPUP_AUTO_CLOSE_DELAY,
    },
  },
  data() {
    return {
      timeouts: [],
      visible: false,
    };
  },
  computed: {
    isVisible: {
      get() {
        return this.visible;
      },
      set(value) {
        this.visible = value;

        if (!value) {
          this.removeWithTimeout();
        }
      },
    },
  },
  mounted() {
    this.visible = true;

    if (this.autoClose) {
      this.timeouts.push(setTimeout(() => this.isVisible = false, this.autoClose));
    }
  },
  beforeDestroy() {
    if (this.timeouts.length) {
      this.timeouts.forEach(timeout => clearTimeout(timeout));
    }
  },
  methods: {
    ...mapActions(['remove']),
    removeWithTimeout() {
      this.timeouts.push(setTimeout(() => this.remove({ id: this.id }), this.$config.VUETIFY_ANIMATION_DELAY));
    },
  },
};
</script>
