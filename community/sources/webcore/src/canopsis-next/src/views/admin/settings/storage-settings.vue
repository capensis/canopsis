<template>
  <v-container>
    <c-page-header />
    <v-card>
      <v-card-text>
        <storage-settings />
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script>

import StorageSettings from '@/components/other/storage-setting/storage-settings.vue';

export default {
  components: {
    StorageSettings,
  },
};
</script>
