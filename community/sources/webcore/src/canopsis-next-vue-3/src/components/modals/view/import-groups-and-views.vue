<template>
  <c-modal-wrapper close>
    <template #title>
      <span>{{ t('modals.importExportViews.title') }}</span>
    </template>

    <template #text>
      <c-layout class="import-groups-and-views__content" align="start">
        <div class="import-groups-and-views__column px-1">
          <div class="text-center mb-2">
            {{ t('modals.importExportViews.groups') }}
          </div>
          <draggable-groups
            :groups="importedGroups"
            pull
            view-pull
            view-put
            @change="importedGroups = $event"
            @change:group="changeImportedGroupHandler"
          />
        </div>

        <div class="import-groups-and-views__column px-1">
          <div class="text-center mb-2">
            {{ t('modals.importExportViews.views') }}
          </div>
          <draggable-group-views :views="importedViews" pull @change="importedViews = $event" />
        </div>

        <v-divider class="mx-1 secondary align-self-stretch" vertical />

        <div class="import-groups-and-views__column">
          <div class="text-center mb-2">
            {{ t('common.result') }}
          </div>
          <draggable-groups
            :groups="currentGroups"
            put
            pull
            view-put
            @change="currentGroups = $event"
            @change:group="changeCurrentGroupHandler"
          />
        </div>
      </c-layout>
    </template>

    <template #actions>
      <v-btn variant="text" @click="close">
        {{ t('common.cancel') }}
      </v-btn>

      <v-btn
        :loading="submitting"
        :disabled="submitting"
        color="primary"
        @click="submit"
      >
        {{ t('common.saveChanges') }}
      </v-btn>
    </template>
  </c-modal-wrapper>
</template>

<script>
import { computed, ref, watch } from 'vue';
import { useRoute } from 'vue-router';

import { MAX_LIMIT, ROUTES_NAMES } from '@/constants';

import {
  prepareCurrentGroupsForImporting,
  prepareImportedGroups,
  prepareImportedViews,
  prepareViewGroupsForImportRequest,
} from '@/helpers/entities/view/form';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useView } from '@/hooks/store/modules/view';

import DraggableGroupViews from '@/components/layout/navigation/partials/groups-side-bar/draggable-group-views.vue';
import DraggableGroups from '@/components/layout/navigation/partials/groups-side-bar/draggable-groups.vue';

export default {
  components: {
    DraggableGroupViews,
    DraggableGroups,
  },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const route = useRoute();
    const modals = useModals();
    const popups = usePopups();
    const { viewGroups, fetchViewGroupsList, importViewWithoutStore } = useView();

    const importedViews = ref([]);
    const importedGroups = ref([]);
    const currentGroups = ref([]);
    const submitting = ref(false);

    const config = computed(() => props.modal.config ?? {});
    const close = () => modals.hide({ id: props.modal.id });
    const isViewsImportExportPage = computed(
      () => route.name === ROUTES_NAMES.adminSettingsViewsImportExport,
    );

    const setDefaultValues = () => {
      const { importedGroups: nextImportedGroups = [], importedViews: nextImportedViews = [] } =
        config.value;

      importedViews.value = prepareImportedViews(nextImportedViews);
      importedGroups.value = prepareImportedGroups(nextImportedGroups);
      currentGroups.value = prepareCurrentGroupsForImporting(
        viewGroups.value.filter(group => !group.is_private),
      );
    };

    const fetchGroups = () => {
      if (!isViewsImportExportPage.value) {
        return fetchViewGroupsList();
      }

      return fetchViewGroupsList({
        params: {
          limit: MAX_LIMIT,
          page: 1,
          with_views: true,
        },
      });
    };

    const changeImportedGroupHandler = (groupIndex, group) =>
      importedGroups.value.splice(groupIndex, 1, group);

    const changeCurrentGroupHandler = (groupIndex, group) =>
      currentGroups.value.splice(groupIndex, 1, group);

    const submit = async () => {
      try {
        submitting.value = true;

        await importViewWithoutStore({
          data: prepareViewGroupsForImportRequest(currentGroups.value),
        });
      } catch (err) {
        console.error(err);

        popups.error({ text: err.description || t('errors.default') });
      } finally {
        await fetchGroups();
        submitting.value = false;
        close();
      }
    };

    watch(viewGroups, setDefaultValues, { immediate: true });

    return {
      t,
      importedViews,
      importedGroups,
      currentGroups,
      submitting,
      close,
      submit,
      changeImportedGroupHandler,
      changeCurrentGroupHandler,
    };
  },
};
</script>

<style lang="scss" scoped>
.import-groups-and-views__content {
  flex-wrap: nowrap;
}

.import-groups-and-views__column {
  flex: 1 1 0;
  min-width: 0;
}
</style>
