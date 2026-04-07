<template>
  <v-layout class="gap-2" column>
    <span class="text-subtitle-2">{{ $t('common.fields') }}</span>
    <v-flex xs12>
      <v-alert
        :value="!form.length"
        :color="errors.has('fields') ? 'error' : 'info'"
      >
        {{ $t('commentTemplate.errors.fieldsRequired') }}
      </v-alert>
    </v-flex>
    <c-card-iterator-field
      v-field="form"
      :handle="`.${dragItemHandleClass}`"
      item-key="key"
    >
      <template #item="{ index }">
        <c-card-iterator-item
          :drag-handle-class="dragItemHandleClass"
          small
          @remove="removeField(index)"
        >
          <template #header>
            <v-layout column>
              <c-name-field
                v-field="form[index].name"
                :label="$t('common.fieldName')"
                :name="`fields[${index}].name`"
                :max-length="255"
                required
              />
              <c-enabled-field
                v-field="form[index].required"
                :label="$t('common.required')"
                class="mt-0"
                hide-details
              />
            </v-layout>
          </template>
        </c-card-iterator-item>
      </template>
    </c-card-iterator-field>
    <v-layout
      class="mt-2"
      justify-end
    >
      <v-tooltip left>
        <template #activator="{ on }">
          <v-btn
            color="primary"
            fab
            small
            v-on="on"
            @click="addField"
          >
            <v-icon>add</v-icon>
          </v-btn>
        </template>
        <span>{{ $t('common.add') }}</span>
      </v-tooltip>
    </v-layout>
  </v-layout>
</template>

<script>
import { watch, nextTick, onMounted, onBeforeUnmount } from 'vue';

import { uid } from '@/helpers/uid';

import { useArrayModelField } from '@/hooks/form/array-model-field';
import { useValidationAttachRequired } from '@/hooks/validator/validation-attach-required';

export default {
  inject: ['$validator'],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Array,
      required: true,
    },
  },
  setup(props, { emit }) {
    const dragItemHandleClass = 'field-drag-handler';

    const { addItemIntoArray, removeItemFromArray } = useArrayModelField(props, emit);
    const { attachRequiredRule, detachRequiredRule, validateRequiredRule } = useValidationAttachRequired('fields');

    const getter = () => !!props.form.length;
    const validate = () => nextTick(validateRequiredRule);

    /**
     * Add new field to the fields array
     */
    const addField = () => addItemIntoArray({
      key: uid(),
      name: '',
      required: false,
    });

    /**
     * Remove field from the fields array
     *
     * @param {number} index - Index of the field to remove
     */
    const removeField = index => removeItemFromArray(index);

    watch(() => props.form, () => validate());

    onMounted(() => {
      if (!props.form.length) {
        addField();
      }

      attachRequiredRule(getter);
    });

    onBeforeUnmount(detachRequiredRule);

    return {
      dragItemHandleClass,
      addField,
      validate,
      removeField,
    };
  },
};
</script>
