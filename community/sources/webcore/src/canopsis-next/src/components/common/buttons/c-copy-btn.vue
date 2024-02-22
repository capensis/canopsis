<template>
  <c-action-btn
    v-bind="$attrs"
    :tooltip="tooltip"
  >
    <template #button="">
      <v-btn
        v-clipboard:copy="value"
        v-clipboard:success="onSuccessCopied"
        v-clipboard:error="onErrorCopied"
        :small="small"
        :fab="fab"
        class="mx-1 ma-0"
        icon
      >
        <v-icon
          :color="color"
          :small="iconSmall"
        >
          {{ icon }}
        </v-icon>
      </v-btn>
    </template>
  </c-action-btn>
</template>

<script>
export default {
  props: {
    icon: {
      type: String,
      default: 'content_copy',
    },
    color: {
      type: String,
      default: '',
    },
    tooltip: {
      type: String,
      default: '',
    },
    value: {
      type: String,
      required: true,
    },
    small: {
      type: Boolean,
      default: false,
    },
    iconSmall: {
      type: Boolean,
      default: false,
    },
    fab: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    onSuccessCopied() {
      this.$emit('success');
    },

    onErrorCopied() {
      this.$emit('error');
    },
  },
};
</script>
