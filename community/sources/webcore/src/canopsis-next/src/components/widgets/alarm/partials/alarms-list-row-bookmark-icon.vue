<template lang="pug">
  v-tooltip(top)
    template(#activator="{ on }")
      v-icon.cursor-default(v-on="on", size="22") bookmark
    span {{ $t('alarm.hasBookmark') }}
</template>
