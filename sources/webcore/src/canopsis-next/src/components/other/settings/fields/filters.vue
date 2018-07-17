<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.filters')}}
    v-container
      v-select(:label="$t('settings.selectAFilter')")
</template>
