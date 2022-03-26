<template lang="pug">
  div.broadcast-message.pa-2(
    :style="{ backgroundColor: color }",
    :title="message"
  )
    v-icon.broadcast-message__icon.mr-1(small) $vuetify.icons.bullhorn
    span.broadcast-message__text(v-html="message")
</template>

<script>
export default {
  props: {
    message: {
      type: String,
      default: '',
    },
    color: {
      type: String,
      default: '',
    },
  },
};
</script>

<style lang="scss" scoped>
  .broadcast-message {
    color: white;
    display: flex;
    width: 100%;
    flex-direction: row;
    align-items: center;
    justify-content: center;

    &__text {
      text-align: center;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;

      & /deep/ p {
        margin: 0 !important;
      }
    }

    &__icon {
      min-width: 16px;
      color: inherit !important;
    }
  }
</style>
