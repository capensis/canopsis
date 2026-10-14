<template>
  <v-layout
    class="gap-3"
    column
  >
    <c-entity-info-property-key-field
      :value="form.name"
      :label="$t('entityInfoProperties.infosKey')"
      :disabled="!isNew"
      name="name"
      required
      @update:selected-items="updateSelectedItems"
    />
    <c-form-block>
      <c-form-block-row :label="$t('common.description')">
        <c-description-field
          v-field="form.description"
          :max-length="255"
          name="description"
        />
      </c-form-block-row>
      <c-form-block-row :label="$t('common.alias')">
        <c-name-field
          v-field="form.alias"
          :max-length="255"
          name="alias"
        />
      </c-form-block-row>
      <c-form-block-row :label="$t('common.type')">
        <entity-info-property-type-field
          v-field="form.type"
          required
        />
      </c-form-block-row>
    </c-form-block>
  </v-layout>
</template>

<script>
import { ref } from 'vue';

import { useModelField } from '@/hooks/form';

import EntityInfoPropertyTypeField from './entity-info-property-type-field.vue';

export default {
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
