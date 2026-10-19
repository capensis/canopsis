<template>
  <v-card>
    <v-card-text>
      <v-layout class="gap-2" column>
        <v-layout align-center>
          <v-text-field
            v-field="form.reference"
            v-validate="'required'"
            :label="$t('externalData.fields.reference')"
            :error-messages="errors.collect(referenceFieldName)"
            :name="referenceFieldName"
            :disabled="disabled"
            class="mr-2"
          >
            <template #append="">
              <c-help-icon
                :text="$t('externalData.tooltips.reference')"
                icon="help"
                left
              />
            </template>
          </v-text-field>
          <v-select
            v-field="form.type"
            :items="availableTypes"
            :label="$t('common.type')"
            :disabled="disabled"
            class="ml-2"
          />
          <v-btn
            v-if="!disabled"
            class="mr-0"
            icon
            @click="remove"
          >
            <v-icon color="error">
              delete
            </v-icon>
          </v-btn>
        </v-layout>
        <external-data-table-form
          v-if="isTableType"
          v-field="form"
          :name="name"
          :disabled="disabled"
          :variables="variables"
          :optionally="optionally"
        />
        <request-form
          v-else
          v-field="form.request"
          :name="`${name}.request`"
          :disabled="disabled"
          :payload-variables="variables"
          :url-variables="variables"
        />
        <c-alert v-if="serverErrorMessages.length" type="error">
          {{ serverErrorMessages.join('\n') }}
        </c-alert>
      </v-layout>
    </v-card-text>
  </v-card>
</template>

<script>
import { computed, onMounted } from 'vue';

import { EXTERNAL_DATA_TYPES } from '@/constants';

import { isTableExternalDataType } from '@/helpers/entities/shared/external-data/entity';

import { useI18n } from '@/hooks/i18n';
import { useValidator } from '@/hooks/validator/validator';

import RequestForm from '@/components/forms/request/request-form.vue';

import ExternalDataTableForm from './external-data-table-form.vue';

export default {
  inject: ['$validator'],
  components: { RequestForm, ExternalDataTableForm },
  model: {
    prop: 'form',
    event: 'input',
  },
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
      default: () => ([]),
    },
    optionally: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();
    const validator = useValidator();

    const availableTypes = computed(() => (
      props.types.length
        ? props.types
        : Object.values(EXTERNAL_DATA_TYPES)
          .map(type => ({ text: t(`externalData.types.${type}`), value: type }))
    ));

    const isTableType = computed(() => isTableExternalDataType(props.form.type));
    const referenceFieldName = computed(() => `${props.name}.reference`);
    const serverErrorMessages = computed(() => validator.errors.collect(props.serverErrorName));

    const remove = () => emit('remove', props.form);

    onMounted(() => {
      if (props.serverErrorName) {
        validator.attach({ name: props.serverErrorName });
      }
    });

    return {
      availableTypes,
      isTableType,
      referenceFieldName,
      serverErrorMessages,

      remove,
    };
  },
};
</script>
