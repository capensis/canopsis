<template lang="pug">
  div.broadcast-message.pa-2(
    :style="{ backgroundColor: color }",
    :title="message"
  )
    span
      v-icon.mr-1(small) $vuetify.icons.bullhorn
      | {{ message }}
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
    text-align: center;
    width: 100%;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;

    & /deep/ svg {
      color: inherit;
    }
  }
</style>
