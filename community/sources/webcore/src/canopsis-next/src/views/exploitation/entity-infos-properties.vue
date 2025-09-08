<template>
  <c-page
    :create-tooltip="$t('modals.createEntityInfoProperty.create.title')"
    :creatable="hasCreateEntityInfoAccess"
    @create="showCreateEntityInfoPropertyModal"
    @refresh="refresh"
  >
    <entity-infos-properties-list
      :entity-infos-properties="entityInfosProperties"
      :pending="pending"
      :options.sync="options"
      :total-items="meta.total_count"
      :updatable="hasUpdateEntityInfoAccess"
      :removable="hasDeleteEntityInfoAccess"
      @edit="showEditEntityInfoPropertyModal"
      @duplicate="showDuplicateEntityInfoPropertyModal"
      @remove="showRemoveEntityInfoPropertyModal"
      @remove-selected="showRemoveSelectedEntityInfoPropertyModal"
    />
  </c-page>
</template>

<script>
import { onMounted } from 'vue';

import { MODALS, USER_PERMISSIONS } from '@/constants';

import { pickIds } from '@/helpers/array';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useObserver } from '@/hooks/observer';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';
import { useCRUDPermissions } from '@/hooks/auth';
import { useEntityInfoProperty } from '@/hooks/store/modules/entity-info-property';

import EntityInfosPropertiesList from '@/components/other/entity-info-property/entity-infos-properties-list.vue';

export default {
  components: { EntityInfosPropertiesList },
  setup() {
    const { t } = useI18n();
    const modals = useModals();
    const {
      hasCreateAccess: hasCreateEntityInfoAccess,
      hasUpdateAccess: hasUpdateEntityInfoAccess,
      hasDeleteAccess: hasDeleteEntityInfoAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.exploitation.entityInfoProperty);

    /**
     * STORE
     */
    const {
      createEntityInfoProperty,
      updateEntityInfoProperty,
      removeEntityInfoProperty,
      fetchEntityInfoPropertiesListWithoutStore,
      bulkRemoveEntityInfoProperty,
    } = useEntityInfoProperty();

    /**
     * QUERY
     */
    const {
      data: entityInfosProperties,
      meta,
      pending,
      options,
      updateOptions,
      fetchList,
    } = useFetchListWithoutStoreWithOptions({
      fetchListHandler: fetchEntityInfoPropertiesListWithoutStore,
    });

    /**
     * Displays a modal for creating a new entity info.
     */
    const showCreateEntityInfoPropertyModal = () => modals.show({
      name: MODALS.createEntityInfoProperty,
      config: {
        action: async (data) => {
          await createEntityInfoProperty({ data });

          return fetchList();
        },
      },
    });

    /**
     * Displays a modal for editing an existing entity info.
     *
     * @param {Object} entityInfoProperty - The entity info to be edited.
     */
    const showEditEntityInfoPropertyModal = entityInfoProperty => modals.show({
      name: MODALS.createEntityInfoProperty,
      config: {
        entityInfoProperty,
        title: t('modals.createEntityInfoProperty.edit.title'),
        action: async (data) => {
          await updateEntityInfoProperty({ id: entityInfoProperty._id, data });

          return fetchList();
        },
      },
    });

    /**
     * Displays a modal for duplicating an entity info.
     *
     * @param {Object} entityInfoProperty - The entity info to be duplicated.
     */
    const showDuplicateEntityInfoPropertyModal = entityInfoProperty => modals.show({
      name: MODALS.createEntityInfoProperty,
      config: {
        entityInfoProperty: { ...entityInfoProperty, _id: undefined },
        title: t('modals.createEntityInfoProperty.duplicate.title'),
        action: async (data) => {
          await createEntityInfoProperty({ data });

          return fetchList();
        },
      },
    });

    /**
     * Displays a confirmation modal for removing an entity info.
     *
     * @param {Object} entityInfoProperty - The entity info to be removed.
     */
    const showRemoveEntityInfoPropertyModal = entityInfoProperty => modals.show({
      name: MODALS.confirmation,
      config: {
        title: t('modals.confirmationRemoveEntityInfoProperty.title'),
        alert: t('modals.confirmationRemoveEntityInfoProperty.alert', { name: entityInfoProperty.name }),
        text: t('modals.confirmationRemoveEntityInfoProperty.text'),
        action: async () => {
          await removeEntityInfoProperty({ id: entityInfoProperty._id });

          return fetchList();
        },
      },
    });

    /**
     * Displays a confirmation modal for removing selected entity infos.
     *
     * @param {Object[]} selected - The entity infos to be removed.
     */
    const showRemoveSelectedEntityInfoPropertyModal = selected => modals.show({
      name: MODALS.confirmation,
      config: {
        text: t('modals.confirmationRemoveEntityInfoProperty.text', { name: selected.length }),
        action: async () => {
          await bulkRemoveEntityInfoProperty({ data: pickIds(selected) });

          return fetchList();
        },
      },
    });

    const { observer } = useObserver({ key: '$refresh' });

    const refresh = () => observer.notify();

    onMounted(() => {
      observer.register(fetchList);
      fetchList();
    });

    return {
      hasCreateEntityInfoAccess,
      hasUpdateEntityInfoAccess,
      hasDeleteEntityInfoAccess,

      entityInfosProperties,
      meta,
      pending,
      options,
      updateOptions,

      refresh,
      showCreateEntityInfoPropertyModal,
      showEditEntityInfoPropertyModal,
      showDuplicateEntityInfoPropertyModal,
      showRemoveEntityInfoPropertyModal,
      showRemoveSelectedEntityInfoPropertyModal,
    };
  },
};
</script>
