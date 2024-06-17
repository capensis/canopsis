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
          <v-expansion-panels dark>
            <v-expansion-panel
              v-for="group in groups"
              :key="group._id"
              class="mt-0 secondary"
            >
              <v-expansion-panel-header>{{ group.title }}</v-expansion-panel-header>
              <v-expansion-panel-content ripple>
                <v-expansion-panels dark>
                  <v-expansion-panel
                    v-for="view in group.views"
                    :key="view._id"
                    class="mt-0 px-2 secondary lighten-1"
                    dark
                  >
                    <v-expansion-panel-header>{{ view.title }}</v-expansion-panel-header>
                    <v-expansion-panel-content ripple>
                      <v-list class="pa-0">
                        <v-list-item
                          v-for="tab in view.tabs"
                          :key="tab._id"
                          class="secondary lighten-2"
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
                </v-expansion-panels>
              </v-expansion-panel-content>
            </v-expansion-panel>
          </v-expansion-panels>
        </v-layout>
      </v-fade-transition>
    </template>
  </modal-wrapper>
</template>

<script>
import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';
import { submittableMixinCreator } from '@/mixins/submittable';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.selectViewTab,
  components: { ModalWrapper },
  mixins: [
    modalInnerMixin,
    entitiesViewGroupMixin,
    submittableMixinCreator({

    }),
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
