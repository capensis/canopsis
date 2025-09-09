<template>
  <v-form data-vv-scope="test-variables">
    <v-layout class="gap-4" column>
      <template-testing-test-variables-form :fields="fields" />
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
import { ref } from 'vue';

import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

import { useValidator } from '@/hooks/validator/validator';

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
      default: TEMPLATE_TESTING_TEST_TYPES.eventFilter,
    },
  },
  setup(props, { emit }) {
    const isGeneralFormValid = ref(true);

    const validator = useValidator();

    const saveAsNew = async () => {
      const isValid = await validator.validateAll('test-variables');

      if (isValid) {
        emit('saveAsNew');
      }
    };
    const save = () => {};
    const runTest = async () => {
      isGeneralFormValid.value = await validator.validateAll('general');
    };

    const updateMainForm = (newMainForm) => {
      emit('input', newMainForm);
    };

    return {
      isGeneralFormValid,
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
