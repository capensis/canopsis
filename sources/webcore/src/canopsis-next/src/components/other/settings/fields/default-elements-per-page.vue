<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.defaultNumberOfElementsPerPage')}}
    v-container
      v-text-field(
      :placeholder="$t('settings.elementsPerPage')",
      type="number"
      )
</template>
