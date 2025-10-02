<template>
  <v-layout class="gap-2" column>
    <c-name-field
      v-field="form.name"
      :label="$t('common.name')"
      name="name"
      required
    />

    <c-description-field
      v-field="form.description"
      :label="$t('common.description')"
      :max-length="500"
      name="description"
    />

    <template-testing-data-type-field
      v-field="form.type"
      :label="$t('common.type')"
      :disabled="!isNew"
      name="type"
      required
    />

    <c-json-field
      v-field="form.body"
      ref="jsonFieldElement"
      :label="$t('common.value')"
      name="value"
      rows="11"
      validate-on="blur"
    />
    <v-flex
      v-if="isEventType"
      class="ml-3"
      align-self-end
    >
      <v-btn
        color="primary"
        outlined
        @click="showSetPreFilledTemplateModal"
      >
        {{ $t('templateTesting.usePreFilledTemplate') }}
      </v-btn>
    </v-flex>
    <c-information-block
      v-else
      :title="$tc('common.header', 2)"
      :help-text="$t('templateTesting.headersHelpText')"
      class="mb-2"
      help-icon="help"
      help-icon-color="grey darken-1"
    >
      <request-headers-field
        v-field="form.headers"
        name="headers"
      />
    </c-information-block>
  </v-layout>
</template>

<script>
import { computed, ref, nextTick } from 'vue';

import { MODALS, TEMPLATE_TESTING_DATA_TYPES, TEMPLATE_TESTING_DATA_EVENT_PRE_FILLED_TEMPLATE } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useModelField } from '@/hooks/form/model-field';

import RequestHeadersField from '@/components/forms/request/fields/request-headers-field.vue';

import TemplateTestingDataTypeField from './fields/template-testing-data-type-field.vue';

export default {
  inject: ['$validator'],
  components: {
    RequestHeadersField,
    TemplateTestingDataTypeField,
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
    isNew: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();
    const modals = useModals();
    const { updateField } = useModelField(props, emit);

    const jsonFieldElement = ref(null);

    const isEventType = computed(() => props.form.type === TEMPLATE_TESTING_DATA_TYPES.event);

    /**
     * Sets the pre-filled template data to the form body field and resets the JSON field
     */
    const setPreFilledTemplate = () => {
      updateField('body', TEMPLATE_TESTING_DATA_EVENT_PRE_FILLED_TEMPLATE);

      nextTick(() => jsonFieldElement.value?.reset?.());
    };

    /**
     * Shows confirmation modal before setting pre-filled template if form has existing content,
     * otherwise sets the template directly
     */
    const showSetPreFilledTemplateModal = () => {
      const bodyToCompare = props.form.body?.trim?.();

      if (!bodyToCompare || bodyToCompare === '{}') {
        setPreFilledTemplate();

        return;
      }

      modals.show({
        name: MODALS.confirmation,
        config: {
          title: t('templateTesting.usePreFilledTemplate'),
          text: t('templateTesting.usePreFilledTemplateWarning'),
          action: setPreFilledTemplate,
        },
      });
    };

    return {
      jsonFieldElement,

      isEventType,

      showSetPreFilledTemplateModal,
    };
  },
};
</script>
