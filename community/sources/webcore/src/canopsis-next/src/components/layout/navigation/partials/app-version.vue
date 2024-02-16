<template>
  <div>
    <v-tooltip left>
      <template #activator="{ on }">
        <div v-on="on">
          {{ version }}
        </div>
      </template>
      <div>
        <div>{{ $t('common.stack') }}:<strong> {{ stack }}</strong></div>
        <div>{{ $t('common.edition') }}:<strong> {{ edition }}</strong></div>
      </div>
    </v-tooltip>
  </div>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

const { mapGetters } = createNamespacedHelpers('info');

export default {
  computed: {
    ...mapGetters(['version', 'stack', 'edition']),
  },
};
</script>
