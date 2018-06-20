<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.filterOnOpenResolved')}}
    v-container
      v-layout
        v-checkbox(
        :label="$t('settings.open')",
        hide-details
        )
        v-checkbox(
        :label="$t('settings.resolved')",
        hide-details
        )
</template>
