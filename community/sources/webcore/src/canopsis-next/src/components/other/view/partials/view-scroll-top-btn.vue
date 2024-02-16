<template>
  <v-fade-transition>
    <v-tooltip
      v-if="pageScrolled"
      top
    >
      <template #activator="{ on }">
        <v-btn
          v-on="on"
          color="secondary lighten-2"
          fab
          dark
          @click="scrollToTop"
        >
          <v-icon>arrow_upward</v-icon>
        </v-btn>
      </template>
      <span>{{ $t('common.toTheTop') }}</span>
    </v-tooltip>
  </v-fade-transition>
</template>

<script>
import { scrollToTopMixin } from '@/mixins/scroll-to-top';

export default {
  mixins: [scrollToTopMixin],
};
</script>
