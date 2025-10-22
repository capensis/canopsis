<template>
  <c-action-btn
    v-bind="$attrs"
    :tooltip="tooltip"
  >
    <template #button="{ on: tooltipOn }">
      <div class="c-action-btn__button-wrapper" v-on="tooltipOn">
        <v-btn
          v-clipboard:copy="value"
          v-clipboard:success="onSuccessCopied"
          v-clipboard:error="onErrorCopied"
          :small="small"
          :fab="fab"
          class="mx-1 ma-0 c-action-btn__button"
          icon
        >
          <v-icon
            :color="color"
            :small="iconSmall"
          >
            {{ icon }}
          </v-icon>
        </v-btn>
      </div>
    </template>
  </c-action-btn>
</template>

<script>
export default {
  inheritAttrs: false,
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
  setup(props, { emit }) {
    const onSuccessCopied = () => emit('success');
    const onErrorCopied = () => emit('error');

    return {
      onSuccessCopied,
      onErrorCopied,
    };
  },
};
</script>
