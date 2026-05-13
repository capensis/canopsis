<template>
  <v-layout class="gap-4" column>
    <v-layout ref="generalSectionRef" class="gap-2" column>
      <c-name-field
        v-field="form.name"
        :autofocus="isDisabledIdField"
        required
      />
      <c-form-block>
        <c-form-block-row :label="$t('common.description')">
          <c-description-field
            v-field="form.description"
            required
          />
        </c-form-block-row>
        <c-form-block-row :label="$t('common.disableDuringPeriods')">
          <c-disable-during-periods-field v-field="form.disable_during_periods" />
        </c-form-block-row>
        <c-form-block-row :label="$t('modals.createDynamicInfo.infosSection.title')">
          <c-label>{{ $t('modals.createDynamicInfo.infosSection.title') }}</c-label>
          <dynamic-info-infos-form
            v-field="form.infos"
            :variables="templateVars.value"
            :copy-variables="copyVars.value"
          />
        </c-form-block-row>
      </c-form-block>
    </v-layout>
  </v-layout>
</template>

<script>

import { useComponentInstance } from '@/hooks/vue';
import { useValidationElementChildren } from '@/hooks/validator/validation-element-children';

import DynamicInfoInfosForm from './dynamic-info-infos-form.vue';

export default {
  inject: ['$validator'],
  components: {
    DynamicInfoInfosForm,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    isDisabledIdField: {
      type: Boolean,
      default: false,
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
    copyVars: {
      type: Object,
      default: () => ({}),
    },
  },
  setup() {
    const instance = useComponentInstance();

    const { hasChildrenError: hasGeneralError } = useValidationElementChildren(instance);

    return {
      hasAnyError: hasGeneralError,
    };
  },
};
</script>
