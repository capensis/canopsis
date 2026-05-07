<template>
  <v-layout class="gap-3 py-3" column>
    <v-card v-for="(item, index) in items" :key="index">
      <v-card-text>
        <dynamic-info-information-item-form
          v-field="items[index]"
          :removable="items.length > 1"
          @remove="removeItem(index)"
        />
      </v-card-text>
    </v-card>
    <v-layout class="gap-2">
      <v-btn color="primary" outlined @click="addItem">
        {{ $t('modals.createDynamicInfo.infosSection.addInfos') }}
      </v-btn>
      <v-btn color="primary" outlined @click="showAddInfosFromTemplateModal">
        {{ $t('modals.createDynamicInfo.infosSection.addInfosFromTemplate') }}
      </v-btn>
    </v-layout>
  </v-layout>
</template>

<script>
import { MODALS } from '@/constants';

import { dynamicInfoInformationToForm } from '@/helpers/entities/dynamic-info/information/form';

import { useModals } from '@/hooks/modals';
import { useArrayModelField } from '@/hooks/form/array-model-field';

import DynamicInfoInformationItemForm from './dynamic-info-information-item-form.vue';

export default {
  inject: ['$validator'],
  components: {
    DynamicInfoInformationItemForm,
  },
  model: {
    prop: 'items',
    event: 'input',
  },
  props: {
    items: {
      type: Array,
      default: () => [],
    },
    existingNames: {
      type: Array,
      default: () => [],
    },
    initialName: {
      type: String,
      default: '',
    },
    variables: {
      type: Array,
      default: () => [],
    },
    copyVariables: {
      type: Array,
      default: () => [],
    },
  },
  setup(props, { emit }) {
    const modals = useModals();
    const { addItemIntoArray, removeItemFromArray } = useArrayModelField(props, emit);

    const addItem = () => addItemIntoArray(dynamicInfoInformationToForm());

    const showAddInfosFromTemplateModal = () => modals.show({
      name: MODALS.addDynamicInfoInfosFromTemplate,
      config: { action: (infosFromTemplate) => { infosFromTemplate.forEach(item => addItemIntoArray(item)); } },
    });

    return {
      addItem,
      removeItem: removeItemFromArray,
      showAddInfosFromTemplateModal,
    };
  },
};
</script>
