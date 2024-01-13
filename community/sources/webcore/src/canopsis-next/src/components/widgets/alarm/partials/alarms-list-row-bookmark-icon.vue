<template>
  <v-tooltip top>
    <template #activator="{ on }">
      <v-icon
        class="cursor-default"
        v-on="on"
        size="22"
      >
        bookmark
      </v-icon>
    </template>
    <span>{{ $t('alarm.hasBookmark') }}</span>
  </v-tooltip>
</template>
