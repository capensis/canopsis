<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.moreInfosModal')}}
    v-container
      v-textarea.pa-0(hide-details, @input="$emit('input', $event)")
</template>

