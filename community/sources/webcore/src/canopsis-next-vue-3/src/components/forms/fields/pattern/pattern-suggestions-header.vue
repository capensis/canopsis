<template>
  <c-layout
    justify-space-between
    align-center
  >
    <v-chip
      color="info darken-1"
      small
    >
      <strong>{{ $t('pattern.suggestions') }}</strong>
    </v-chip>
    <v-btn
      color="error"
      outlined
      @click="rejectAll"
    >
      <v-icon class="mr-2">
        close
      </v-icon>
      {{ $t('pattern.rejectAllSuggestions') }}
    </v-btn>
  </c-layout>
</template>

<script>
export default {
  setup(props, { emit }) {
    const rejectAll = () => emit('reject:all');

    return {
      rejectAll,
    };
  },
};
</script>
