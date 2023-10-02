<template lang="pug">
  div
    v-tooltip(left)
      template(#activator="{ on }")
        div(v-on="on") {{ version }}
      div
        div {{ $t('common.stack') }}:
          strong  {{ stack }}
        div {{ $t('common.edition') }}:
          strong  {{ edition }}
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
