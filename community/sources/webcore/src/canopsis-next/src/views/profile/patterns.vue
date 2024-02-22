<template>
  <div>
    <c-page-header>{{ $t('pattern.patterns') }}</c-page-header>
    <v-layout wrap>
      <v-flex xs12>
        <v-card class="ma-2">
          <v-tabs
            v-if="hasReadAnyCorporatePatternAccess"
            v-model="activeTab"
            slider-color="primary"
            fixed-tabs
          >
            <v-tab :href="`#${$constants.PATTERN_TABS.patterns}`">
              {{ $t('pattern.myPatterns') }}
            </v-tab>
            <v-tab-item :value="$constants.PATTERN_TABS.patterns">
              <patterns
                @edit="showEditPatternModal"
                @remove="showDeletePatternModal"
                @remove-selected="showDeleteSelectedPatternsModal"
              />
            </v-tab-item>
            <v-tab :href="`#${$constants.PATTERN_TABS.corporatePatterns}`">
              {{ $t('pattern.corporatePatterns') }}
            </v-tab>
            <v-tab-item :value="$constants.PATTERN_TABS.corporatePatterns">
              <corporate-patterns
                @edit="showEditPatternModal"
                @remove="showDeletePatternModal"
                @remove-selected="showDeleteSelectedPatternsModal"
              />
            </v-tab-item>
          </v-tabs>
          <patterns
            v-else
            @edit="showEditPatternModal"
            @remove="showDeletePatternModal"
            @remove-selected="showDeleteSelectedPatternsModal"
          />
        </v-card>
      </v-flex>
    </v-layout>
    <c-fab-expand-btn
      :has-access="hasAccessToCreatePattern"
      @refresh="refresh"
    >
      <c-action-fab-btn
        :tooltip="createPbehaviorTitle"
        icon="pause"
        color="grey darken-1"
        top
        @click.stop="showCreatePbehaviorPatternModal"
      />
      <c-action-fab-btn
        :tooltip="createEntityTitle"
        icon="perm_identity"
        color="blue"
        top
        @click.stop="showCreateEntityPatternModal"
      />
      <c-action-fab-btn
        :tooltip="createAlarmTitle"
        icon="notification_important"
        color="error"
        top
        @click.stop="showCreateAlarmPatternModal"
      />
    </c-fab-expand-btn>
  </div>
</template>

<script>
import { MODALS, PATTERN_TABS, PATTERN_TYPES } from '@/constants';

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
      activeTab: PATTERN_TABS.patterns,
    };
  },
  computed: {
    isCorporatePatternsTab() {
      return this.activeTab === PATTERN_TABS.corporatePatterns;
    },

    hasAccessToCreatePattern() {
      return !this.isCorporatePatternsTab || this.hasCreateAnyCorporatePatternAccess;
    },

    createAlarmTitle() {
      return this.isCorporatePatternsTab
        ? this.$t('modals.createCorporateAlarmPattern.create.title')
        : this.$t('modals.createAlarmPattern.create.title');
    },

    createEntityTitle() {
      return this.isCorporatePatternsTab
        ? this.$t('modals.createCorporateEntityPattern.create.title')
        : this.$t('modals.createEntityPattern.create.title');
    },

    createPbehaviorTitle() {
      return this.isCorporatePatternsTab
        ? this.$t('modals.createCorporatePbehaviorPattern.create.title')
        : this.$t('modals.createPbehaviorPattern.create.title');
    },
  },
  methods: {
    refresh() {
      if (this.isCorporatePatternsTab) {
        return this.fetchCorporatePatternsListWithPreviousParams();
      }

      return this.fetchPatternsListWithPreviousParams();
    },

    showCreateAlarmPatternModal() {
      this.$modals.show({
        name: MODALS.createPattern,
        config: {
          pattern: { is_corporate: this.isCorporatePatternsTab },
          title: this.createAlarmTitle,
          type: PATTERN_TYPES.alarm,
          action: async (pattern) => {
            await this.createPattern({ data: pattern });

            return this.refresh();
          },
        },
      });
    },

    showCreateEntityPatternModal() {
      this.$modals.show({
        name: MODALS.createPattern,
        config: {
          pattern: { is_corporate: this.isCorporatePatternsTab },
          title: this.createEntityTitle,
          type: PATTERN_TYPES.entity,
          action: async (pattern) => {
            await this.createPattern({ data: pattern });

            return this.refresh();
          },
        },
      });
    },

    showCreatePbehaviorPatternModal() {
      this.$modals.show({
        name: MODALS.createPattern,
        config: {
          pattern: { is_corporate: this.isCorporatePatternsTab },
          title: this.createPbehaviorTitle,
          type: PATTERN_TYPES.pbehavior,
          action: async (pattern) => {
            await this.createPattern({ data: pattern });

            return this.refresh();
          },
        },
      });
    },

    getEditPatternModalTitle(pattern) {
      if (pattern.is_corporate) {
        return {
          [PATTERN_TYPES.alarm]: this.$t('modals.createCorporateAlarmPattern.edit.title'),
          [PATTERN_TYPES.entity]: this.$t('modals.createCorporateEntityPattern.edit.title'),
          [PATTERN_TYPES.pbehavior]: this.$t('modals.createCorporatePbehaviorPattern.edit.title'),
        }[pattern.type];
      }

      return {
        [PATTERN_TYPES.alarm]: this.$t('modals.createAlarmPattern.edit.title'),
        [PATTERN_TYPES.entity]: this.$t('modals.createEntityPattern.edit.title'),
        [PATTERN_TYPES.pbehavior]: this.$t('modals.createPbehaviorPattern.edit.title'),
      }[pattern.type];
    },

    showEditPatternModal(editablePattern) {
      this.$modals.show({
        name: MODALS.createPattern,
        config: {
          pattern: editablePattern,
          type: editablePattern.type,
          title: this.getEditPatternModalTitle(editablePattern),
          action: async (pattern) => {
            await this.updatePattern({ id: editablePattern._id, data: pattern });

            return this.refresh();
          },
        },
      });
    },

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
            await this.bulkRemovePatterns({
              data: selected.map(({ _id: id }) => ({ _id: id })),
            });

            return this.refresh();
          },
        },
      });
    },
  },
};
</script>
