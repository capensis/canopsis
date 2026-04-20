<template>
  <v-layout
    class="pattern-suggestions__wrapper gap-2 pa-3"
    column
  >
    <pattern-suggestions-header @reject:all="rejectAll" />

    <pattern-suggestions-list
      :suggestions="suggestions"
      :entity-attributes="entityAttributes"
      :optimized-fields-regexps="optimizedFieldsRegexps"
      @apply="applySuggestion"
      @show:entities-comparison="showEntitiesComparisonModal"
    />
  </v-layout>
</template>

<script>
import PatternSuggestionsHeader from './pattern-suggestions-header.vue';
import PatternSuggestionsList from './pattern-suggestions-list.vue';

export default {
  components: {
    PatternSuggestionsHeader,
    PatternSuggestionsList,
  },
  props: {
    suggestions: {
      type: Array,
      default: () => [],
    },
    entityAttributes: {
      type: Array,
      default: () => [],
    },
    optimizedFieldsRegexps: {
      type: Array,
      default: () => [],
    },
  },
  setup(props, { emit }) {
    const rejectAll = () => emit('reject:all');
    const applySuggestion = index => emit('apply:suggestion', index);
    const showEntitiesComparisonModal = suggestion => emit('show:entities-comparison', suggestion);

    return {
      rejectAll,
      applySuggestion,
      showEntitiesComparisonModal,
    };
  },
};
</script>

<style lang="scss">
.pattern-suggestions__wrapper {
  border-radius: 10px;
  border: 2px solid var(--v-info-base);
}
</style>
