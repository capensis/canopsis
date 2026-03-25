<template>
  <v-layout class="gap-3" column>
    <c-enabled-field
      v-field="form.enabled"
      hide-details
    />

    <c-name-field
      v-field="form.name"
      :label="$t('llm.modelName')"
      name="llm_name"
      required
    />

    <v-text-field
      v-field="form.type"
      :label="$t('llm.modelType')"
      disabled
      readonly
    />

    <c-password-field
      v-field="form.api_key"
      :label="$t('llm.apiKey')"
      :required="isNew"
      :placeholder="$t('llm.apiKeyPlaceholder')"
      name="api_key"
      replaceable
      visibility
    />

    <v-layout class="gap-3">
      <v-flex xs6>
        <llm-model-field
          v-field="form.model"
          :items="models"
          :loading="modelsPending"
          name="model"
          required
        />
      </v-flex>
      <v-flex xs6>
        <llm-thinking-level-field
          v-field="form.thinking_level"
          :items="thinkingLevels"
          name="thinking_level"
        />
      </v-flex>
    </v-layout>

    <c-enabled-field
      v-field="form.default"
      :label="$t('llm.isDefaultModel')"
    />

    <div
      v-if="currentDefaultModelName"
      class="text-caption text--secondary mt-2"
    >
      {{ $t('llm.currentDefaultModelLine', { name: currentDefaultModelName }) }}
    </div>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { useLlmModelsListForSelect } from '@/components/other/llm/hooks/llm-models-list';

import LlmModelField from './fields/llm-model-field.vue';
import LlmThinkingLevelField from './fields/llm-thinking-level-field.vue';

export default {
  inject: ['$validator'],
  components: {
    LlmModelField,
    LlmThinkingLevelField,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    isNew: {
      type: Boolean,
      default: false,
    },
    currentDefaultModelName: {
      type: String,
      default: '',
    },
  },
  setup(props) {
    const { items: models, itemsByName: modelsByName, pending: modelsPending } = useLlmModelsListForSelect();
    const emptyThinkingLevels = [];

    const thinkingLevels = computed(() => modelsByName.value[props.form.model]?.thinking_levels ?? emptyThinkingLevels);

    return {
      models,
      modelsPending,
      thinkingLevels,
    };
  },
};
</script>
