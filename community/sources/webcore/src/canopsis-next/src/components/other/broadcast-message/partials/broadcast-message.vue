<template>
  <div
    :style="{ backgroundColor: color }"
    :title="message"
    class="broadcast-message pa-2"
  >
    <c-compiled-template
      :template="message"
      class="broadcast-message__text"
      parent-element="span"
    />
    <div
      class="broadcast-message__actions"
      title=""
    >
      <slot name="actions" />
    </div>
  </div>
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
    position: relative;
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

    &__actions {
      display: flex;
      align-items: center;
      position: absolute;
      right: 0;
      top: 0;
      bottom: 0;
    }
  }
</style>
