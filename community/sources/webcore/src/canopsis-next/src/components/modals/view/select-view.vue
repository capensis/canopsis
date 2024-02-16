<template>
  <modal-wrapper close>
    <template #title="">
      <span>{{ $t('modals.selectView.title') }}</span>
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
          <v-expansion-panels
            accordion
            dark
          >
            <v-expansion-panel
              class="secondary"
              v-for="group in groups"
              :key="group._id"
            >
              <v-expansion-panel-header>
                {{ group.title }}
              </v-expansion-panel-header>
              <v-expansion-panel-content ripple>
                <v-list class="py-0 px-2 secondary">
                  <v-list-item
                    class="secondary lighten-1"
                    v-for="view in group.views"
                    :key="view._id"
                    ripple
                    @click="selectView(view._id)"
                  >
                    <v-list-item-title class="text-body-1">
                      {{ view.title }}
                    </v-list-item-title>
                  </v-list-item>
                </v-list>
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

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.selectView,
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
    async selectView(viewId) {
      if (this.config.action) {
        await this.config.action(viewId);
      }

      this.$modals.hide();
    },
  },
};
</script>
