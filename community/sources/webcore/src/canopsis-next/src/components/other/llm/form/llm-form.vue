<template>
  <c-layout gap="3" column>
    <c-enabled-field
      :model-value="form.enabled"
      with-background
      hide-details
      @update:model-value="updateEnabled"
    />

    <c-name-field
      :model-value="form.name"
      :label="$t('llm.modelName')"
      :max-length="FORM_MAX_STRING_LENGTH"
      name="llm_name"
      required
      @update:model-value="updateField('name', $event)"
    />

    <c-form-block>
      <c-form-block-row :label="$t('llm.modelType')">
        <llm-model-type-field
          :model-value="form.type"
          :disabled="!isNew"
          name="type"
          required
          @update:model-value="updateType"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('llm.apiKey')">
        <c-password-field
          :model-value="form.api_key"
          :label="$t('llm.apiKey')"
          :required="isNew"
          :placeholder="$t('llm.apiKeyPlaceholder')"
          :replaceable="!isNew"
          name="api_key"
          visibility
          @update:model-value="updateField('api_key', $event)"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('llm.model')">
        <v-row>
          <v-col cols="6" class="pa-0">
            <llm-model-field
              :model-value="form.model"
              :items="models"
              :loading="modelsPending"
              name="model"
              required
              @update:model-value="updateField('model', $event)"
            />
          </v-col>
          <v-col v-if="thinkingLevels.length" cols="6" class="pa-0">
            <llm-thinking-level-field
              :model-value="form.thinking_level"
              :items="thinkingLevels"
              name="thinking_level"
              @update:model-value="updateField('thinking_level', $event)"
            />
          </v-col>
        </v-row>
      </c-form-block-row>

      <c-form-block-row :label="$t('llm.isDefaultModel')" align-center>
        <c-layout column>
          <c-enabled-field
            :model-value="form.default"
            :label="$t('llm.isDefaultModel')"
            :disabled="!form.enabled"
            hide-details
            no-margin
            @update:model-value="updateField('default', $event)"
          />
          <div v-if="defaultLlm" class="text-caption text-medium-emphasis mt-1">
            {{ $t('llm.currentDefaultModelLine') }} <strong>{{ defaultLlm.name }}</strong>
          </div>
        </c-layout>
      </c-form-block-row>
    </c-form-block>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { FORM_MAX_STRING_LENGTH } from '@/constants';

import { useModelField } from '@/hooks/form/model-field';

import { useLlmModelsListForSelect } from '@/components/other/llm/hooks/llm-models-list';

import LlmModelField from './fields/llm-model-field.vue';
import LlmModelTypeField from './fields/llm-model-type-field.vue';
import LlmThinkingLevelField from './fields/llm-thinking-level-field.vue';

export default {
  components: {
    LlmModelField,
    LlmModelTypeField,
    LlmThinkingLevelField,
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
      default: null,
    },
  },

  emits: ['update:form'],

  setup(props, { emit }) {
    const { updateModel, updateField } = useModelField(props, emit);
    const type = computed(() => props.form.type);
    const {
      items: models,
      itemsByName: modelsByName,
      pending: modelsPending,
    } = useLlmModelsListForSelect(type);
    const emptyThinkingLevels = [];

    const thinkingLevels = computed(
      () => modelsByName.value[props.form.model]?.thinking_levels ?? emptyThinkingLevels,
    );

    const updateType = (value) => {
      updateModel({
        ...props.form,
        type: value,
        model: '',
        thinking_level: '',
      });
    };

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
      FORM_MAX_STRING_LENGTH,
      models,
      modelsPending,
      thinkingLevels,

      updateEnabled,
      updateField,
      updateType,
    };
  },
};
</script>
