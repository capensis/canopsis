<template lang="pug">
  div.broadcast-message.pa-2(
    :style="{ backgroundColor: color }",
    :title="message"
  )
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
      width: 100%;
      text-align: center;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;

      & ::v-deep p {
        margin: 0 !important;
      }
    }
  }
</style>
