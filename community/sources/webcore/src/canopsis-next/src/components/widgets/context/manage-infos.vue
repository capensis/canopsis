<template>
  <div>
    <v-layout justify-end>
      <c-action-btn
        :tooltip="$t('entity.addInformation')"
        icon="add"
        @click="showAddInfoModal"
      />
    </v-layout>
    <v-data-table
      :items="infos"
      :headers="tableHeaders"
      :no-data-text="$t('entity.emptyInfos')"
      :options.sync="options"
      item-key="name"
    >
      <template #item="{ item, index }">
        <tr>
          <td>{{ item.name }}</td>
          <td>{{ item.description }}</td>
          <td>{{ item.value }}</td>
          <td>
            <v-layout>
              <c-action-btn
                type="edit"
                @click="showEditInfoModal(index, item)"
              />
              <c-action-btn
                type="delete"
                @click="removeItemFromArray(index)"
              />
            </v-layout>
          </td>
        </tr>
      </template>
    </v-data-table>
  </div>
</template>

<script>
import { ref, computed } from 'vue';

import { MODALS } from '@/constants';

import { useArrayModelField } from '@/hooks/form/array-model-field';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

export default {
  model: {
    prop: 'infos',
    event: 'input',
  },
  props: {
    infos: {
      type: Array,
      default: () => [],
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();
    const modals = useModals();
    const {
      addItemIntoArray,
      updateItemInArray,
      removeItemFromArray,
    } = useArrayModelField(props, emit);

    const options = ref({});

    const tableHeaders = computed(() => [
      { text: t('common.name'), value: 'name' },
      { text: t('common.description'), value: 'description' },
      { text: t('common.value'), value: 'value' },
      { text: t('common.actionsLabel'), value: 'actions', sortable: false },
    ]);

    /**
     * Opens the create entity info modal. On submit, adds the new info to the infos array.
     */
    const showAddInfoModal = () => modals.show({
      name: MODALS.createEntityInfo,
      config: {
        infos: props.infos,
        action: info => addItemIntoArray(info),
      },
    });

    /**
     * Opens the edit entity info modal for the given item. On submit, updates the info at the given index.
     *
     * @param {number} index - Index of the info item in the infos array
     * @param {Object} info - The entity info object to edit (name, description, value)
     */
    const showEditInfoModal = (index, info) => modals.show({
      name: MODALS.createEntityInfo,
      config: {
        infos: props.infos,
        entityInfo: info,
        title: t('modals.createEntityInfo.edit.title'),
        action: (editedInfo) => {
          const realIndex = (options.value.page - 1) * options.value.itemsPerPage + index;

          updateItemInArray(realIndex, editedInfo);
        },
      },
    });

    return {
      options,
      tableHeaders,
      showAddInfoModal,
      showEditInfoModal,
      removeItemFromArray,
    };
  },
};
</script>
