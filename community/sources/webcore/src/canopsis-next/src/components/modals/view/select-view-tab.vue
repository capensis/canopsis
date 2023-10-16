<template>
  <modal-wrapper close>
    <template #title="">
      <span>{{ $t('modals.selectViewTab.title') }}</span>
    </template>
    <template #text="">
      <v-fade-transition>
        <v-layout
          v-if="pending"
          justify-center
        >
          <v-progress-circular
            color="primary"
            indeterminate
          />
        </v-layout>
        <v-layout v-else>
          <v-expansion-panel dark>
            <v-expansion-panel-content
              class="secondary"
              v-for="group in groups"
              :key="group._id"
              ripple
            >
              <template #header="">
                <div>{{ group.title }}</div>
              </template>
              <v-expansion-panel
                class="px-2"
                dark
              >
                <v-expansion-panel-content
                  class="secondary lighten-1"
                  v-for="view in group.views"
                  :key="view._id"
                  ripple
                >
                  <template #header="">
                    <div>{{ view.title }}</div>
                  </template>
                  <v-list class="pa-0">
                    <v-list-item
                      class="secondary lighten-2"
                      v-for="tab in view.tabs"
                      :key="tab._id"
                      ripple
                      @click="selectTab(tab._id, view._id)"
                    >
                      <v-list-item-title class="text-body-1 pl-4">
                        {{ tab.title }}
                      </v-list-item-title>
                    </v-list-item>
                  </v-list>
                </v-expansion-panel-content>
              </v-expansion-panel>
            </v-expansion-panel-content>
          </v-expansion-panel>
        </v-layout>
      </v-fade-transition>
    </template>
  </modal-wrapper>
</template>

<script>
import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.selectViewTab,
  components: { ModalWrapper },
  mixins: [
    modalInnerMixin,
    entitiesViewGroupMixin,
  ],
  data() {
    return {
      pending: true,
    };
  },
  async mounted() {
    this.pending = true;

    await this.fetchAllGroupsListWithWidgetsWithCurrentUser();

    this.pending = false;
  },
  methods: {
    async selectTab(tabId, viewId) {
      if (this.config.action) {
        await this.config.action({ tabId, viewId });
      }

      this.$modals.hide();
    },
  },
};
</script>
