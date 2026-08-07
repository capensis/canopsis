<template>
  <c-modal-wrapper close>
    <template #title>
      <span>{{ $t('modals.selectView.title') }}</span>
    </template>

    <template #text>
      <v-fade-transition>
        <c-layout v-if="pending" justify="center">
          <v-progress-circular color="primary" indeterminate />
        </c-layout>

        <c-layout v-else>
          <v-expansion-panels>
            <v-expansion-panel
              v-for="group in publicGroups"
              :key="group._id"
              class="secondary"
            >
              <v-expansion-panel-title>
                {{ group.title }}
              </v-expansion-panel-title>

              <v-expansion-panel-text>
                <v-list class="py-0 px-2 secondary">
                  <v-list-item
                    v-for="view in group.views"
                    :key="view._id"
                    class="secondary lighten-1"
                    ripple
                    @click="selectView(view._id)"
                  >
                    <v-list-item-title class="text-body-1">
                      {{ view.title }}
                    </v-list-item-title>
                  </v-list-item>
                </v-list>
              </v-expansion-panel-text>
            </v-expansion-panel>
          </v-expansion-panels>
        </c-layout>
      </v-fade-transition>
    </template>
  </c-modal-wrapper>
</template>

<script>
import { mapActions, mapState } from 'pinia';

import { MAX_LIMIT } from '@/constants';
import { useViewStore } from '@/stores/modules/view';

export default {
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      pending: true,
    };
  },
  computed: {
    ...mapState(useViewStore, {
      groups: 'items',
    }),

    config() {
      return this.modal.config ?? {};
    },

    publicGroups() {
      return this.groups.filter(group => !group.is_private);
    },
  },
  async mounted() {
    this.pending = true;

    await this.fetchViewGroupsList({
      params: {
        limit: MAX_LIMIT,
        page: 1,
        with_views: true,
      },
    });

    this.pending = false;
  },
  methods: {
    ...mapActions(useViewStore, {
      fetchViewGroupsList: 'fetchList',
    }),

    async selectView(viewId) {
      if (this.config.action) {
        await this.config.action(viewId);
      }

      this.$modals.hide({ id: this.modal.id });
    },
  },
};
</script>
