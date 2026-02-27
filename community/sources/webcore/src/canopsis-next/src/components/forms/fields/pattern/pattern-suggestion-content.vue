<template>
  <v-layout
    :class="{ 'pattern-suggestions__content--difference': hasDifference }"
    class="pattern-suggestions__content gap-4"
    column
  >
    <c-alert
      :value="true"
      type="info"
    >
      <v-layout justify-space-between align-center>
        <v-layout class="gap-2" column>
          <div class="font-weight-regular">
            {{ $t('pattern.conditionsOptimized') }}
          </div>
          <ul class="pl-4 font-weight-regular">
            <li
              v-for="{ field, regexp } in optimizedFieldsRegexps"
              :key="field"
            >
              Infos.{{ field }}.Value <strong>{{ $t('common.regexp') }}</strong> {{ regexp }}
            </li>
          </ul>
        </v-layout>
        <v-layout justify-end>
          <v-btn color="success" @click="apply">
            <v-icon class="mr-2" color="white">
              check
            </v-icon>
            {{ $t('pattern.applyThisSuggestion') }}
          </v-btn>
        </v-layout>
      </v-layout>
    </c-alert>

    <div class="pattern-suggestions__pattern">
      <c-entity-patterns-field
        :patterns="suggestionPatternsForm"
        :attributes="entityAttributes"
        :readonly="true"
        :disabled="true"
        with-type
      />
    </div>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { PATTERNS_FIELDS, PATTERN_TYPES } from '@/constants';

import { patternToForm } from '@/helpers/entities/pattern/form';

export default {
  props: {
    suggestion: {
      type: Object,
      required: true,
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
    const hasDifference = computed(() => !!props.suggestion.difference);
    const suggestionPatternsForm = computed(() => patternToForm({
      type: PATTERN_TYPES.entity,
      [PATTERNS_FIELDS.entity]: props.suggestion[PATTERNS_FIELDS.entity],
    }));

    const apply = () => emit('apply');

    return {
      hasDifference,
      suggestionPatternsForm,

      apply,
    };
  },
};
</script>

<style lang="scss">
.pattern-suggestions__content {
  border-radius: 10px;
  border-top-left-radius: 0;
  padding: 16px;
  border: 2px solid var(--v-primary-base);

  &--difference {
    border-color: var(--v-info-base);
  }
}
</style>
