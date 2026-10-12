<template>
  <c-simple-tooltip :content="$t('alarm.hasBookmark')" top>
    <template #activator="{ on }">
      <v-icon
        class="cursor-default"
        size="22"
        v-on="on"
      >
        bookmark
      </v-icon>
    </template>
  </c-simple-tooltip>
</template>
