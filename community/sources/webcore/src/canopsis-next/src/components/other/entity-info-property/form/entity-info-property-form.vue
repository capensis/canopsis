<template>
  <v-layout column>
    <c-entity-info-property-key-field
      :value="form.name"
      :label="$t('entityInfoProperties.infosKey')"
      :disabled="!isNew"
      name="name"
      required
      @update:selected-items="updateSelectedItems"
    />
    <c-description-field
      v-field="form.description"
      :label="$t('common.description')"
      :max-length="255"
      name="description"
    />
    <c-name-field
      v-field="form.alias"
      :label="$t('common.alias')"
      :max-length="255"
      name="alias"
    />
    <entity-info-property-type-field
      v-field="form.type"
      required
    />
  </v-layout>
</template>

<script>
import { ref } from 'vue';

import { useModelField } from '@/hooks/form';

import EntityInfoPropertyTypeField from './entity-info-property-type-field.vue';

export default {
  inject: ['$validator'],
  components: {
    EntityInfoPropertyTypeField,
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
    isNew: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { updateModel } = useModelField(props, emit);

    const wasInit = ref(props.isNew);

    const updateSelectedItems = (selectedItems) => {
      if (!wasInit.value) {
        wasInit.value = true;

        return;
      }

      if (!selectedItems.length) {
        updateModel({ ...props.form, name: '' });

        return;
      }

      updateModel({
        ...props.form,

        name: selectedItems[0].value,
        type: selectedItems[0].type,
      });
    };

    return {
      updateSelectedItems,
    };
  },
};
</script>
