<template>
  <div v-if="suggestions.length">
    <v-layout
      class="mb-4"
      align-center
      wrap
    >
      <v-tabs
        v-model="activeSuggestionIndex"
        slider-color="success"
        class="pattern-suggestions__tabs"
        hide-slider
      >
        <v-tab
          v-for="(suggestion, index) in suggestions"
          :key="index"
        >
          <pattern-suggestion-tab
            :number="index + 1"
            :entities-count="suggestion.entitiesCount"
            :active="index === activeSuggestionIndex"
          />
        </v-tab>
      </v-tabs>
    </v-layout>

    <v-tabs-items v-model="activeSuggestionIndex">
      <v-tab-item
        v-for="(suggestion, index) in suggestions"
        :key="index"
        :value="index"
      >
        <pattern-suggestion-content
          :suggestion="suggestion"
          :index="index"
          @apply="handleApply"
        />
      </v-tab-item>
    </v-tabs-items>
  </div>
</template>

<script>
import { ref } from 'vue';

import PatternSuggestionTab from './pattern-suggestion-tab.vue';
import PatternSuggestionContent from './pattern-suggestion-content.vue';

export default {
  components: {
    PatternSuggestionTab,
    PatternSuggestionContent,
  },
  props: {
    suggestions: {
      type: Array,
      default: () => [],
    },
  },
  setup() {
    const activeSuggestionIndex = ref(0);

    return {
      activeSuggestionIndex,
    };
  },
  methods: {
    handleApply(suggestion, index) {
      this.$emit('apply', suggestion, index);
    },
  },
};
</script>

<style lang="scss">
.pattern-suggestions__tabs {
  .v-tab {
    min-width: 0;
    padding: 0;
    overflow: hidden;

    .pattern-suggestions__tabs__number {
      height: 100%;
      padding: 0 16px;
      background-color: var(--v-info-background-base);
      color: var(--v-text-light-primary);
      transition: background-color 0.3s ease, color 0.3s ease;
    }

    &, &:before {
      border-top-left-radius: 10px;
      border-top-right-radius: 10px;
    }

    &--active .pattern-suggestions__tabs__number {
      background-color: var(--v-success-base);
      color: var(--v-text-dark-primary);
    }
  }

  .v-tabs-bar__content {
    column-gap: 12px;
  }

  &__found {
    overflow-wrap: normal;
    white-space: nowrap;
  }
}
</style>
