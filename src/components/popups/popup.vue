<template lang="pug">
  v-alert(v-model="isVisible", :type="type", transition="fade-transition", dismissible)
    div(v-html="text")
</template>

<script>
// LIBS
import { createNamespacedHelpers } from 'vuex';
// OTHERS
import { POPUP_AUTO_CLOSE_DELAY, VUETIFY_ANIMATION_DELAY } from '@/config';

const { mapActions } = createNamespacedHelpers('popup');

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
      this.timeouts.push(setTimeout(() => this.remove({ id: this.id }), VUETIFY_ANIMATION_DELAY));
    },
  },
};
</script>
