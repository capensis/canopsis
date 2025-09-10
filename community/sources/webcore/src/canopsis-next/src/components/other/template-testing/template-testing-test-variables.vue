<template>
  <v-form data-vv-scope="test-variables">
    <v-layout class="gap-4" column>
      <template-testing-test-variables-form
        v-model="validationForm"
        :fields="fields"
        :type="type"
      />
      <c-alert :value="!isGeneralFormValid" type="error">
        {{ $t('templateTesting.mainFormHasErrors') }}
      </c-alert>
      <v-layout class="gap-2" justify-end>
        <v-btn
          color="secondary"
          outlined
          @click="saveAsNew"
        >
          {{ $t('templateTesting.saveTestAsNew') }}
        </v-btn>
        <v-btn
          color="secondary"
          @click="save"
        >
          {{ $t('templateTesting.saveTest') }}
        </v-btn>
        <v-btn
          :loading="running"
          :disabled="running"
          color="primary"
          @click="runTest"
        >
          {{ $t('templateTesting.runTest') }}
        </v-btn>
      </v-layout>
    </v-layout>
  </v-form>
</template>

<script>
import { computed, ref, watch } from 'vue';

import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

import { formToService } from '@/helpers/entities/service/form';
import { formToEventFilter } from '@/helpers/entities/event-filter/rule/form';
import { formToLinkRule } from '@/helpers/entities/link/form';
import { formToScenario } from '@/helpers/entities/scenario/form';
import { formToWidget } from '@/helpers/entities/widget/form';
import { formToDeclareTicketRule } from '@/helpers/entities/declare-ticket/rule/form';
import { formToDynamicInfo } from '@/helpers/entities/dynamic-info/rule/form';
import { formToRemediationInstruction } from '@/helpers/entities/remediation/instruction/form';
import { formToRemediationJob } from '@/helpers/entities/remediation/job/form';
import { formToMetaAlarmRule } from '@/helpers/entities/meta-alarm/rule/form';
import { formToTemplateTestingTestValidateForm } from '@/helpers/entities/template-testing-test/form';

import { usePendingHandler } from '@/hooks/query/pending';
import { useValidator } from '@/hooks/validator/validator';
import { useTemplateValidation } from '@/hooks/store/modules/template-validation';

import TemplateTestingTestVariablesForm from './template-testing-test-variables-form.vue';

export default {
  components: {
    TemplateTestingTestVariablesForm,
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
    fields: {
      type: Array,
      default: () => [],
    },
    isNew: {
      type: Boolean,
      default: false,
    },
    type: {
      type: Number,
      required: false,
    },
  },
  setup(props, { emit }) {
    const validationForm = ref([]);
    const isGeneralFormValid = ref(true);

    const validator = useValidator();

    const {
      validateEntityServices,
      validateEventFilters,
      validateScenarios,
      validateLinkRules,
      validateWidgets,
      validateDeclareTicketRules,
      validateDynamicInfos,
      validateInstructions,
      validateJobs,
      validateMetaAlarmRules,
    } = useTemplateValidation();

    watch(() => props.form, (newForm) => {
      validationForm.value = formToTemplateTestingTestValidateForm(newForm, props.type);
    }, { immediate: true });

    const validateHandler = computed(() => ({
      [TEMPLATE_TESTING_TEST_TYPES.eventFilter]: validateEventFilters,
      [TEMPLATE_TESTING_TEST_TYPES.linkRule]: validateLinkRules,
      [TEMPLATE_TESTING_TEST_TYPES.scenario]: validateScenarios,
      [TEMPLATE_TESTING_TEST_TYPES.widget]: validateWidgets,
      [TEMPLATE_TESTING_TEST_TYPES.declareTicketRule]: validateDeclareTicketRules,
      [TEMPLATE_TESTING_TEST_TYPES.dynamicInfo]: validateDynamicInfos,
      [TEMPLATE_TESTING_TEST_TYPES.instruction]: validateInstructions,
      [TEMPLATE_TESTING_TEST_TYPES.job]: validateJobs,
      [TEMPLATE_TESTING_TEST_TYPES.metaAlarmRule]: validateMetaAlarmRules,
    })[props.type] ?? validateEntityServices);

    const formToRequest = computed(() => ({
      [TEMPLATE_TESTING_TEST_TYPES.eventFilter]: formToEventFilter,
      [TEMPLATE_TESTING_TEST_TYPES.linkRule]: formToLinkRule,
      [TEMPLATE_TESTING_TEST_TYPES.scenario]: formToScenario,
      [TEMPLATE_TESTING_TEST_TYPES.widget]: formToWidget,
      [TEMPLATE_TESTING_TEST_TYPES.declareTicketRule]: formToDeclareTicketRule,
      [TEMPLATE_TESTING_TEST_TYPES.dynamicInfo]: formToDynamicInfo,
      [TEMPLATE_TESTING_TEST_TYPES.instruction]: formToRemediationInstruction,
      [TEMPLATE_TESTING_TEST_TYPES.job]: formToRemediationJob,
      [TEMPLATE_TESTING_TEST_TYPES.metaAlarmRule]: formToMetaAlarmRule,
    })[props.type] ?? formToService);

    const saveAsNew = async () => {
      const isValid = await validator.validateAll('test-variables');

      if (isValid) {
        emit('saveAsNew');
      }
    };

    const { pending: running, handler: runTest } = usePendingHandler(async () => {
      isGeneralFormValid.value = await validator.validateAll('general');

      if (!isGeneralFormValid.value) {
        return;
      }

      await validateHandler.value({
        data: {
          rule: formToRequest.value(props.form),
        },
      });
    });

    const save = () => {};

    const updateMainForm = (newMainForm) => {
      emit('input', newMainForm);
    };

    return {
      isGeneralFormValid,
      validationForm,
      running,
      saveAsNew,
      save,
      runTest,
      updateMainForm,
    };
  },
};
</script>

<style lang="scss" scoped>
.variables-item {
  display: flex;
  align-items: center;
  gap: 4px;

  &__name {
    font-weight: 500;
  }

  &__separator {
    color: rgba(0, 0, 0, 0.6);
  }

  &__value {
    flex: 1;
    min-width: 0;
  }

  &__empty-object {
    color: rgba(0, 0, 0, 0.6);
  }
}

.theme--dark {
  .variables-item {
    &__separator {
      color: rgba(255, 255, 255, 0.7);
    }

    &__empty-object {
      color: rgba(255, 255, 255, 0.7);
    }
  }
}
</style>
