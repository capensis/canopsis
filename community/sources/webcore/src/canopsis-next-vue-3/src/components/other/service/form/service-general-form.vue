<template>
  <c-layout gap="3" column>
    <c-name-field
      :model-value="form.name"
      autofocus
      required
      @update:model-value="updateField('name', $event)"
    />
    <c-form-block>
      <c-form-block-row :label="$t('common.category')">
        <c-entity-category-field
          :category="form.category"
          addable
          required
          @update:category="updateField('category', $event)"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('entity.availabilityState')">
        <c-alarm-state-field
          :value="form.sli_avail_state"
          :label="$t('entity.availabilityState')"
          required
          @update:value="updateField('sli_avail_state', $event)"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('entity.impact')">
        <c-impact-level-field
          :value="form.impact_level"
          required
          @update:value="updateField('impact_level', $event)"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.coordinates')">
        <c-coordinates-field
          :value="form.coordinates"
          row
          @update:value="updateField('coordinates', $event)"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('service.outputTemplate')">
        <text-editor-field
          :model-value="form.output_template"
          :label="$t('service.outputTemplate')"
          :variables="templateVars.output"
          name="output_template"
          required
          @update:model-value="updateField('output_template', $event)"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('stateSetting.title')">
        <entity-state-setting
          :form="form"
          :preparer="prepareStateSettingForm"
        />
      </c-form-block-row>
    </c-form-block>
  </c-layout>
</template>

<script>
import { useModelField } from '@/hooks/form/model-field';

import TextEditorField from '@/components/forms/fields/text-editor-field.vue';
import EntityStateSetting from '@/components/other/state-setting/entity-state-setting.vue';

export default {
  components: {
    TextEditorField,
    EntityStateSetting,
  },

  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    prepareStateSettingForm: {
      type: Function,
      default: data => data,
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },

  emits: ['update:form'],

  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);

    return {
      updateField,
    };
  },
};
</script>
