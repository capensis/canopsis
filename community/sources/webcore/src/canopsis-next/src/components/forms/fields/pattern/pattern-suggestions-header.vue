<template>
  <v-layout
    justify-space-between
    align-center
  >
    <v-chip
      color="info"
      small
    >
      <strong>{{ $t('pattern.suggestions') }}</strong>
    </v-chip>
    <v-btn
      color="error"
      outlined
      @click="handleRejectAll"
    >
      <v-icon class="mr-2">
        close
      </v-icon>
      {{ $t('pattern.rejectAllSuggestions') }}
    </v-btn>
  </v-layout>
</template>

<script>
export default {
  setup(props, { emit }) {
    const handleRejectAll = () => {
      emit('reject-all');
    };

    return {
      handleRejectAll,
    };
  },
};
</script>
