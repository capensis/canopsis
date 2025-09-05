<template>
  <v-layout class="gap-4" column>
    <template-testing-test-variables-form :main-form="form" @update:mainForm="updateMainForm" />
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
</template>

<script>
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
  },
  setup(props, { emit }) {
    const saveAsNew = () => {};
    const save = () => {};
    const runTest = () => {};

    const updateMainForm = (newMainForm) => {
      emit('input', newMainForm);
    };

    return {
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
