<template>
  <c-modal-wrapper close>
    <template #title>
      <span>{{ title }}</span>
    </template>

    <template #text>
      <div class="select-view-tab-modal__content">
        <c-progress-overlay :pending="pending" />

        <v-row no-gutters>
          <v-col cols="12">
            <v-expansion-panels>
              <v-expansion-panel v-for="group in visibleGroups" :key="group._id" class="secondary">
                <v-expansion-panel-title>
                  {{ group.title }}
                </v-expansion-panel-title>

                <v-expansion-panel-text>
                  <v-expansion-panels>
                    <v-expansion-panel
                      v-for="view in group.views"
                      :key="view._id"
                      class="secondary lighten-1"
                    >
                      <v-expansion-panel-title>
                        {{ view.title }}
                      </v-expansion-panel-title>

                      <v-expansion-panel-text>
                        <v-list class="pa-0">
                          <v-list-item
                            v-for="tab in view.tabs"
                            :key="tab._id"
                            :disabled="pending || submitting"
                            class="secondary lighten-2"
                            @click="selectTab(tab._id, view._id)"
                          >
                            <v-list-item-title class="text-body-1 pl-4">
                              {{ tab.title }}
                            </v-list-item-title>
                          </v-list-item>
                        </v-list>
                      </v-expansion-panel-text>
                    </v-expansion-panel>
                  </v-expansion-panels>
                </v-expansion-panel-text>
              </v-expansion-panel>
            </v-expansion-panels>
          </v-col>
        </v-row>
      </div>
    </template>
  </c-modal-wrapper>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue';

import { MAX_LIMIT } from '@/constants';

import { useAuth } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { usePendingHandler } from '@/hooks/query/pending';
import { useView } from '@/hooks/store/modules/view';

const props = defineProps({
  modal: {
    type: Object,
    required: true,
  },
});

const { t } = useI18n();
const { config, close } = useInnerModal(props);
const { fetchCurrentUser } = useAuth();
const { fetchViewGroupsListWithoutStore } = useView();

const groups = ref([]);

const title = computed(() => config.value.title ?? t('modals.selectViewTab.title'));
const visibleGroups = computed(() => (
  config.value.publicOnly
    ? groups.value.filter(group => !group.is_private)
    : groups.value
));

const { pending: submitting, handler: selectTab } = usePendingHandler(async (tabId, viewId) => {
  try {
    const view = groups.value
      .flatMap(group => group.views ?? [])
      .find(item => item._id === viewId);
    const tab = view?.tabs?.find(item => item._id === tabId);

    await config.value.action?.({ tabId, viewId, tab, view });

    close();
  } catch (err) {
    console.warn(err);
  }
});

const { pending, handler: fetchList } = usePendingHandler(async () => {
  const [groupsResponse] = await Promise.all([
    fetchViewGroupsListWithoutStore({
      params: {
        limit: MAX_LIMIT,
        page: 1,
        with_views: true,
        with_tabs: true,
        with_widgets: true,
        with_flags: true,
        with_private: true,
      },
    }),
    fetchCurrentUser(),
  ]);

  groups.value = groupsResponse?.data ?? groupsResponse ?? [];
}, true);

onMounted(fetchList);
</script>

<style lang="scss" scoped>
.select-view-tab-modal__content {
  position: relative;
  min-height: 160px;
}
</style>
