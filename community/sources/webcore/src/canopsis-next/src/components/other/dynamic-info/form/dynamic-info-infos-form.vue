<template>
  <v-layout class="gap-3 py-3" column>
    <v-card v-for="(item, index) in items" :key="item.key">
      <v-card-text>
        <dynamic-info-infos-item-form
          v-field="items[index]"
          :name="`items.${index}`"
          :variables="variables"
          :copy-variables="copyVariables"
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

import DynamicInfoInfosItemForm from './fields/dynamic-info-infos-item-form.vue';

export default {
  inject: ['$validator'],
  components: {
    DynamicInfoInfosItemForm,
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
    const { updateModel, addItemIntoArray, removeItemFromArray } = useArrayModelField(props, emit);

    /**
     * Appends a new dynamic info row with default form fields to the list.
     */
    const addItem = () => addItemIntoArray(dynamicInfoInformationToForm());

    /**
     * Opens the modal to pick a dynamic info template; on confirm, adds one row per template name preset.
     */
    const showAddInfosFromTemplateModal = () => modals.show({
      name: MODALS.addDynamicInfoInfosFromTemplate,
      config: {
        action: template => (
          updateModel([...props.items, ...template.names.map(name => dynamicInfoInformationToForm({ name }))])
        ),
      },
    });

    return {
      addItem,
      removeItem: removeItemFromArray,
      showAddInfosFromTemplateModal,
    };
  },
};
</script>
