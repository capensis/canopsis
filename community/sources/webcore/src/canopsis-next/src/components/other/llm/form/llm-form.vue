<template>
  <v-layout class="gap-3" column>
    <c-enabled-field
      :value="form.enabled"
      with-background
      hide-details
      @change="updateEnabled"
    />

    <c-name-field
      v-field="form.name"
      :label="$t('llm.modelName')"
      :max-length="255"
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
      :replaceable="!isNew"
      name="api_key"
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
      :disabled="!form.enabled"
      hide-details
    />

    <div
      v-if="defaultLlm"
      class="text-caption text--secondary"
    >
      {{ $t('llm.currentDefaultModelLine') }} <strong>{{ defaultLlm.name }}</strong>
    </div>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { useModelField } from '@/hooks/form/model-field';

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
    defaultLlm: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props, { emit }) {
    const { updateModel } = useModelField(props, emit);
    const { items: models, itemsByName: modelsByName, pending: modelsPending } = useLlmModelsListForSelect();
    const emptyThinkingLevels = [];

    const thinkingLevels = computed(() => modelsByName.value[props.form.model]?.thinking_levels ?? emptyThinkingLevels);

    const updateEnabled = (value) => {
      const newForm = {
        ...props.form,
        enabled: value,
      };

      if (!newForm.enabled) {
        newForm.default = false;
      }

      updateModel(newForm);
    };

    return {
      models,
      modelsPending,
      thinkingLevels,
      updateEnabled,
    };
  },
};
</script>
