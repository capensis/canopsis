<template>
  <v-layout column>
    <c-name-field
      v-field="form.name"
      :label="$t('common.name')"
      name="name"
      required
    />
    <v-layout
      class="mt-3"
      column
    >
      <span class="mb-2">{{ $t('common.fields') }}</span>
      <c-draggable-list-field
        v-field="form.fields"
        handle=".field-drag-handler"
      >
        <v-layout
          v-for="(field, index) in form.fields"
          :key="field.key"
          class="mb-2"
          align-center
        >
          <v-icon
            class="field-drag-handler draggable mr-2"
            size="20"
          >
            drag_indicator
          </v-icon>
          <v-flex>
            <v-card flat>
              <v-card-text>
                <v-layout
                  align-center
                  justify-space-between
                >
                  <v-flex xs10>
                    <c-name-field
                      v-field="form.fields[index].name"
                      :label="$t('common.fieldName')"
                      :name="`fields[${index}].name`"
                      hide-details
                      required
                    />
                  </v-flex>
                  <v-flex
                    class="text-right"
                    xs2
                  >
                    <c-action-btn
                      type="delete"
                      @click="removeField(index)"
                    />
                  </v-flex>
                </v-layout>
                <v-layout
                  class="mt-2"
                  align-center
                >
                  <c-enabled-field
                    v-field="form.fields[index].required"
                    :label="$t('common.required')"
                    hide-details
                  />
                </v-layout>
              </v-card-text>
            </v-card>
          </v-flex>
        </v-layout>
      </c-draggable-list-field>
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
  </v-layout>
</template>

<script>
import { computed } from 'vue';

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
      type: Object,
      required: true,
    },
  },
  setup(props, { emit }) {
    const fields = computed({
      get: () => props.form.fields,
      set: value => emit('input', { ...props.form, fields: value }),
    });

    const { addItemIntoArray, removeItemFromArray } = useArrayModelField({ form: fields }, emit);

    /**
     * Add new field to the fields array
     */
    const addField = () => {
      addItemIntoArray({
        key: uid(),
        name: '',
        required: false,
      });
    };

    /**
     * Remove field from the fields array
     *
     * @param {number} index - Index of the field to remove
     */
    const removeField = (index) => {
      removeItemFromArray(index);
    };

    return {
      addField,
      removeField,
    };
  },
};
</script>
