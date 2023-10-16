<template>
  <v-tooltip left="left">
    <v-btn
      slot="activator"
      color="secondary"
      data-test="refreshButton"
      fab="fab"
      dark="dark"
      icon="icon"
      @click.stop="$emit('click', $event)"
    >
      <v-icon>refresh</v-icon>
    </v-btn><span>{{ $t('common.refresh') }}</span>
  </v-tooltip>
</template>
