<template>
  <c-simple-tooltip :content="$t('alarm.hasBookmark')" top>
    <template #activator="{ props }">
      <v-icon
        class="cursor-default"
        size="20"
        v-bind="props"
      >
        bookmark
      </v-icon>
    </template>
  </c-simple-tooltip>
</template>
