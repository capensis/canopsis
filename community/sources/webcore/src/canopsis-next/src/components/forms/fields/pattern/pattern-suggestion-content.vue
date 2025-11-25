<template>
  <v-layout
    class="c-pattern-suggestions__content gap-4"
    column
  >
    <c-alert
      :value="true"
      type="info"
    >
      <div class="mb-2 font-weight-regular">
        {{ $t('pattern.conditionsOptimized') }}
      </div>
      <ul class="mb-0 pl-4 font-weight-regular">
        <li
          v-for="(condition, conditionIndex) in originalConditions"
          :key="conditionIndex"
        >
          {{ condition }}
        </li>
      </ul>
      <v-layout
        class="mt-3"
        justify-end
      >
        <v-btn
          color="success"
          @click="handleApply"
        >
          <v-icon
            class="mr-2"
            color="white"
          >
            check
          </v-icon>
          {{ $t('pattern.applyThisSuggestion') }}
        </v-btn>
      </v-layout>
    </c-alert>

    <div class="c-pattern-suggestions__pattern">
      <pattern-groups-field
        :groups="optimizedPattern.groups || []"
        :attributes="attributes"
        :readonly="true"
        :disabled="true"
      />
    </div>
  </v-layout>
</template>

<script>
import PatternGroupsField from './pattern-groups-field.vue';

export default {
  components: { PatternGroupsField },
  props: {
    suggestion: {
      type: Object,
      required: true,
    },
    attributes: {
      type: Array,
      default: () => [],
    },
    index: {
      type: Number,
      required: true,
    },
  },
  computed: {
    originalConditions() {
      return this.suggestion.originalConditions || [];
    },
    optimizedPattern() {
      return this.suggestion.optimizedPattern || {};
    },
  },
  methods: {
    handleApply() {
      this.$emit('apply', this.suggestion, this.index);
    },
  },
};
</script>
