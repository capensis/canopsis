<template>
  <v-card class="ma-3" variant="flat">
    <v-card-text>
      <v-data-table
        :items="info.infos"
        :headers="infosTableHeaders"
        :no-data-text="$t('common.noData')"
        density="compact"
      />
    </v-card-text>
  </v-card>
</template>

<script>
export default {
  props: {
    info: {
      type: Object,
      required: true,
    },
  },

  computed: {
    infosTableHeaders() {
      return [
        { title: this.$t('common.name'), key: 'name' },
        { title: this.$t('common.value'), key: 'value' },
      ];
    },
  },
};
</script>
