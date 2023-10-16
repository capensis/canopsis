<template>
  <v-tooltip left>
    <v-btn
      slot="activator"
      color="secondary"
      data-test="refreshButton"
      fab
      dark
      icon
      @click.stop="$emit('click', $event)"
    >
      <v-icon>refresh</v-icon>
    </v-btn>
    <span>{{ $t('common.refresh') }}</span>
  </v-tooltip>
</template>
