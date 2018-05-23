<template lang="pug">
  div
    v-alert(v-model="isVisible", :type="type", dismissible)
      div {{ text }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

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
  },
  data() {
    return {
      visible: true,
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
          this.remove({ id: this.id });
        }
      },
    },
  },
  methods: {
    ...mapActions(['remove']),
  },
};
</script>
