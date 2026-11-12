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
      item-key="name"
    >
      <template #item="{ item }">
        <tr>
          <td>{{ item.name }}</td>
          <td>{{ item.description }}</td>
          <td>{{ item.value }}</td>
          <td>
            <v-layout>
              <c-action-btn
                type="edit"
                @click="showEditInfoModal(item)"
              />
              <c-action-btn
                type="delete"
                @click="removeInfo(item)"
              />
            </v-layout>
          </td>
        </tr>
      </template>
    </v-data-table>
  </div>
</template>

<script>
import { computed } from 'vue';

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
     * Opens the edit entity info modal for the given item. On submit, finds and updates the info by name.
     *
     * @param {Object} info - The entity info object to edit (name, description, value)
     */
    const showEditInfoModal = info => modals.show({
      name: MODALS.createEntityInfo,
      config: {
        infos: props.infos,
        entityInfo: info,
        title: t('modals.createEntityInfo.edit.title'),
        action: (editedInfo) => {
          const index = props.infos.findIndex(({ name }) => name === info.name);

          if (index !== -1) {
            updateItemInArray(index, editedInfo);
          }
        },
      },
    });

    /**
     * Finds and removes the given info by name.
     *
     * @param {Object} info - The entity info object to remove
     */
    const removeInfo = (info) => {
      const index = props.infos.findIndex(({ name }) => name === info.name);

      if (index !== -1) {
        removeItemFromArray(index);
      }
    };

    return {
      tableHeaders,
      showAddInfoModal,
      showEditInfoModal,
      removeInfo,
    };
  },
};
</script>
