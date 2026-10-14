<template>
  <c-layout gap="3" column>
    <c-entity-info-property-key-field
      v-model:value="name"
      :label="$t('entityInfoProperties.infosKey')"
      :disabled="!isNew"
      name="name"
      required
      @update:selected-items="updateSelectedItems"
    />

    <c-form-block>
      <c-form-block-row :label="$t('common.description')" indented>
        <c-description-field
          v-model="description"
          :max-length="FORM_MAX_STRING_LENGTH"
          name="description"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.alias')">
        <c-name-field v-model="alias" :max-length="FORM_MAX_STRING_LENGTH" name="alias" />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.type')">
        <entity-info-property-type-field v-model:value="type" required />
      </c-form-block-row>
    </c-form-block>
  </c-layout>
</template>

<script>
import { ref } from 'vue';

import { FORM_MAX_STRING_LENGTH } from '@/constants';

import { useFormModels } from '@/hooks/form';

import EntityInfoPropertyTypeField from './entity-info-property-type-field.vue';

export default {
  components: {
    EntityInfoPropertyTypeField,
  },

  props: {
    isNew: {
      type: Boolean,
      default: false,
    },
  },

  setup(props) {
    const { name, description, alias, type } = useFormModels([
      'name',
      'description',
      'alias',
      'type',
    ]);

    const wasInit = ref(props.isNew);

    const updateSelectedItems = (selectedItems) => {
      if (!wasInit.value) {
        wasInit.value = true;

        return;
      }

      if (!selectedItems.length) {
        name.value = '';

        return;
      }

      const [selectedItem] = selectedItems;

      name.value = selectedItem.value;

      if (selectedItem.type !== undefined) {
        type.value = selectedItem.type;
      }
    };

    return {
      name,
      description,
      alias,
      type,
      FORM_MAX_STRING_LENGTH,

      updateSelectedItems,
    };
  },
};
</script>
