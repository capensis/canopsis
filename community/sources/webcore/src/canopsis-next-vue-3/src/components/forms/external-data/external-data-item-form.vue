<template>
  <v-card>
    <v-card-text>
      <c-layout gap="2" column>
        <c-layout align="center" gap="2">
          <c-text-field
            :model-value="form.reference"
            :label="$t('externalData.fields.reference')"
            :name="referenceFieldName"
            :disabled="disabled"
            required
            @update:model-value="updateField('reference', $event)"
          >
            <template #append-inner="">
              <c-help-icon :text="$t('externalData.tooltips.reference')" icon="help" left />
            </template>
          </c-text-field>
          <c-select
            :model-value="form.type"
            :items="availableTypes"
            :label="$t('common.type')"
            :disabled="disabled"
            item-title="text"
            item-value="value"
            @update:model-value="updateField('type', $event)"
          />
          <v-btn
            v-if="!disabled"
            class="mr-0"
            density="compact"
            size="small"
            variant="text"
            icon
            @click="remove"
          >
            <v-icon color="error">
              delete
            </v-icon>
          </v-btn>
        </c-layout>
        <external-data-table-form
          v-if="isTableType"
          :form="form"
          :name="name"
          :disabled="disabled"
          :variables="variables"
          :optionally="optionally"
          @update:form="updateModel"
        />
        <request-form
          v-else
          :form="form.request"
          :name="`${name}.request`"
          :disabled="disabled"
          :payload-variables="variables"
          :url-variables="variables"
          @update:form="updateField('request', $event)"
        />
        <c-alert v-if="serverErrorMessages.length" type="error">
          {{ serverErrorMessages.join('\n') }}
        </c-alert>
      </c-layout>
    </v-card-text>
  </v-card>
</template>

<script>
import { computed, onBeforeUnmount, onMounted } from 'vue';

import { EXTERNAL_DATA_TYPES } from '@/constants';

import { isTableExternalDataType } from '@/helpers/entities/shared/external-data/entity';

import { useFormContext, useFormErrorsMap } from '@/hooks/form';
import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';
import { useComponentInstance } from '@/hooks/vue';

import RequestForm from '@/components/forms/request/request-form.vue';

import ExternalDataTableForm from './external-data-table-form.vue';

export default {
  emits: ['update:form'],
  components: { RequestForm, ExternalDataTableForm },
  props: {
    form: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      required: true,
    },
    serverErrorName: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    types: {
      type: Array,
      default: () => [],
    },
    variables: {
      type: Array,
      default: () => [],
    },
    optionally: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();
    const formContext = useFormContext();
    const instance = useComponentInstance();
    const { collectErrors } = useFormErrorsMap();
    const { updateField, updateModel } = useModelField(props, emit);

    const availableTypes = computed(() => props.types.length
        ? props.types
        : Object.values(EXTERNAL_DATA_TYPES).map(type => ({
            text: t(`externalData.types.${type}`),
            value: type,
          })),
    );

    const isTableType = computed(() => isTableExternalDataType(props.form.type));
    const referenceFieldName = computed(() => `${props.name}.reference`);
    const serverErrorMessages = computed(() => collectErrors(props.serverErrorName));
    const remove = () => emit('remove', props.form);

    onMounted(() => {
      if (props.serverErrorName) {
        formContext.addFieldInstance(props.serverErrorName, instance);
      }
    });

    onBeforeUnmount(() => {
      if (props.serverErrorName) {
        formContext.removeFieldInstance(props.serverErrorName);
      }
    });

    return {
      availableTypes,
      isTableType,
      referenceFieldName,
      serverErrorMessages,

      updateField,
      updateModel,
      remove,
    };
  },
};
</script>
