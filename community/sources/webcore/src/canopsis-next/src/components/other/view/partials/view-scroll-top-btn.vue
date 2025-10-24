<template>
  <v-fade-transition>
    <c-simple-tooltip
      v-if="pageScrolled"
      :message="$t('common.toTheTop')"
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
    </c-simple-tooltip>
  </v-fade-transition>
</template>

<script>
import { scrollToTopMixin } from '@/mixins/scroll-to-top';

export default {
  mixins: [scrollToTopMixin],
};
</script>
