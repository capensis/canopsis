<template>
  <v-layout class="gap-3" column>
    <c-name-field
      v-field="form.name"
      autofocus
      required
    />
    <c-form-block>
      <c-form-block-row :label="$t('common.category')">
        <c-entity-category-field
          v-field="form.category"
          addable
          required
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('entity.availabilityState')">
        <c-alarm-state-field
          v-field="form.sli_avail_state"
          :label="$t('entity.availabilityState')"
          required
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('entity.impact')">
        <c-impact-level-field
          v-field="form.impact_level"
          required
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.coordinates')">
        <c-coordinates-field
          v-field="form.coordinates"
          row
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('service.outputTemplate')">
        <text-editor-field
          v-validate="'required'"
          v-field="form.output_template"
          :label="$t('service.outputTemplate')"
          :error-messages="errors.collect('output_template')"
          :variables="templateVars.output"
          name="output_template"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('stateSetting.title')">
        <entity-state-setting
          :form="form"
          :preparer="prepareStateSettingForm"
        />
      </c-form-block-row>
    </c-form-block>
  </v-layout>
</template>

<script>
import TextEditorField from '@/components/forms/fields/text-editor-field.vue';
import EntityStateSetting from '@/components/other/state-setting/entity-state-setting.vue';

export default {
  inject: ['$validator'],
  components: {
    TextEditorField,
    EntityStateSetting,
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
    prepareStateSettingForm: {
      type: Function,
      default: data => data,
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
