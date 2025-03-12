<template>
  <div class="secondary pa-3">
    <v-card>
      <v-card-text />
    </v-card>
  </div>
</template>
<script>

export default {
  props: {
    externalDataTable: {
      type: Object,
      required: true,
    },
  },
  setup() {
    return {};
  },
};
</script>
