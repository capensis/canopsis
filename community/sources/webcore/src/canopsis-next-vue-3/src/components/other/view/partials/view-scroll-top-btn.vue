<template>
  <v-fade-transition>
    <div v-if="pageScrolled" class="view-scroll-top-btn">
      <v-tooltip location="top">
        <template #activator="{ props }">
          <v-btn
            v-bind="props"
            color="secondary-lighten-2"
            size="x-large"
            icon
            @click="scrollToTop"
          >
            <v-icon icon="arrow_upward" color="white" />
          </v-btn>
        </template>

        <span>{{ $t('common.toTheTop') }}</span>
      </v-tooltip>
    </div>
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

<style lang="scss" scoped>
.view-scroll-top-btn {
  display: inline-flex;
}
</style>
