<template>
  <c-modal-wrapper close>
    <template #title>
      <span>{{ $t('modals.selectViewTab.title') }}</span>
    </template>

    <template #text>
      <div class="select-view-tab-modal__content">
        <c-progress-overlay :pending="pending" />

        <c-layout>
          <v-expansion-panels>
            <v-expansion-panel v-for="group in groups" :key="group._id" class="secondary">
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
        </c-layout>
      </div>
    </template>
  </c-modal-wrapper>
</template>

<script>
import { mapActions } from 'pinia';

import { MAX_LIMIT, MODALS } from '@/constants';

import { useAuthStore } from '@/stores/modules/auth';
import { useViewStore } from '@/stores/modules/view';

export default {
  name: MODALS.selectViewTab,

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  data() {
    return {
      groups: [],
      pending: true,
      submitting: false,
    };
  },

  computed: {
    config() {
      return this.modal.config ?? {};
    },
  },

  async mounted() {
    this.pending = true;

    const [groups] = await Promise.all([
      this.fetchViewGroupsListWithoutStore({
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
      this.fetchCurrentUser(),
    ]);

    this.groups = groups?.data ?? groups ?? [];
    this.pending = false;
  },

  methods: {
    ...mapActions(useViewStore, {
      fetchViewGroupsListWithoutStore: 'fetchListWithoutStore',
    }),
    ...mapActions(useAuthStore, ['fetchCurrentUser']),

    async selectTab(tabId, viewId) {
      try {
        this.submitting = true;

        if (this.config.action) {
          await this.config.action({ tabId, viewId });
        }

        this.$modals.hide({ id: this.modal.id });
      } catch (err) {
        console.warn(err);
      } finally {
        this.submitting = false;
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.select-view-tab-modal__content {
  position: relative;
  min-height: 160px;
}
</style>
