<template>
  <v-layout class="gap-2" column>
    <span class="text-subtitle-2">{{ $t('common.fields') }}</span>
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
                hide-details
                required
              />
              <c-enabled-field
                v-field="form[index].required"
                :label="$t('common.required')"
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
      <v-btn
        color="primary"
        fab
        small
        @click="addField"
      >
        <v-icon>add</v-icon>
      </v-btn>
    </v-layout>
  </v-layout>
</template>

<script>
import { uid } from '@/helpers/uid';

import { useArrayModelField } from '@/hooks/form/array-model-field';

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

    return {
      dragItemHandleClass,
      addField,
      removeField,
    };
  },
};
</script>
