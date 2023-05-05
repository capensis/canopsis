<template lang="pug">
  v-fade-transition
    v-tooltip(v-if="pageScrolled", top)
      template(#activator="{ on }")
        v-btn(
          v-on="on",
          color="secondary lighten-2",
          fab,
          dark,
          @click="scrollToTop"
        )
          v-icon arrow_upward
      span {{ $t('common.toTheTop') }}
</template>

<script>
import { scrollToTopMixin } from '@/mixins/scroll-to-top';

export default {
  mixins: [scrollToTopMixin],
};
</script>
