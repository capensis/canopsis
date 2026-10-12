<template>
  <v-fade-transition>
    <c-simple-tooltip
      v-if="pageScrolled"
      :content="$t('common.toTheTop')"
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
import { useScrollToTop } from '@/hooks/scroll-to-top';

export default {
  setup() {
    const { pageScrolled, scrollToTop } = useScrollToTop();

    return {
      pageScrolled,
      scrollToTop,
    };
  },
};
</script>
