<template>
  <v-layout
    class="ai-chat-greeting gap-3"
    column
    align-center
    justify-end
  >
    <v-icon
      color="grey"
      x-large
    >
      $vuetify.icons.robot
    </v-icon>
    <span class="text-h5 font-weight-thin grey--text">{{ $t('llm.chat.howCanIHelp') }}</span>
  </v-layout>
</template>

<script>
export default {};
</script>
