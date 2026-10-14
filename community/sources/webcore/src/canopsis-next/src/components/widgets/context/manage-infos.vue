<template>
  <div>
    <c-layout justify="end">
      <v-btn color="primary" variant="outlined" @click="showAddInfoModal">
        {{ t('entity.addInfos') }}
      </v-btn>
    </c-layout>
    <v-data-table
      v-model:options="options"
      :items="infos"
      :headers="tableHeaders"
      :no-data-text="t('entity.emptyInfos')"
      item-value="name"
    >
      <template #item="{ item, index }">
        <tr>
          <td>{{ getInfoItem(item).name }}</td>
          <td>{{ getInfoItem(item).description }}</td>
          <td>{{ getInfoItem(item).value }}</td>
          <td>
            <c-layout>
              <c-action-btn
                type="edit"
                @click="showEditInfoModal(index, getInfoItem(item))"
              />
              <c-action-btn
                type="delete"
                @click="removeItemFromArray(index)"
              />
            </c-layout>
          </td>
        </tr>
      </template>
    </v-data-table>
  </div>
</template>

<script>
import { computed, ref } from 'vue';

import { MODALS } from '@/constants';

import { useArrayModelField } from '@/hooks/form/array-model-field';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

export default {
  props: {
    infos: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['update:infos'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const modals = useModals();
    const {
      addItemIntoArray,
      updateItemInArray,
      removeItemFromArray,
    } = useArrayModelField(props, emit);

    const options = ref({
      page: 1,
      itemsPerPage: 10,
    });

    const tableHeaders = computed(() => [
      { title: t('common.name'), key: 'name' },
      { title: t('common.description'), key: 'description' },
      { title: t('common.value'), key: 'value' },
      { title: t('common.actionsLabel'), key: 'actions', sortable: false },
    ]);

    const getInfoItem = item => item?.raw ?? item;

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
      t,
      options,
      tableHeaders,
      getInfoItem,
      showAddInfoModal,
      showEditInfoModal,
      removeItemFromArray,
    };
  },
};
</script>
