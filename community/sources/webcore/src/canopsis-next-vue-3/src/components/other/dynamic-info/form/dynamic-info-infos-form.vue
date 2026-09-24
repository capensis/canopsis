<template>
  <c-layout class="py-3" gap="3" column>
    <v-card v-for="(field, index) in infosField" :key="field.key">
      <v-card-text>
        <dynamic-info-infos-item-form
          :form="field"
          :name="`${name}.${field.key}`"
          :unique-names="uniqueNames(index)"
          :variables="variables"
          :copy-variables="copyVariables"
          :removable="infosField.length > 1"
          @update:form="updateItem(index, $event)"
          @remove="remove(index)"
        />
      </v-card-text>
    </v-card>

    <c-layout gap="2">
      <v-btn color="primary" variant="outlined" @click="addItem">
        {{ $t('modals.createDynamicInfo.infosSection.addInfos') }}
      </v-btn>

      <v-btn color="primary" variant="outlined" @click="showAddInfosFromTemplateModal">
        {{ $t('modals.createDynamicInfo.infosSection.addInfosFromTemplate') }}
      </v-btn>
    </c-layout>
  </c-layout>
</template>

<script>
import { MODALS } from '@/constants';

import { dynamicInfoInformationToForm } from '@/helpers/entities/dynamic-info/information/form';

import { useFormContext } from '@/hooks/form';
import { useModals } from '@/hooks/modals';

import DynamicInfoInfosItemForm from './fields/dynamic-info-infos-item-form.vue';

export default {
  components: {
    DynamicInfoInfosItemForm,
  },

  props: {
    variables: {
      type: Array,
      default: () => [],
    },
    copyVariables: {
      type: Array,
      default: () => [],
    },
    name: {
      type: String,
      default: 'infos',
    },
  },

  setup() {
    const modals = useModals();
    const form = useFormContext();
    const infosField = form.array('infos');

    const uniqueNames = index => infosField.value
      .filter((_, itemIndex) => itemIndex !== index)
      .map(({ name }) => name)
      .filter(Boolean);

    const addItem = () => infosField.value.push(dynamicInfoInformationToForm());
    const remove = index => form.remove('infos', index);
    const updateItem = (index, value) => form.replace('infos', index, value);
    const showAddInfosFromTemplateModal = () => modals.show({
        name: MODALS.addDynamicInfoInfosFromTemplate,
        config: {
          action: template => template.names.forEach(
            name => infosField.value.push(dynamicInfoInformationToForm({ name })),
          ),
        },
      });

    return {
      infosField,
      uniqueNames,
      addItem,
      showAddInfosFromTemplateModal,
      updateItem,
      remove,
    };
  },
};
</script>
