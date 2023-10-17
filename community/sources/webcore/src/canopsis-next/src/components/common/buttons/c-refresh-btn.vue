<template>
  <v-tooltip left>
    <template #activator="{ on }">
      <v-btn
        v-on="on"
        color="secondary"
        fab
        @click.stop="$emit('click', $event)"
      >
        <v-icon>refresh</v-icon>
      </v-btn>
    </template>
    <span>{{ $t('common.refresh') }}</span>
  </v-tooltip>
</template>
