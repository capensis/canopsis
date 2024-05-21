<template>
  <v-tooltip top disable-resize>
    <template #activator="{ on }">
      <v-icon
        class="cursor-default"
        size="22"
        v-on="on"
      >
        bookmark
      </v-icon>
    </template>
    <span>{{ $t('alarm.hasBookmark') }}</span>
  </v-tooltip>
</template>
