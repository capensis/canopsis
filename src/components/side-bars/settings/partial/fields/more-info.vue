<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.moreInfosModal')}}
    v-container
      v-text-field(class="pa-0", textarea, hide-details)
</template>
