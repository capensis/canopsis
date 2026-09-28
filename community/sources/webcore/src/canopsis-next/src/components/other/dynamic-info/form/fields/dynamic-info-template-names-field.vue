<template>
  <v-layout class="gap-2" column>
    <c-label
      :label="$t('modals.createDynamicInfoTemplate.fields.names')"
      required
    />
    <v-layout
      v-for="(name, index) in names"
      :key="name.key"
      align-center
      justify-space-between
    >
      <v-flex :xs11="names.length > 1" :xs12="names.length <= 1">
        <c-name-field
          v-field="names[index].value"
          :label="$t('common.name')"
          :name="`names[${name.key}]`"
          required
        />
      </v-flex>
      <v-flex
        v-if="names.length > 1"
        shrink
      >
        <v-btn
          color="error"
          icon
          @click="removeName(index)"
        >
          <v-icon>delete</v-icon>
        </v-btn>
      </v-flex>
    </v-layout>
    <div>
      <v-btn
        color="primary"
        outlined
        @click="addName"
      >
        {{ $t('modals.createDynamicInfoTemplate.buttons.addName') }}
      </v-btn>
    </div>
  </v-layout>
</template>

<script>
import { generateTemplateFormName } from '@/helpers/entities/dynamic-info/template/form';

import { useArrayModelField } from '@/hooks/form/array-model-field';

export default {
  model: {
    prop: 'names',
    event: 'input',
  },
  props: {
    names: {
      type: Array,
      default: () => [],
    },
  },
  setup(props, { emit }) {
    const { addItemIntoArray, removeItemFromArray } = useArrayModelField(props, emit);

    const addName = () => addItemIntoArray(generateTemplateFormName());

    return {
      addName,
      removeName: removeItemFromArray,
    };
  },
};
</script>
