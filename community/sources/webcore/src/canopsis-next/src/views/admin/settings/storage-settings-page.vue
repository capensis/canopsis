<template>
  <c-page>
    <v-card-text class="pa-4">
      <storage-settings />
    </v-card-text>
  </c-page>
</template>

<script>
import StorageSettings from '@/components/other/storage-setting/storage-settings.vue';

export default {
  components: { StorageSettings },
};
</script>
