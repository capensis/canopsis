<template lang="pug">
  div
    c-page-header {{ $t('pattern.patterns') }}
    v-layout(row, wrap)
      v-flex(xs12)
        v-card.ma-2
          v-tabs(
            v-if="hasReadAnyCorporatePatternAccess",
            v-model="activeTab",
            slider-color="primary",
            fixed-tabs
          )
            v-tab {{ $t('pattern.myPatterns') }}
            v-tab-item(lazy)
              patterns(
                @edit="showEditPatternModal",
                @remove="showDeletePatternModal",
                @remove-selected="showDeleteSelectedPatternsModal"
              )
            v-tab {{ $t('pattern.corporatePatterns') }}
            v-tab-item(lazy)
              corporate-patterns(
                @edit="showEditPatternModal",
                @remove="showDeletePatternModal",
                @remove-selected="showDeleteSelectedPatternsModal"
              )
          patterns(
            v-else,
            @edit="showEditPatternModal",
            @remove="showDeletePatternModal",
            @remove-selected="showDeleteSelectedPatternsModal"
          )

    c-fab-expand-btn(@refresh="refresh", :has-access="hasAccessToCreatePattern")
      v-tooltip(top)
        v-btn(
          slot="activator",
          color="grey darken-1",
          fab,
          dark,
          small,
          @click.stop="showCreatePbehaviorPatternModal"
        )
          v-icon pause
        span Create pbehavior filter
      v-tooltip(top)
        v-btn(
          slot="activator",
          color="blue",
          fab,
          dark,
          small,
          @click.stop="showCreateEntityPatternModal"
        )
          v-icon perm_identity
        span Create entity filter
      v-tooltip(top)
        v-btn(
          slot="activator",
          color="error",
          fab,
          dark,
          small,
          @click.stop="showCreateAlarmPatternModal"
        )
          v-icon notification_important
        span Create alarm filter
</template>

<script>
import { MODALS } from '@/constants';

import { localQueryMixin } from '@/mixins/query-local/query';
import { entitiesPatternsMixin } from '@/mixins/entities/pattern';
import { entitiesCorporatePatternsMixin } from '@/mixins/entities/pattern/corporate';
import {
  permissionsTechnicalProfileCorporatePatternMixin,
} from '@/mixins/permissions/technical/profile/corporate-pattern';

import Patterns from '@/components/other/pattern/patterns.vue';
import CorporatePatterns from '@/components/other/pattern/corporate-patterns.vue';

export default {
  components: {
    Patterns,
    CorporatePatterns,
  },
  mixins: [
    localQueryMixin,
    entitiesPatternsMixin,
    entitiesCorporatePatternsMixin,
    permissionsTechnicalProfileCorporatePatternMixin,
  ],
  data() {
    return {
      activeTab: 0,
    };
  },
  computed: {
    hasAccessToCreatePattern() {
      return !this.activeTab || (this.activeTab && this.hasCreateAnyCorporatePatternAccess);
    },
  },
  methods: {
    refresh() {
      if (this.activeTab) { // TODO: change to keys
        return this.fetchCorporatePatternsListWithPreviousParams();
      }

      return this.fetchPatternsListWithPreviousParams();
    },

    showCreateAlarmPatternModal() {},

    showCreateEntityPatternModal() {},

    showCreatePbehaviorPatternModal() {},

    showEditPatternModal() {},

    showDeletePatternModal(id) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.removePattern({ id });

            return this.refresh();
          },
        },
      });
    },

    showDeleteSelectedPatternsModal(selected) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await Promise.all(selected.map(({ _id: id }) => this.removePattern({ id })));

            return this.refresh();
          },
        },
      });
    },
  },
};
</script>
