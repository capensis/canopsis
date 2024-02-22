<template>
  <v-fade-transition>
    <v-tooltip
      v-if="pageScrolled"
      top
    >
      <template #activator="{ on }">
        <v-btn
          color="secondary lighten-2"
          fab
          dark
          v-on="on"
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
