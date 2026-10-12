<template>
  <c-layout
    class="ai-chat-greeting"
    gap="3"
    align="center"
    justify="end"
    column
  >
    <v-icon class="ai-chat-greeting__icon" color="grey" size="64">
      $vuetify.icons.robot
    </v-icon>
    <span class="text-grey text-h5 font-weight-thin">
      {{ $t('llm.chat.howCanIHelp') }}
    </span>
  </c-layout>
</template>

<script>
export default {};
</script>

<style lang="scss" scoped>
.ai-chat-greeting {
  align-items: center;
  justify-content: flex-end;
  min-width: 0;
  height: 100%;
  text-align: center;

  &__icon {
    align-self: center;
    display: inline-flex;
    width: 64px;
    height: 64px;
    font-size: 64px;
    justify-content: center;

    :deep(svg) {
      width: 64px;
      height: 64px;
    }
  }

  .text-h5 {
    font-size: 24px !important;
    line-height: 32px;
  }
}
</style>
