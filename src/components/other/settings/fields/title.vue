<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('common.title')}}
    v-container
      v-text-field(:placeholder="$t('settings.widgetTitle')")
</template>
